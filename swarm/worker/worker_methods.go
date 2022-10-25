package worker

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type WorkerService struct {
	pb.UnimplementedWorkerServiceServer

	w *Worker
}

func NewWorkerService(w *Worker) *WorkerService {
	return &WorkerService{w: w}
}

func (*WorkerService) authorizeMasterRequest(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "could not get peer data")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "could not get metadata")
	}

	id := md.Get("x-master-id")
	if len(id) != 1 {
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(id[0]) == "" {
		return "", status.Errorf(codes.InvalidArgument, "")
	}

	gLog.Debug().Str("master_ip", p.Addr.String()).Str("master_id", id[0]).
		Str("master_ua", md.Get("user-agent")[0]).Msg("master request accepted, authorizing...")

	ah := md.Get("x-authentication-hash")
	if len(ah) != 1 {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(ah[0]) == "" {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}

	wmac, e := hex.DecodeString(ah[0])
	if e != nil {
		return "", status.Errorf(codes.Internal, e.Error())
	}

	mac := hmac.New(sha256.New, []byte(gCli.String("swarm-master-secret")))
	mac.Write([]byte(id[0]))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(wmac, expectedMAC) {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("master_id", id[0]).Msg("the master's request has been authorized")
	return id[0], nil
}

func (m *WorkerService) authorizeServiceReply(ctx context.Context) error {
	mac := hmac.New(sha256.New, []byte(gCli.String("swarm-master-secret")))
	io.WriteString(mac, m.w.id)

	md := metadata.New(map[string]string{
		"x-worker-id":           m.w.id,
		"x-authentication-hash": hex.EncodeToString(mac.Sum(nil)),
	})

	return grpc.SendHeader(ctx, md)
}

func (m *WorkerService) Init(ctx context.Context, _ *emptypb.Empty) (*pb.InitReply, error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		// !!!
		gLog.Warn().Msg("aborting application because of inital phase is failed")
		gAbort()
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	var trrs []*structpb.Struct
	if trrs, e = m.w.getTorrents(); e != nil {
		// TODO: may be remove it...
		gLog.Warn().Msg("aborting application because of inital phase is failed")
		gAbort()

		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.InitReply{
		WorkerId:      m.w.id,
		WorkerVersion: gCli.App.Version,
		WDFreeSpace:   utils.CheckDirectoryFreeSpace(gCli.String("deluge-data-path")),
		Torrent:       trrs,
	}, e
}

func (*WorkerService) Ping(_ context.Context, _ *emptypb.Empty) (_ *emptypb.Empty, _ error) {
	// wid, e := m.authorizeWorker(ctx)
	// if e != nil {
	// 	return &emptypb.Empty{}, e
	// }

	// if !m.isWorkerRegistered(wid) {
	// 	gLog.Info().Str("worker_id", wid).Msg("worker is not registered, returning 403...")
	// 	return nil, status.Errorf(codes.PermissionDenied, "")
	// }

	// gLog.Info().Str("worker_id", wid).Msg("received ping from worker")
	return &emptypb.Empty{}, nil
}

func (m *WorkerService) GetTorrents(ctx context.Context, _ *emptypb.Empty) (_ *pb.TorrentsReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	var trrs []*structpb.Struct
	if trrs, e = m.w.getTorrents(); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.TorrentsReply{
		Torrent: trrs,
	}, e
}

func (m *WorkerService) GetTorrentScore(ctx context.Context, req *pb.TorrentScoreRequest) (_ *pb.TorrentScoreReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	if req.GetHash() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect hash")
	}

	var trr *deluge.Torrent
	if trr, e = gDeluge.TorrentStatus(req.GetHash()); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if req.GetName() != trr.GetName() {
		return nil, status.Errorf(codes.InvalidArgument, "given name is not equal torrent name")
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.TorrentScoreReply{
		Score: trr.GetVKScore(),
		Ratio: trr.Ratio,
	}, e
}

func (m *WorkerService) DropTorrent(ctx context.Context, req *pb.TorrentDropRequest) (_ *pb.TorrentDropReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	if req.GetHash() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect hash")
	}

	var trr *deluge.Torrent
	if trr, e = gDeluge.TorrentStatus(req.GetHash()); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if req.GetName() != trr.GetName() {
		return nil, status.Errorf(codes.InvalidArgument, "given name is not equal torrent name")
	}

	var fspace uint64
	if req.GetWithData() {
		fspace = uint64(trr.TotalSize)
		gLog.Warn().Str("master_id", mid).Msg("torrent removing with data detected")
	}

	var ok bool
	if ok, e = gDeluge.RemoveTorrent(trr.Hash, req.GetWithData()); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if !ok {
		return nil, status.Errorf(codes.Internal, "undefined internal error in torrent removing; ok is not true")
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.TorrentDropReply{
		FreedSpace: fspace,
		FreeSpace:  utils.CheckDirectoryFreeSpace(gCli.String("torrentfiles-dir")),
	}, e
}

func (m *WorkerService) SaveTorrentFile(ctx context.Context, req *pb.TFileSaveRequest) (_ *pb.TFileSaveReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	buf := bytes.NewBuffer(req.GetPayload())
	n, e := gDeluge.SaveTorrentFile(req.GetFilename(), buf)
	if e != nil {
		return nil, status.Errorf(codes.FailedPrecondition, e.Error())
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	gLog.Debug().Str("master_id", mid).Int64("written_bytes", n).Msg("the requested method has been processed")
	return &pb.TFileSaveReply{
		WrittenBytes: n,
	}, e
}

func (m *WorkerService) GetSystemFreeSpace(ctx context.Context, _ *emptypb.Empty) (_ *pb.SystemSpaceReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.SystemSpaceReply{
		FreeSpace: utils.CheckDirectoryFreeSpace(gCli.String("torrentfiles-dir")),
	}, e
}

func (m *WorkerService) ForceReannounce(ctx context.Context, _ *emptypb.Empty) (_ *emptypb.Empty, e error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	var thashes []string
	if thashes, e = gDeluge.GetTorrentsHashes(); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if e = gDeluge.ForceReannounce(thashes...); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	if e = m.authorizeServiceReply(ctx); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &emptypb.Empty{}, e
}
