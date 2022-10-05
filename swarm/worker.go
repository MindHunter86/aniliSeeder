package swarm

import (
	"context"
	"crypto/x509"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	uuid "github.com/satori/go.uuid"
)

type Minion struct{}

func NewMinion() *Minion {
	return &Minion{}
}

func (*Minion) Bootstrap() error {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.InitialPhase(ctx, &pb.MasterRequest{AccessKey: "fuckyounigga"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetVersion())
	return nil
}

// =================

type Worker struct {
	gConn *grpc.ClientConn

	id string
}

func NewWorker(c *cli.Context, l *zerolog.Logger, ctx context.Context) *Worker {
	gCli, gLog, gCtx = c, l, ctx
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

	c := pb.NewMasterServiceClient(m.gConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.RegistrationRequest{}

	var rpl *pb.RegistrationReply
	if rpl, e = c.Register(ctx, nil); e != nil {
		return
	}

	//

	return m.run()
}

func (m *Worker) run() error {
	<-gCtx.Done()
	return m.desctruct()
}

func (m *Worker) desctruct() error {
	return m.gClient.Close()
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

	return &pb.RegistrationRequest{
		WorkerId:      m.id,
		WorkerVersion: gCli.App.Version,
		AccessSecret:  gCli.String("swarm-master-secret"),
		// !!
	}, e
}
