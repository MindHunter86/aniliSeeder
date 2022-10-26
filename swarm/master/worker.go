package master

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
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

	masterId string

	trrs        []*deluge.Torrent
	version     string
	wdFreeSpace uint64

	mu sync.RWMutex
	id string
}

func newWorker(ms *yamux.Session, mid string) *worker {
	return &worker{
		msess:    ms,
		masterId: mid,
	}
}

func (m *worker) connect() (e error) {
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

	opts = append(opts, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return m.msess.Open()
	}))

	gLog.Debug().Msg("trying to connect to worker over a mux...")

	ctx, cancel := context.WithTimeout(context.Background(), gCli.Duration("grpc-connect-timeout"))
	defer cancel()

	if m.gconn, e = grpc.DialContext(ctx, "", opts...); e != nil {
		return
	}

	gLog.Debug().Msg("connection with the master server has been established; registering grpc services...")
	m.gservice = pb.NewWorkerServiceClient(m.gconn)

	if _, e = m.getInitialServiceData(); e != nil {
		gLog.Debug().Err(e).Msg("got an error while gathering initial service data")
		return
	}

	gLog.Debug().Msg("registration completed; the worker has been initialized")
	return
}

func (m *worker) reconnect() (e error) {
	if _, e = m.msess.Ping(); e != nil {
		return
	}

	m.gservice = nil
	if e = m.gconn.Close(); e != nil {
		gLog.Debug().Err(e).Msg("got an error while processing grpc conn.Close()")
	}

	if e = m.connect(); e != nil {
		return
	}

	return
}

// func (m *worker) disconnect() (e error) {
// 	if e = m.gconn.Close(); e != nil {
// 		gLog.Warn().Err(e).Msg("")
// 	}
// 	if e = m.msess.Close(); e != nil {
// 		gLog.Warn().Err(e).Msg("")
// 	}

// 	return
// }

func (m *worker) getId() (id string) {
	m.mu.RLock()
	id = m.id
	m.mu.RUnlock()

	return
}

func (m *worker) setId(id string) {
	m.mu.Lock()
	m.id = id
	m.mu.Unlock()
}

func (m *worker) newServiceRequest(d time.Duration) (context.Context, context.CancelFunc) {
	mac := hmac.New(sha256.New, []byte(gCli.String("master-secret")))
	io.WriteString(mac, m.masterId)

	md := metadata.New(map[string]string{
		"x-master-id":           m.masterId,
		"x-authentication-hash": hex.EncodeToString(mac.Sum(nil)),
	})

	return context.WithTimeout(
		metadata.NewOutgoingContext(context.Background(), md),
		d,
	)
}

func (m *worker) authorizeSerivceReply(md *metadata.MD) (e error) {
	id := md.Get("x-worker-id")
	if len(id) != 1 {
		return status.Errorf(codes.InvalidArgument, "there is no metadata in the reply")
	}
	if strings.TrimSpace(id[0]) == "" {
		return status.Errorf(codes.InvalidArgument, "there is no worker-id in the reply")
	}
	if m.getId() != "" && m.getId() != id[0] {
		return status.Errorf(codes.InvalidArgument, "given worker id is not equal to registration id")
	} else if m.getId() == "" {
		m.setId(id[0])
	}

	gLog.Debug().Str("worker_id", m.getId()).Msg("worker reply accepted, authorizing...")

	ah := md.Get("x-authentication-hash")
	if len(ah) != 1 {
		gLog.Info().Str("worker_id", m.getId()).Msg("worker authorization failed")
		return status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(ah[0]) == "" {
		gLog.Info().Str("worker_id", m.getId()).Msg("worker authorization failed")
		return status.Errorf(codes.InvalidArgument, "")
	}

	mmac, e := hex.DecodeString(ah[0])
	if e != nil {
		return status.Errorf(codes.Internal, e.Error())
	}

	mac := hmac.New(sha256.New, []byte(gCli.String("master-secret")))
	mac.Write([]byte(m.getId()))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(mmac, expectedMAC) {
		gLog.Info().Str("worker_id", m.getId()).Msg("worker authorization failed")
		return status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("worker_id", m.getId()).Msg("the worker's reply has been authorized")
	return
}

func (m *worker) getRPCErrors(err error) error {
	estatus, _ := status.FromError(err)

	switch estatus.Code() {
	case codes.OK:
		return nil

	case codes.Unavailable:
		gLog.Warn().Msg("trying to reconnect to the worker service...")
		if e := m.reconnect(); e != nil {
			gLog.Error().Err(e).Msg("could not reconnect to the worker service")
		}

	default:
		gLog.Warn().Str("error_code", estatus.Code().String()).Str("error_message", estatus.Message()).
			Msg("abnormal response from worker service")
	}

	return err
}

// methods

func (m *worker) getInitialServiceData() (_ string, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	var md metadata.MD
	var rpl *pb.InitReply
	if rpl, e = m.gservice.Init(ctx, &emptypb.Empty{}, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	if m.id != rpl.GetWorkerId() {
		gLog.Warn().Str("worker_id", m.id).Msg("abnormal reply from init method of the worker service; drop worker")
		return m.id, errWorkerInvalidId
	}

	// unpack Torrent
	var buf []byte
	if buf, e = json.Marshal(rpl.GetTorrent()); e != nil {
		return
	}

	if e = json.Unmarshal(buf, &m.trrs); e != nil {
		return
	}

	m.wdFreeSpace = rpl.GetWDFreeSpace()
	m.version = rpl.GetWorkerVersion()

	return
}

// get torrents from workers without any caches
func (m *worker) getTorrents() (trrs []*deluge.Torrent, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	var md metadata.MD
	var rpl *pb.TorrentsReply
	if rpl, e = m.gservice.GetTorrents(ctx, &emptypb.Empty{}, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	var buf []byte
	if buf, e = json.Marshal(rpl.GetTorrent()); e != nil {
		return
	}

	if e = json.Unmarshal(buf, &trrs); e != nil {
		return
	}

	m.trrs = trrs

	gLog.Debug().Int("torrents_count", len(trrs)).Msg("got reply from the worker with torrents list")
	return
}

func (m *worker) getFreeSpace() (_ uint64, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	var md metadata.MD
	var rpl *pb.SystemSpaceReply
	if rpl, e = m.gservice.GetSystemFreeSpace(ctx, &emptypb.Empty{}, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	gLog.Debug().Uint64("worker_fspace", rpl.FreeSpace).Msg("got reply from the worker with system free space")
	return rpl.FreeSpace, e
}

func (m *worker) saveTorrentFile(fname string, fbytes *[]byte) (_ int64, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	req := &pb.TFileSaveRequest{
		Filename: fname,
		Payload:  *fbytes,
	}

	var md metadata.MD
	var rpl *pb.TFileSaveReply
	if rpl, e = m.gservice.SaveTorrentFile(ctx, req, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	gLog.Debug().Int64("written_bytes", rpl.WrittenBytes).Msg("got reply from the worker with written bytes")
	return rpl.WrittenBytes, e
}

func (m *worker) deleteTorrent(hash, name string, withData bool) (_ uint64, _ uint64, e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	req := &pb.TorrentDropRequest{
		Name:     name,
		Hash:     hash,
		WithData: withData,
	}

	var md metadata.MD
	var rpl *pb.TorrentDropReply
	if rpl, e = m.gservice.DropTorrent(ctx, req, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	gLog.Debug().Uint64("worker_fspace", rpl.GetFreeSpace()).Uint64("worker_freed_space", rpl.FreedSpace).
		Msg("got reply from the worker with deleted bytes")

	return rpl.FreedSpace, rpl.FreeSpace, e
}

func (m *worker) forceReannounce() (e error) {
	ctx, cancel := m.newServiceRequest(gCli.Duration("grpc-request-timeout"))
	defer cancel()

	var md metadata.MD
	if _, e = m.gservice.ForceReannounce(ctx, &emptypb.Empty{}, grpc.Header(&md)); m.getRPCErrors(e) != nil {
		return
	}

	if e = m.authorizeSerivceReply(&md); e != nil {
		return
	}

	return
}
