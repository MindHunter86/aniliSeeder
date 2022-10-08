package master

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/hashicorp/yamux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	errWorkerInvalidId = errors.New("the worker service is replied with invalid id")
)

type worker struct {
	msess    *yamux.Session
	gconn    *grpc.ClientConn
	gservice pb.WorkerServiceClient

	id    string
	token string
	trrs  []*deluge.Torrent
}

func newWorker(ms *yamux.Session) *worker {
	return &worker{
		msess: ms,
	}
}

func (m *worker) newServiceRequest(d time.Duration) (context.Context, context.CancelFunc) {
	mac := hmac.New(sha256.New, []byte(gCli.String("swarm-master-secret")))
	mac.Write([]byte(gMasterId))

	md := metadata.New(map[string]string{
		"x-master-id":           gMasterId,
		"x-authentication-hash": string(mac.Sum(nil)),
	})

	return context.WithTimeout(
		metadata.NewOutgoingContext(context.Background(), md),
		d,
	)
}

func (m *worker) authorizeSerivceReply(ctx context.Context) (_ string, e error) {
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

	gLog.Debug().Str("worker_ip", p.Addr.String()).Str("worker_id", id[0]).
		Str("worker_ua", md.Get("user-agent")[0]).Msg("worker reply accepted, authorizing...")

	ah := md.Get("x-authentication-hash")
	if len(ah) != 1 {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(ah[0]) == "" {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}

	mac := hmac.New(sha256.New, []byte(gCli.String("swarm-master-secret")))
	mac.Write([]byte(id[0]))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal([]byte(ah[0]), expectedMAC) {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("worker_id", id[0]).Msg("the worker's reply has been authorized")
	return id[0], nil
}

func (m *worker) getRPCErrors(err error) error {
	estatus, _ := status.FromError(err)

	switch estatus.Code() {
	case codes.OK:
		return nil

	// // !! EXPERIMENTAL
	// case codes.Unavailable:
	// 	gLog.Warn().Msg("trying to reconnect to the master server...")
	// 	m.gconn.ResetConnectBackoff()

	default:
		gLog.Warn().Str("error_code", estatus.Code().String()).Str("error_message", estatus.Message()).
			Msg("abnormal response from master server")
	}

	return err
}

func (m *worker) getInitialServiceData() (_ string, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	var rpl *pb.InitReply
	if rpl, e = m.gservice.Init(ctx, &emptypb.Empty{}); m.getRPCErrors(e) != nil {
		return
	}

	if m.id, e = m.authorizeSerivceReply(ctx); e != nil {
		return
	}

	if m.id != rpl.GetWorkerId() {
		gLog.Warn().Str("worker_id", m.id).Msg("abnormal reply from init method of the orker service; drop worker")
		return m.id, errWorkerInvalidId
	}

	return
}

type workerPool struct {
	pool sync.Pool

	sync.RWMutex
	workers map[string]*worker
}

func newWorkerPool() *workerPool {
	return &workerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &worker{}
			},
		},
		workers: make(map[string]*worker),
	}
}

func (m *workerPool) newWorker(msess *yamux.Session) (wrk *worker, e error) {
	var opts []grpc.DialOption

	if !gCli.Bool("grpc-insecure") {
		gLog.Debug().Msg("trying access to ca...")

		var cpool *x509.CertPool
		if cpool, e = getCACertPool(); e != nil {
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

	opts = append(opts, grpc.WithDialer(func(tgt string, timeout time.Duration) (net.Conn, error) {
		return msess.Open()
	}))

	gLog.Debug().Msg("trying to connect to worker over a mux...")

	//todo : remove sync Pool
	// wrk := m.pool.Get().(*worker)
	wrk = newWorker(msess)

	ctx, cancel := context.WithTimeout(context.Background(), gCli.Duration("grpc-connect-timeout"))
	defer cancel()

	if wrk.gconn, e = grpc.DialContext(ctx, "", opts...); e != nil {
		return
	}

	gLog.Debug().Msg("connection with the master server has been established; registering grpc services...")
	wrk.gservice = pb.NewWorkerServiceClient(wrk.gconn)

	var wid string
	if wid, e = wrk.getInitialServiceData(); e != nil {
		return
	}

	gLog.Debug().Msg("registration completed; the worker has been initialized")
	m.Lock()
	m.workers[wid] = wrk
	m.Unlock()

	return
}

func (m *workerPool) dropWorker(w *worker) (e error) {
	if e = w.gconn.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}
	if e = w.msess.Close(); e != nil {
		gLog.Warn().Err(e).Msg("")
	}

	m.Lock()
	m.workers[w.id] = nil
	m.Unlock()

	return
}
