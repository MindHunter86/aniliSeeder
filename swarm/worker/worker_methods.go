package worker

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"strings"

	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"golang.org/x/net/context"
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
		return "", status.Errorf(codes.DataLoss, "")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "")
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

	mac := hmac.New(sha256.New, []byte(gCli.String("swarm-master-secret")))
	mac.Write([]byte(id[0]))
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal([]byte(ah[0]), expectedMAC) {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("master_id", id[0]).Msg("the master's request has been authorized")
	return id[0], nil
}

func (m *WorkerService) Init(ctx context.Context, _ *emptypb.Empty) (*pb.InitReply, error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	var trrs []*structpb.Struct
	if trrs, e = m.w.getTorrents(); e != nil {
		return nil, status.Errorf(codes.Internal, e.Error())
	}

	return &pb.InitReply{
		WorkerId:      m.w.id,
		WorkerVersion: gCli.App.Version,
		WDFreeSpace:   utils.CheckDirectoryFreeSpace("deluge-data-path"),
		Torrent:       trrs,
	}, e
}

func (m *WorkerService) Ping(ctx context.Context, _ *emptypb.Empty) (_ *emptypb.Empty, _ error) {
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

	if req.GetName() != trr.Name {
		return nil, status.Errorf(codes.InvalidArgument, "given name is not equal torrent name")
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

	if req.GetName() != trr.Name {
		return nil, status.Errorf(codes.InvalidArgument, "given name is not equal torrent name")
	}

	return &pb.TorrentDropReply{
		FreedSpace: uint64(trr.TotalSize),
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

	return &pb.SystemSpaceReply{
		FreeSpace: utils.CheckDirectoryFreeSpace(gCli.String("torrentfiles-dir")),
	}, e
}
