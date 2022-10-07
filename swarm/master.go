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
	"strings"
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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

	gLog.Debug().Msg("trying to open grpc socket for master listening...")
	if m.ln, e = net.Listen("tcp", gCli.String("swarm-master-listen")); e != nil {
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

	m.gserver = grpc.NewServer(opts...)
	pb.RegisterMasterServiceServer(m.gserver, m)

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

func (m *Master) isWorkerRegistered(id string) bool {
	return m.workers[id] != nil
}

func (m *Master) authorizeWorker(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "")
	}

	id := md.Get("x-worker-id")
	if len(id) != 1 {
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(id[0]) == "" {
		return "", status.Errorf(codes.InvalidArgument, "")
	}

	gLog.Debug().Str("worker_ip", p.Addr.String()).Str("worker_id", md.Get("x-worker-id")[0]).
		Str("worker_ua", md.Get("user-agent")[0]).Msg("worker connect accepted, authorizing...")

	ak := md.Get("x-access-token")
	if len(ak) != 1 {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(ak[0]) == "" {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if ak[0] != gCli.String("swarm-master-secret") {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("worker_id", md.Get("x-worker-id")[0]).Msg("the worker's connect has been authorized")
	return id[0], nil
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

func (m *Master) Ping(ctx context.Context, _ *emptypb.Empty) (_ *emptypb.Empty, _ error) {
	wid, e := m.authorizeWorker(ctx)
	if e != nil {
		return &emptypb.Empty{}, e
	}

	if !m.isWorkerRegistered(wid) {
		gLog.Info().Str("worker_id", wid).Msg("worker is not registered, returning 403...")
		return nil, status.Errorf(codes.PermissionDenied, "")
	}

	gLog.Info().Str("worker_id", wid).Msg("received ping from worker")
	return &emptypb.Empty{}, nil
}

func (m *Master) Register(ctx context.Context, req *pb.RegistrationRequest) (_ *pb.RegistrationReply, e error) {
	var wid string
	if wid, e = m.authorizeWorker(ctx); e != nil {
		return
	}

	gLog.Info().Str("worker_id", wid).Msg("new client validation phase running...")

	switch {
	case m.workers[wid] != nil:
		return nil, status.Errorf(codes.AlreadyExists, "")
	case req.WorkerVersion != gCli.App.Version:
		gLog.Warn().Str("worker_id", wid).Str("worker_ver", req.WorkerVersion).
			Msg("connected client has higher/lower version")
	case req.WDFreeSpace == 0:
		return nil, status.Errorf(codes.InvalidArgument, "")
	}

	gLog.Info().Str("worker_id", wid).Msg("trying parse torrent list from new client...")
	var trrs []*deluge.Torrent

	var buf []byte
	if buf, e = json.Marshal(req.Torrent); e != nil {
		gLog.Error().Err(e).Msg("there is an error while processing new client's torrent list")
		return nil, status.Errorf(codes.Internal, "")
	}

	if e = json.Unmarshal(buf, &trrs); e != nil {
		gLog.Error().Err(e).Msg("there is an error while processing new client's torrent list")
		return nil, status.Errorf(codes.Internal, "")
	}

	gLog.Debug().Str("worker_id", wid).Int("torrents_count", len(trrs)).Msg("torrent list parsing from the client has been completed")
	gLog.Info().Str("worker_id", wid).Msg("client validation seems ok; registering new worker...")

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

	if gCli.Bool("http-debug") {
		log.Println(req.WDFreeSpace)
		log.Println(req.WorkerVersion)
		log.Println(wtrrs)
	}

	m.Lock()
	m.workers[wid] = &Worker{
		Version:     req.WorkerVersion,
		WDFreeSpace: req.WDFreeSpace,
		Torrents:    wtrrs,
	}
	m.Unlock()

	gLog.Info().Str("worker_id", wid).Msg("new client registration has been completed")
	return &pb.RegistrationReply{Config: &structpb.Struct{}}, e
}
