package master

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/MindHunter86/aniliSeeder/deluge"
	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/hashicorp/yamux"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	gCli *cli.Context
	gLog *zerolog.Logger
	gCtx context.Context

	gMasterId string
)

// gRPC bypass picture
// https://camo.githubusercontent.com/79dbaaa1fb7d239f1d21d4be23985b831babfc4b95538413298dce1c5c2600e1/68747470733a2f2f6c68332e676f6f676c6575736572636f6e74656e742e636f6d2f2d544c51596b4975396a49412f57664c61644c4a357a55492f41414141414141415967382f3745716c72613754764b38676f736b44465142637777394c6e584369794a387977434c63424741732f73313630302f494d475f383032372e6a7067

type Master struct {
	rawListener net.Listener
	workerPool  *workerPool
}

func NewMaster(ctx context.Context) *Master {
	gCtx = ctx
	gLog = gCtx.Value(utils.ContextKeyLogger).(*zerolog.Logger)
	gCli = gCtx.Value(utils.ContextKeyCliContext).(*cli.Context)
	gMasterId = uuid.NewV4().String()

	return &Master{
		workerPool: newWorkerPool(),
	}
}

func (m *Master) handleIncomingConnection(conn net.Conn) (e error) {
	gLog.Debug().Str("master_listen", gCli.String("master-addr")).
		Msg("trying to initialize mux session...")

	var muxsess *yamux.Session
	if muxsess, e = yamux.Client(conn, yamux.DefaultConfig()); e != nil {
		return
	}

	var d time.Duration
	if d, e = muxsess.Ping(); e != nil {
		return
	}

	gLog.Debug().Str("ping_time", d.String()).Msg("mux session is alive")

	gLog.Debug().Str("master_listen", gCli.String("master-addr")).
		Msg("trying to initialize gRPC client...")

	if _, e = m.workerPool.newWorker(muxsess); e != nil {
		gLog.Debug().Err(e).Msg("got an error while processing new worker; drop mux session ...")
		muxsess.Close()
		return
	}

	return
}

func (m *Master) Bootstrap() (e error) {
	gLog.Debug().Str("master_listen", gCli.String("master-addr")).
		Msg("initializing the tcp server for further muxing")

	if m.rawListener, e = net.Listen("tcp", gCli.String("master-addr")); e != nil {
		return
	}

	return m.run()
}

func (m *Master) run() (e error) {
	gLog.Debug().Str("master_listen", gCli.String("master-addr")).
		Msg("initializing net acceptor; starting listening for incoming TCP connections...")

	var wg sync.WaitGroup

	var clock sync.RWMutex
	var conns []net.Conn

LOOP:
	for {
		select {
		case <-gCtx.Done():
			gLog.Warn().Msg("context done() has been caught; closing grpc server socket...")
			break LOOP
		default:
			conn, e := m.rawListener.Accept()
			if e != nil {
				gLog.Error().Err(e).Msg("got some error with processing a new tcp client")
			}

			gLog.Debug().Str("master_listen", gCli.String("master-addr")).Str("client_addr", conn.RemoteAddr().String()).
				Msg("new incoming connection; processing...")

			go func(cn net.Conn) {
				if err := m.handleIncomingConnection(cn); e != nil {
					gLog.Warn().Err(err).Msg("got error while handling the workers connection")
					cn.Close()
				}
			}(conn)
		}
	}

	clock.Lock()
	defer clock.Unlock()

	for _, conn := range conns {
		gLog.Debug().Str("client_addr", conn.RemoteAddr().String()).Msg("trying to close the accepted client connection")
		if e = conn.Close(); e != nil {
			gLog.Warn().Str("client_addr", conn.RemoteAddr().String()).Err(e).
				Msg("got some errors while closing client connection")
		}
	}

	gLog.Debug().Msg("waiting for closing all accepted connection...")
	wg.Wait()

	return
}

func (*Master) authorizeWorker(ctx context.Context) (string, error) {
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
	if ak[0] != gCli.String("master-secret") {
		gLog.Info().Str("worker_id", id[0]).Msg("worker authorization failed")
		return "", status.Errorf(codes.Unauthenticated, "")
	}

	gLog.Debug().Str("worker_id", md.Get("x-worker-id")[0]).Msg("the worker's connect has been authorized")
	return id[0], nil
}

func (m *Master) Register(ctx context.Context, req *pb.RegistrationRequest) (_ *emptypb.Empty, e error) {
	var wid string
	if wid, e = m.authorizeWorker(ctx); e != nil {
		return
	}

	gLog.Info().Str("worker_id", wid).Msg("new client validation phase running...")

	switch {
	// case m.workers[wid] != nil:
	// 	return nil, status.Errorf(codes.AlreadyExists, "")
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

	// m.Lock()
	// m.workers[wid] = &worker.Worker{
	// 	Version:     req.WorkerVersion,
	// 	WDFreeSpace: req.WDFreeSpace,
	// 	Torrents:    wtrrs,
	// }
	// m.Unlock()

	gLog.Info().Str("worker_id", wid).Msg("new client registration has been completed")
	return &emptypb.Empty{}, e
}

func (*Master) IsMaster() bool {
	return true
}
