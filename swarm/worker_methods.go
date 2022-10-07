package swarm

import (
	"strings"

	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*Worker) authorizeMasterRequest(ctx context.Context) (string, error) {
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
		Str("master_ua", md.Get("user-agent")[0]).Msg("master connect accepted, authorizing...")

	ak := md.Get("x-access-token")
	if len(ak) != 1 {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if strings.TrimSpace(ak[0]) == "" {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.InvalidArgument, "")
	}
	if ak[0] != gCli.String("swarm-master-secret") {
		gLog.Info().Str("master_id", id[0]).Msg("master authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("master_id", md.Get("x-master-id")[0]).Msg("the master's connect has been authorized")
	return id[0], nil
}

func (m *Worker) GetTorrents(ctx context.Context, _ *emptypb.Empty) (_ *pb.TorrentsReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")
	return
}

func (m *Worker) GetTorrentScore(ctx context.Context, req *pb.TorrentScoreReply) (_ *pb.TorrentScoreReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")
	return
}

func (m *Worker) DropTorrent(ctx context.Context, req *pb.TorrentDropRequest) (_ *pb.TorrentDropReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")
	return
}

func (m *Worker) SaveTorrentFile(ctx context.Context, req *pb.TFileSaveRequest) (_ *pb.TFileSaveReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")
	return
}

func (m *Worker) GetSystemFreeSpace(ctx context.Context, _ *emptypb.Empty) (_ *pb.SystemSpaceReply, _ error) {
	mid, e := m.authorizeMasterRequest(ctx)
	if e != nil {
		return nil, e
	}

	gLog.Debug().Str("master_id", mid).Msg("processing master request...")

	return &pb.SystemSpaceReply{
		FreeSpace: utils.CheckDirectoryFreeSpace(gCli.String("torrentfiles-dir")),
	}, e
}
