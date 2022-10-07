package swarm

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	uuid "github.com/satori/go.uuid"
)

type Minion struct{}

func NewMinion() *Minion {
	return &Minion{}
}

// func (*Minion) Bootstrap() error {
// 	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	c := pb.NewMasterClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	r, err := c.InitialPhase(ctx, &pb.MasterRequest{AccessKey: "fuckyounigga"})
// 	if err != nil {
// 		log.Fatalf("could not greet: %v", err)
// 	}
// 	log.Printf("Greeting: %s", r.GetVersion())
// 	return nil
// }

// =================

type Worker struct {
	Version     string
	WDFreeSpace uint64
	Torrents    map[string]*deluge.Torrent

	gConn        *grpc.ClientConn
	masterClient pb.MasterServiceClient

	id     string
	token  string
	config *WorkerConfig

	sync.RWMutex
	tickerPinging bool
}

type WorkerConfig struct{}

func NewWorker(ctx context.Context) Swarm {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gDeluge = gCtx.Value(utils.ContextKeyDelugeClient).(*deluge.Client)

	return &Worker{
		id:    uuid.NewV4().String(),
		token: gCli.String("swarm-master-secret"),
	}
}

func (m *Worker) Bootstrap() (e error) {
	var opts []grpc.DialOption

	if !gCli.Bool("grpc-insecure") {
		gLog.Debug().Msg("trying access to ca...")

		var cpool *x509.CertPool
		if cpool, e = m.getCACertPool(); e != nil {
			return
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cpool, "")))
	} else {
		gLog.Warn().Msg("ATTENTION! gRPC connection is unsecure! do at your own risk")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                gCli.Duration("http2-ping-time"),
		Timeout:             gCli.Duration("http2-ping-timeout"),
		PermitWithoutStream: true,
	}))
	opts = append(opts, grpc.WithBlock())

	opts = append(opts, grpc.WithTimeout(gCli.Duration("grpc-connect-timeout")))

	gLog.Debug().Msg("trying to connect to master...")
	if m.gConn, e = grpc.Dial(gCli.String("swarm-master-addr"), opts...); e != nil {
		return
	}

	gLog.Debug().Msg("seems that connection has been established")
	gLog.Debug().Msg("trying to complete init phase with master")

	m.masterClient = pb.NewMasterServiceClient(m.gConn)
	ctx, cancel := m.getNewRPCContext(time.Second)
	defer cancel()

	var req = &pb.RegistrationRequest{}
	if req, e = m.getRegistrationRequest(); e != nil {
		return
	}

	var rpl *pb.RegistrationReply
	if rpl, e = m.masterClient.Register(ctx, req); m.getRPCErrors(e) != nil {
		gLog.Warn().Err(e).Msg("there is some errors while processing grpc request in registration phase")
		return
	}

	gLog.Debug().Msg("registration has been completed; parsing config data from master...")

	// var cfg *WorkerConfig
	if _, e = m.parseRegistrationReply(rpl); e != nil {
		return
	}

	// if e = m.Setup(cfg); e != nil {
	// 	return
	// }

	gLog.Debug().Msg("the registration phase finished; waiting for commands from the master")
	return m.run()
}

// TODO
// func (*Worker) Setup(cfg *WorkerConfig) (e error) {
// 	return
// }

func (m *Worker) run() error {
	ticker := time.NewTicker(time.Second)
	ticker.Stop() // !!
	// todo refactor ?

	if gCli.Duration("grpc-ping-interval") != 0*time.Second {
		ticker.Reset(gCli.Duration("grpc-ping-interval"))
	}

LOOP:
	for {
		select {
		case <-gCtx.Done():
			break LOOP
		case <-ticker.C:
			m.RLock()
			if m.tickerPinging {
				gLog.Debug().Msg("skipping ping call because of the last call is not completed yet")
				continue
			}
			m.RUnlock()

			if e := m.ping(); e != nil {
				gLog.Warn().Err(e).Msg("grpc ping has been failed")
			}
		}
	}

	ticker.Stop()
	return m.desctruct()
}

func (m *Worker) desctruct() error {
	return m.gConn.Close()
}

func (*Worker) getCACertPool() (*x509.CertPool, error) {
	// TODO
	// if gCli.String(CA-PATH) != "" -->> loadCAFromFS()
	return x509.SystemCertPool()
}

// TODO
// if gCli.String(CA-PATH) != "" -->> loadCAFromFS()
//--------------------------------------------------
// func (*Worker) loadTLSCertificate(path string) (_ io.Reader, e error) {
// 	var fInfo os.FileInfo

// 	if fInfo, e = os.Stat(path); e != nil {
// 		if os.IsNotExist(e) {
// 			gLog.Error().Err(e).Msg("could not load ca because of invalid given file path")
// 			return
// 		}

// 		return
// 	}

// 	if fInfo.IsDir() {
// 		gLog.Error().Msg("could not load ca because of give file path is a directory")
// 	}

// 	return
// }

func (m *Worker) getNewRPCContext(d time.Duration) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{
		"x-worker-id":    m.id,
		"x-access-token": m.token,
	})

	return context.WithTimeout(
		metadata.NewOutgoingContext(context.Background(), md),
		d,
	)
}

func (*Worker) parseRegistrationReply(rpl *pb.RegistrationReply) (_ *WorkerConfig, e error) {
	var cfg *WorkerConfig

	var buf []byte
	if buf, e = json.Marshal(rpl); e != nil {
		return
	}

	if e = json.Unmarshal(buf, &cfg); e != nil {
		return
	}

	return cfg, e
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

func (m *Worker) ping() (e error) {
	timer := time.Now()

	m.Lock()
	m.tickerPinging = true
	m.Unlock()

	defer func() {
		m.Lock()
		m.tickerPinging = false
		m.Unlock()
	}()

	ctx, cancel := m.getNewRPCContext(gCli.Duration("grpc-ping-timeout"))
	defer cancel()

	if _, e = m.masterClient.Ping(ctx, &emptypb.Empty{}); m.getRPCErrors(e) != nil {
		return
	}

	gLog.Debug().Str("ping_time", time.Since(timer).String()).Msg("ping/pong method completed")
	return
}

func (m *Worker) getRPCErrors(err error) error {
	estatus, _ := status.FromError(err)

	switch estatus.Code() {
	case codes.OK:
		return nil

	// !! EXPERIMENTAL
	case codes.Unavailable:
		gLog.Warn().Msg("trying to reconnect to the master server...")
		m.gConn.ResetConnectBackoff()

	default:
		gLog.Warn().Str("error_code", estatus.Code().String()).Str("error_message", estatus.Message()).
			Msg("abnormal response from master server")
	}

	return err
}

// Debug func
// func (*Worker) CheckGRPCPayload(payload []*structpb.Struct) (_ bool, e error) {

// 	var trrs = make([]*deluge.Torrent, 100)

// 	var buf []byte
// 	if buf, e = json.Marshal(payload); e != nil {
// 		return
// 	}

// 	if e = json.Unmarshal(buf, &trrs); e != nil {
// 		return
// 	}

// 	for _, trr := range trrs {
// 		log.Println(trr.Name)
// 	}

// 	return true, e
// }
