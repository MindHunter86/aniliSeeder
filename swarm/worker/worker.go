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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/MindHunter86/aniliSeeder/anilibria"
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

	sync.RWMutex
	pingerDisable bool
}

var (
	gCli    *cli.Context
	gLog    *zerolog.Logger
	gCtx    context.Context
	gDeluge *deluge.Client
	gAniApi *anilibria.ApiClient
)

func NewWorker(ctx context.Context) swarm.Swarm {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gDeluge = gCtx.Value(utils.ContextKeyDelugeClient).(*deluge.Client)

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
	go m.run()

	gLog.Debug().Msg("starting grpc master server ...")
	return m.gserver.Serve(m.msession)
}

func (m *Worker) ping() (e error) {
	if _, e = m.msession.Ping(); e != nil {
		if e = m.connect(); e != nil {
			return
		}
	}

	// TODO
	return
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

	gLog.Debug().Msg("grpc master server has been setuped; initialize destructor")
	return
}

func (m *Worker) reconnect() (e error) {
	if e = m.msession.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}
	if e = m.rawconn.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}

	return
}

func (m *Worker) run() {
	<-gCtx.Done()
	gLog.Info().Msg("context done() has been caught; closing grpc server, mux session, tcp conn...")

	m.gserver.Stop()

	var e error
	if e = m.msession.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}
	if e = m.rawconn.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}
}

func (m *Worker) run2() (e error) {
	ticker := time.NewTicker(time.Second)
	ticker.Stop() // !!
	// todo refactor ?

	if gCli.Duration("grpc-ping-interval") != 0*time.Second {
		ticker.Reset(gCli.Duration("grpc-ping-interval"))
	}

	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-gCtx.Done():
			break LOOP
		case <-ticker.C:
			m.RLock()
			if m.pingerDisable {
				gLog.Debug().Msg("skipping ping call because of the last call is not completed yet")
				continue
			}
			m.RUnlock()

			// if e = m.ping(); e != nil {
			// 	gLog.Warn().Err(e).Msg("grpc ping has been failed; close application...")
			// 	return
			// }
		}
	}

	return
}

//

func (m *Worker) getNewRPCContext(d time.Duration) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{
		"x-worker-id": m.id,
	})

	return context.WithTimeout(
		metadata.NewOutgoingContext(context.Background(), md),
		d,
	)
}

func (m *Worker) getRegistrationRequest() (_ *pb.RegistrationRequest, e error) {
	var trrs []*structpb.Struct
	if trrs, e = m.getTorrents(); e != nil {
		return
	}

	return &pb.RegistrationRequest{
		WorkerVersion: gCli.App.Version,
		WDFreeSpace:   utils.CheckDirectoryFreeSpace(gCli.String("torrentfiles-dir")),
		Torrent:       trrs,
	}, e
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

// 	m.enablePing()
// 	return nil
// }

// func (m *Worker) disablePing() {
// 	m.Lock()
// 	m.pingerDisable = true
// 	m.Unlock()
// }
// func (m *Worker) enablePing() {
// 	m.Lock()
// 	m.pingerDisable = false
// 	m.Unlock()
// }