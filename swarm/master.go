package swarm

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Master struct {
	pb.UnimplementedMasterServiceServer

	ln      net.Listener
	gserver *grpc.Server
	workers map[string]*Worker
}

func NewMaster(ctx context.Context) *Master {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gAniApi = gCtx.Value(utils.ContextKeyAnilibriaClient).(*anilibria.ApiClient)

	return &Master{}
}

func (m *Master) Bootstrap() (e error) {
	gLog.Debug().Msg("generating pub\\priv key pair...")

	var crt tls.Certificate
	if crt, e = m.getTLSCertificate(); e != nil {
		return
	}

	gLog.Debug().Msg("trying to open grpc socket for master listening...")

	if m.ln, e = net.Listen("tcp", gCli.String("swarm-master-listen")); e != nil {
		return
	}

	gLog.Debug().Msg("grpc socket seems is ok, setuping grpc...")

	var creds = credentials.NewServerTLSFromCert(&crt)
	m.gserver = grpc.NewServer(grpc.Creds(creds))
	pb.RegisterMasterServiceServer(m.gserver, m)

	gLog.Debug().Msg("grpc server has been setuped; starting listening for worker connections...")

	defer func() {
		if err := m.close(); err != nil {
			gLog.Warn().Err(err).Msg("there are somee errors after closing grpc server socket")
		}
	}()

	gLog.Debug().Msg("grpc master server has been setuped")
	return
}

func (m *Master) Serve(done func()) (e error) {
	defer done()

	gLog.Debug().Msg("starting grpc master server ...")
	return m.gserver.Serve(m.ln)
}

func (m *Master) close() error {
	<-gCtx.Done()
	gLog.Warn().Msg("context done() has been caught; closing grpc server socket...")

	m.gserver.Stop()
	return m.ln.Close()
}

func (m *Master) getTLSCertificate() (_ tls.Certificate, e error) {
	// TODO
	// if gCli.String("swarm-master-custom-ca") != ""
	// get key, get pub, return

	var cbytes, kbytes []byte
	if cbytes, kbytes, e = m.createPublicPrivatePair(); e != nil {
		return
	}

	return tls.X509KeyPair(cbytes, kbytes)
}

func (*Master) createPublicPrivatePair() (_, _ []byte, e error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		return
	}

	certBytes, e := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if e != nil {
		return
	}

	var cbuf = new(bytes.Buffer)
	pem.Encode(cbuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var p []byte
	if p, e = x509.MarshalECPrivateKey(priv); e != nil {
		return
	}

	var pbuf = new(bytes.Buffer)
	pem.Encode(pbuf, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p,
	})

	if gCli.Bool("http-debug") {
		fmt.Println("\n" + cbuf.String())
		fmt.Println("\n" + pbuf.String())
	}

	return cbuf.Bytes(), pbuf.Bytes(), nil
}

// func (*Master) InitialPhase(ctx context.Context, in *pb.MasterRequest) (*pb.MasterReply, error) {
// 	log.Printf("Received: %v", in.GetAccessKey())
// 	return &pb.MasterReply{Version: "Hello " + in.GetAccessKey()}, nil
// }

// func (*Master) Bootstrap() error {
// 	lis, err := net.Listen("tcp", "localhost:8081")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	pb.RegisterMasterServer(s, &Master{})
// 	log.Printf("server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}

// 	return err
// }

func (*Master) Register(ctx context.Context, req *pb.RegistrationRequest) (_ *pb.RegistrationReply, e error) {
	//
	return
}
