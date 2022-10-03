package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type App struct{}

var (
	gCli *cli.Context
	gLog *zerolog.Logger

	gCtx   context.Context
	gAbort context.CancelFunc

	gAniApi *anilibria.ApiClient
	gDeluge *deluge.Client
)

func NewApp(c *cli.Context, l *zerolog.Logger) *App {
	gCli, gLog = c, l

	return &App{}
}

func (m *App) Bootstrap() (e error) {
	kernSignal := make(chan os.Signal, 1)
	signal.Notify(kernSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)

	gCtx, gAbort = context.WithCancel(context.WithValue(context.Background(), contextKeyKernSignal, kernSignal))

	var wg = sync.WaitGroup{}
	defer wg.Wait()
	defer gLog.Debug().Msg("waiting for opened goroutines")
	defer gAbort()

	// main event loop
	wg.Add(1)
	go m.loop(wg.Done)

	// anilibria API
	if gAniApi, e = anilibria.NewApiClient(gCli, gLog); e != nil {
		return
	}

	// deluge RPC client
	if gDeluge, e = deluge.NewClient(gCli, gLog); e != nil {
		return
	}

	// another subsystems
	// ...

	// socket cmds server
	var sServer = NewSockServer()
	if e = sServer.Bootstrap(); e != nil {
		return
	}

	wg.Add(1)
	go sServer.Serve(wg.Done)

	wg.Wait()
	return
}

func (*App) loop(done func()) {
	defer done()

	kernSignal := gCtx.Value(contextKeyKernSignal).(chan os.Signal)

	gLog.Debug().Msg("initiate main event loop")
	defer gLog.Debug().Msg("main event loop has been closed")

LOOP:
	for {
		select {
		case <-kernSignal:
			gLog.Info().Msg("kernel signal has been caught; initiate application closing...")
			gAbort()
			break LOOP
		case <-gCtx.Done():
			break LOOP
		}
	}
}

// todo
// func (m *App) Close() error {
// 	return nil
// }
