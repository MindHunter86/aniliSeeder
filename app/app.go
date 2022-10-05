package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/swarm"
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
	gRPC    swarm.Swarm
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
	var echan = make(chan error, 32)

	defer m.checkErrorsBeforeClosing(echan)
	defer wg.Wait()
	defer gLog.Debug().Msg("waiting for opened goroutines")
	defer gAbort()

	// main event loop
	wg.Add(1)
	go m.loop(echan, wg.Done)

	if gCli.Bool("swarm-is-master") {
		// deluge RPC client
		if gDeluge, e = deluge.NewClient(gCli, gLog); e != nil {
			return
		}

		gRPC = swarm.NewWorker(gCli, gLog, gCtx)
	} else {
		// anilibria API
		if gAniApi, e = anilibria.NewApiClient(gCli, gLog); e != nil {
			return
		}

		// gRPC = swarm.NewMaster(gCli, gLog, gCtx)
	}

	// grpc master/worker setup
	go func(errs chan error, done func()) {
		defer done()
		if err := gRPC.Bootstrap(); err != nil {
			errs <- err
		}
	}(echan, wg.Done)

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

func (*App) loop(errs chan error, done func()) {
	defer done()

	// ??
	// todo : review
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
		case err := <-errs:
			gLog.Error().Err(err).Msg("there are internal errors from one of application submodule")
			gLog.Warn().Msg("calling abort()...")
			gAbort()
		case <-gCtx.Done():
			gLog.Info().Msg("internal abort() has been caught; initiate application closing...")
			break LOOP
		}
	}
}

func (*App) checkErrorsBeforeClosing(errs chan error) {
	for err := range errs {
		gLog.Warn().Err(err).Msg("an error has been detected while application trying close the submodules")
	}
}

// todo
// func (m *App) Close() error {
// 	return nil
// }
