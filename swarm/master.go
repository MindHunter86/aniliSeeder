package swarm

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"google.golang.org/grpc"
)

type Master struct {
	pb.UnimplementedMasterServer
}

func NewMaster() *Master {
	return &Master{}
}

func (*Master) InitialPhase(ctx context.Context, in *pb.MasterRequest) (*pb.MasterReply, error) {
	log.Printf("Received: %v", in.GetAccessKey())
	return &pb.MasterReply{Version: "Hello " + in.GetAccessKey()}, nil
}

func (*Master) Bootstrap() error {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServer(s, &Master{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
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
