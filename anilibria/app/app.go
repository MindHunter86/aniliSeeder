package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type App struct{}

var (
	gCli *cli.Context
	gLog *zerolog.Logger

	gCtx   context.Context
	gAbort context.CancelFunc
)

func NewApp(c *cli.Context, l *zerolog.Logger) *App {
	gCli, gLog = c, l

	return &App{}
}

func (m *App) Bootstrap() error {
	kernSignal := make(chan os.Signal, 1)
	signal.Notify(kernSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)

	gCtx, gAbort = context.WithCancel(context.WithValue(context.Background(), contextKeyKernSignal, kernSignal))

	var wg = sync.WaitGroup{}
	defer wg.Wait()

	// main event loop
	wg.Add(1)
	go m.loop(wg.Done)

	// another subsystems
	// ...

	return nil
}

func (m *App) loop(done func()) {
	defer done()

	kernSignal := gCtx.Value(contextKeyKernSignal).(chan os.Signal)

	gLog.Debug().Msg("initiate main event loop")
	defer gLog.Debug().Msg("initiate main event loop")

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

func (m *App) Close() error {
	return nil
}
