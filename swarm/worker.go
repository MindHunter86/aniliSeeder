package swarm

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	gConn *grpc.ClientConn

	id string
}

func NewWorker(ctx context.Context) Swarm {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gDeluge = gCtx.Value(utils.ContextKeyDelugeClient).(*deluge.Client)

	return &Worker{}
}

func (m *Worker) Bootstrap() (e error) {
	gLog.Debug().Msg("trying access to ca...")
	var cpool *x509.CertPool
	if cpool, e = m.getCACertPool(); e != nil {
		return
	}

	gLog.Debug().Msg("trying to connect to master...")
	if m.gConn, e = grpc.Dial(gCli.String("swarm-master-addr"), grpc.WithTransportCredentials(
		credentials.NewClientTLSFromCert(cpool, "")),
	); e != nil {
		return
	}

	gLog.Debug().Msg("seems that connection has been established")
	gLog.Debug().Msg("trying to complete init phase with master")

	gLog.Debug().Msg("DEBUGGING!")

	var req = &pb.RegistrationRequest{}
	if req, e = m.getRegistrationRequest(); e != nil {
		return
	}

	fmt.Println(req.Torrent)

	gLog.Debug().Msg("DEBUGGING2!")

	// c := pb.NewMasterServiceClient(m.gConn)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// var req = &pb.RegistrationRequest{}

	// var rpl *pb.RegistrationReply
	// if rpl, e = c.Register(ctx, nil); e != nil {
	// 	return
	// }

	//

	return m.run()
}

func (m *Worker) run() error {
	<-gCtx.Done()
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

func (m *Worker) getRegistrationRequest() (_ *pb.RegistrationRequest, e error) {
	m.id = uuid.NewV4().String()

	var trrs []*structpb.Struct
	if trrs, e = m.getTorrents(); e != nil {
		return
	}

	// if _, e = m.CheckGRPCPayload(trrs); e != nil {
	// 	return
	// }

	return &pb.RegistrationRequest{
		WorkerId:      m.id,
		WorkerVersion: gCli.App.Version,
		AccessSecret:  gCli.String("swarm-master-secret"),
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
