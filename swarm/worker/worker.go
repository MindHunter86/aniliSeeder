package worker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/swarm"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/hashicorp/yamux"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	uuid "github.com/satori/go.uuid"
)

type Worker struct {
	Version     string
	WDFreeSpace uint64
	Torrents    map[string]*deluge.Torrent

	rawconn  net.Conn
	msession *yamux.Session
	gserver  *grpc.Server

	id string
}

var (
	gCli    *cli.Context
	gLog    *zerolog.Logger
	gCtx    context.Context
	gDeluge *deluge.Client
	gAbort  context.CancelFunc
)

func NewWorker(ctx context.Context) swarm.Swarm {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gDeluge = gCtx.Value(utils.ContextKeyDelugeClient).(*deluge.Client)
	gAbort = gCtx.Value(utils.ContextKeyAbortFunc).(context.CancelFunc)

	return &Worker{
		id: uuid.NewV4().String(),
	}
}

func (m *Worker) Bootstrap() (e error) {
	if e = m.connect(); e != nil {
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Wait()
	defer gLog.Debug().Msg("waiting for destructor...")

	return m.run(wg.Done)
}

func (m *Worker) connect() (e error) {
	gLog.Debug().Str("master_addr", gCli.String("swarm-master-addr")).
		Msg("trying to establish raw tcp connection with the master server")

	if m.rawconn, e = net.DialTimeout("tcp", gCli.String("swarm-master-addr"), gCli.Duration("grpc-connect-timeout")); e != nil {
		return
	}

	gLog.Debug().Str("master_addr", gCli.String("swarm-master-addr")).Msg("trying to initialize mux session...")
	if m.msession, e = yamux.Server(m.rawconn, yamux.DefaultConfig()); e != nil {
		return
	}

	gLog.Debug().Msg("grpc socket seems is ok, setuping grpc...")

	var opts []grpc.ServerOption

	if !gCli.Bool("grpc-insecure") {
		gLog.Debug().Msg("generating pub\\priv key pair...")

		var crt tls.Certificate
		if crt, e = m.getTLSCertificate(); e != nil {
			return
		}

		var creds = credentials.NewServerTLSFromCert(&crt)
		opts = append(opts, grpc.Creds(creds))

	} else {
		gLog.Warn().Msg("ATTENTION! gRPC connection is unsecure! do at your own risk")
	}

	if gCli.Duration("http2-conn-max-age") != 0*time.Second {
		gLog.Debug().Msg("set keepalive for the master server...")

		enforcement := keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}

		opts = append(opts, grpc.KeepaliveEnforcementPolicy(enforcement))
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      gCli.Duration("http2-conn-max-age"),
			MaxConnectionAgeGrace: gCli.Duration("http2-conn-max-age") + 10*time.Second,
		}))
	}

	var wservice = NewWorkerService(m)

	m.gserver = grpc.NewServer(opts...)
	pb.RegisterWorkerServiceServer(m.gserver, wservice)

	gLog.Debug().Msg("grpc master server has been started")
	return
}

func (m *Worker) reconnect() (e error) {
	if e = m.msession.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}

	return m.connect()
}

func (m *Worker) serve(done func()) {
	gLog.Debug().Msg("starting grpc master server ...")
	if e := m.gserver.Serve(m.msession); e != nil {
		gLog.Warn().Err(e).Msg("got some errors from grpc serve")
	}

	gLog.Debug().Msg("grpc server has been stopped")
	done()
}

func (m *Worker) run(done func()) error {
	defer done()

	var wg sync.WaitGroup

	defer wg.Wait()
	defer gLog.Debug().Msg("waiting for wg.Done")

	wg.Add(1)
	go m.serve(wg.Done)

	var reconnFreeze bool

	// avoid panic if cli arg returns 0
	ticker := time.NewTicker(time.Second)
	ticker.Stop()

	if i := gCli.Duration("grpc-ping-interval"); i != 0 {
		ticker.Reset(i)
	}

LOOP:
	for {
		select {
		case <-gCtx.Done():
			gLog.Info().Msg("context done() has been caught; closing grpc server, mux session, tcp conn...")
			break LOOP
		case <-ticker.C:
			ticker.Stop()

			if reconnFreeze {
				gLog.Debug().Msg("restoring ping ticker after reconnection hold")
				reconnFreeze = false

				ticker.Reset(gCli.Duration("grpc-ping-interval"))
				continue
			}

			var e error
			if reconnFreeze, e = m.ping(); e != nil {
				gLog.Warn().Err(e).Msg("aborting application due to ping and reconnection failures")
				gAbort()
			}

			if reconnFreeze {
				wg.Add(1)
				go m.serve(wg.Done)

				if h := gCli.Duration("grpc-ping-reconnect-hold"); h != 0*time.Second {
					gLog.Debug().Msg("reconnection detected; holding for N seconds")

					ticker.Stop()
					ticker.Reset(h)

					continue
				}
			}

			ticker.Reset(gCli.Duration("grpc-ping-interval"))
		}
	}

	gLog.Debug().Msg("trying to stop grpc server")
	m.gserver.Stop()

	gLog.Debug().Msg("trying to stop mux session")
	return m.msession.Close()
}

func (m *Worker) ping() (reconn bool, e error) {
	var d time.Duration
	if d, e = m.msession.Ping(); e != nil {
		gLog.Debug().Err(e).Msg("got an error while pinging the mux session")

		if gCli.Bool("grpc-disable-reconnect") {
			gLog.Debug().Msg("aborting application because of grpc-disable-reconnect flag set")
			return
		}

		if e = m.reconnect(); e != nil {
			gLog.Debug().Err(e).Msg("got an error while reconnecting to the master server")
			return
		}

		reconn = true
	}

	gLog.Debug().Str("ping_time", d.String()).Msg("mux ping duration")
	return
}

func (*Worker) getTorrents() (_ []*structpb.Struct, e error) {
	var trrs []*deluge.Torrent
	var strmap = make([]*structpb.Struct, len(trrs))

	if trrs, e = gDeluge.GetTorrentsV2(); e != nil {
		return
	}

	var buf []byte
	if buf, e = json.Marshal(trrs); e != nil {
		return
	}

	if e = json.Unmarshal(buf, &strmap); e != nil {
		return
	}

	return strmap, e
}

func (*Worker) GetConnectedWorkers() (_ map[string]*swarm.SwarmWorker) {
	return nil
}

func (*Worker) IsMaster() bool {
	return false
}

// todo
// ? refactor
// func (m *Worker) ping() (e error) {
// 	timer := time.Now()

// 	m.disablePing()

// 	ctx, cancel := m.getNewRPCContext(gCli.Duration("grpc-ping-timeout"))
// 	defer cancel()

// 	if _, e = m.masterClient.Ping(ctx, &emptypb.Empty{}); m.getRPCErrors(e) == nil {
// 		gLog.Debug().Str("ping_time", time.Since(timer).String()).Msg("ping/pong method completed")

// 		m.enablePing()
// 		return
// 	}

// 	if code, ok := status.FromError(e); !ok || code.Code() == codes.PermissionDenied {
// 		gLog.Warn().Msg("the master says that worker isn't registered")

// 		// if e := m.registerInMaster(); e != nil {
// 		// 	gLog.Error().Err(e).Msg("reregistration has been failed")
// 		// 	return e
// 		// }

// 		gLog.Warn().Msg("registraion has been completed")
// 	}
