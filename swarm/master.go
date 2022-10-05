package swarm

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Master struct {
	pb.UnimplementedMasterServiceServer

	ln      net.Listener
	gserver *grpc.Server

	sync.RWMutex
	workers map[string]*Worker
}

func NewMaster(ctx context.Context) *Master {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gAniApi = gCtx.Value(utils.ContextKeyAnilibriaClient).(*anilibria.ApiClient)

	return &Master{
		workers: make(map[string]*Worker),
	}
}

func (m *Master) Bootstrap() (e error) {
	gLog.Debug().Msg("generating pub\\priv key pair...")

	// var crt tls.Certificate
	// if crt, e = m.getTLSCertificate(); e != nil {
	// 	return
	// }

	gLog.Debug().Msg("trying to open grpc socket for master listening...")

	if m.ln, e = net.Listen("tcp", gCli.String("swarm-master-listen")); e != nil {
		return
	}

	gLog.Debug().Msg("grpc socket seems is ok, setuping grpc...")

	// var creds = credentials.NewServerTLSFromCert(&crt)
	// m.gserver = grpc.NewServer(grpc.Creds(creds))
	m.gserver = grpc.NewServer()
	pb.RegisterMasterServiceServer(m.gserver, m)

	gLog.Debug().Msg("grpc server has been setuped; starting listening for worker connections...")

	go func() {
		if err := m.close(); err != nil {
			gLog.Warn().Err(err).Msg("there are some errors after closing grpc server socket")
		}
	}()

	gLog.Debug().Msg("grpc master server has been setuped")

	gLog.Debug().Msg("starting grpc master server ...")
	return m.gserver.Serve(m.ln)
}

func (m *Master) close() error {
	<-gCtx.Done()
	gLog.Info().Msg("context done() has been caught; closing grpc server socket...")

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

func (m *Master) Register(ctx context.Context, req *pb.RegistrationRequest) (_ *pb.RegistrationReply, e error) {
	p, _ := peer.FromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)

	gLog.Info().Str("worker_id", req.WorkerId).Str("worker_ip", p.Addr.String()).Strs("client_ua", md.Get("user-agent")).
		Msg("new client has been connected")
	gLog.Info().Str("worker_id", req.WorkerId).Msg("new client validation phase running...")

	switch {
	case req.GetWorkerId() == "":
		return nil, status.Errorf(codes.InvalidArgument, "")
	case m.workers[req.GetWorkerId()] != nil:
		return nil, status.Errorf(codes.AlreadyExists, "")
	case req.AccessSecret != gCli.String("swarm-master-secret"):
		return nil, status.Errorf(codes.Unauthenticated, "")
	case req.WorkerVersion != gCli.App.Version:
		gLog.Warn().Str("worker_id", req.WorkerId).Str("worker_ver", req.WorkerVersion).
			Msg("connected client has higher/lower version")
	case req.WDFreeSpace == 0:
		return nil, status.Errorf(codes.InvalidArgument, "")
	}

	gLog.Info().Str("worker_id", req.WorkerId).Msg("trying parse torrent list from new client...")
	var trrs []*deluge.Torrent

	var buf []byte
	if buf, e = json.Marshal(req.Torrent); e != nil {
		gLog.Error().Err(e).Msg("there is an error while proccessing new client's torrent list")
		return nil, status.Errorf(codes.Internal, "")
	}

	if e = json.Unmarshal(buf, &trrs); e != nil {
		gLog.Error().Err(e).Msg("there is an error while proccessing new client's torrent list")
		return nil, status.Errorf(codes.Internal, "")
	}

	gLog.Debug().Str("worker_id", req.WorkerId).Int("torrents_count", len(trrs)).Msg("torrent list parsing from the client has been completed")
	gLog.Info().Str("worker_id", req.WorkerId).Msg("client validation seems ok; registering new worker...")

	var wtrrs = make(map[string]*deluge.Torrent)
	for _, t := range trrs {
		if t == nil || t.Hash == "" {
			gLog.Warn().Msg("there is strange torrent in the list from the client")
		}

		wtrrs[t.Hash] = &deluge.Torrent{
			ActiveTime:    t.ActiveTime,
			Ratio:         t.Ratio,
			IsFinished:    t.IsFinished,
			IsSeed:        t.IsSeed,
			Name:          t.Name,
			NumPeers:      t.NumPeers,
			NumPieces:     t.NumPieces,
			NumSeeds:      t.NumSeeds,
			PieceLength:   t.PieceLength,
			SeedingTime:   t.SeedingTime,
			State:         t.State,
			TotalPeers:    t.TotalPeers,
			TotalSeeds:    t.TotalSeeds,
			TotalDone:     t.TotalDone,
			TotalUploaded: t.TotalUploaded,
			TotalSize:     t.TotalSize,
		}
	}

	log.Println(req.WDFreeSpace)
	log.Println(req.WorkerVersion)
	log.Println(wtrrs)

	m.Lock()
	m.workers[req.WorkerId] = &Worker{
		Version:     req.WorkerVersion,
		WDFreeSpace: req.WDFreeSpace,
		Torrents:    wtrrs,
	}
	m.Unlock()

	gLog.Info().Str("worker_id", req.WorkerId).Msg("new client registration has been completed")
	return &pb.RegistrationReply{Config: &structpb.Struct{}}, e
}
