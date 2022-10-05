package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/swarm"
	"github.com/MindHunter86/aniliSeeder/utils"
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
	var wg = sync.WaitGroup{}
	var echan = make(chan error, 32)

	gCtx, gAbort = context.WithCancel(context.Background())
	gCtx = context.WithValue(gCtx, utils.ContextKeyLogger, gLog)
	gCtx = context.WithValue(gCtx, utils.ContextKeyCliContext, gCli)

	defer m.checkErrorsBeforeClosing(echan)
	defer wg.Wait() // !!
	defer gLog.Debug().Msg("waiting for opened goroutines")
	defer gAbort()

	if gCli.Bool("swarm-is-master") {
		// anilibria API
		if gAniApi, e = anilibria.NewApiClient(gCli, gLog); e != nil {
			return
		}

		gCtx = context.WithValue(gCtx, utils.ContextKeyAnilibriaClient, gAniApi)
		// gRPC = swarm.NewMaster(gCli, gLog, gCtx)
	} else {
		// deluge RPC client
		if gDeluge, e = deluge.NewClient(gCli, gLog); e != nil {
			return
		}

		gCtx = context.WithValue(gCtx, utils.ContextKeyDelugeClient, gDeluge)
		gRPC = swarm.NewWorker(gCtx)
	}

	// grpc master/worker setup
	go func(errs chan error, done func()) {
		log.Println("1")
		if err := gRPC.Bootstrap(); err != nil {
			errs <- err
		}

		done()
	}(echan, wg.Done)

	// another subsystems
	// ...

	// main event loop
	wg.Add(1)
	go m.loop(echan, wg.Done)

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

	kernSignal := make(chan os.Signal, 1)
	signal.Notify(kernSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)

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
	if len(errs) == 0 {
		return
	}

	for err := range errs {
		gLog.Warn().Err(err).Msg("an error has been detected while application trying close the submodules")
	}
}

// todo
// func (m *App) Close() error {
// 	return nil
// }
