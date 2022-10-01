package main

import (
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var log zerolog.Logger
var version = "devel" // -ldflags="-X 'main.version=X.X.X'"

func main() {
	app := cli.NewApp()
	cli.VersionFlag = &cli.BoolFlag{Name: "print-version", Aliases: []string{"V"}}

	app.Name = "aniliSeeder"
	app.Version = version
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "MindHunter86",
			Email: "admin@vkom.cc",
		},
	}
	app.Copyright = "(c) 2022 mindhunter86"
	app.Usage = "N\\A"

	app.Flags = []cli.Flag{
		// common settings
		// cli.DurationFlag{
		// 	Name:  "http-client-timeout",
		// 	Usage: "Internal HTTP client connection `TIMEOUT` (format: 1000ms, 1s)",
		// 	Value: 10 * time.Second,
		// },
		// cli.BoolFlag{
		// 	Name:  "http-client-insecure",
		// 	Usage: "Flag for TLS certificate verification disabling",
		// },

		&cli.IntFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   5,
			Usage:   "Verbose `LEVEL` (value from 5(debug) to 0(panic) and -1 for log disabling(quite mode))",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "",
		},
		&cli.BoolFlag{
			Name:    "quite",
			Aliases: []string{"q"},
			Usage:   "Flag is equivalent to verbose -1",
		},

		// queue settings
		// application settings

		// billmanager opts
		&cli.StringFlag{
			Name: "command",
		},
		&cli.StringFlag{
			Name: "subcommand",
		},
	}

	app.Action = DefaultAction("")

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if e := app.Run(os.Args); e != nil {
		log.Fatal().Err(e).Msg("")
	}
}

func DefaultAction(name string) cli.ActionFunc {
	return func(c *cli.Context) (e error) {
		log := zerolog.New(zerolog.ConsoleWriter{
			Out: os.Stderr,
		}).With().Timestamp().Logger().Hook(SeverityHook{})
		zerolog.TimeFieldFormat = time.RFC3339Nano
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		log.Debug().Msg("ready...")

		// TODO
		// if c.Int("verbose") < -1 || c.Int("verbose") > 5 {
		// 	log.Fatal().Msg("There is invalid data in verbose option. Option supports values for -1 to 5")
		// }

		// zerolog.SetGlobalLevel(zerolog.Level(int8((c.Int("verbose") - 5) * -1)))
		// if c.Int("verbose") == -1 || c.Bool("quite") {
		// 	zerolog.SetGlobalLevel(zerolog.Disabled)
		// }

		log.Debug().Strs("args", os.Args).Msg("")
		return nil
	}
}

type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.DebugLevel {
		return
	}

	rfn := "unknown"
	pcs := make([]uintptr, 1)

	if runtime.Callers(4, pcs) != 0 {
		if fun := runtime.FuncForPC(pcs[0] - 1); fun != nil {
			rfn = fun.Name()
		}
	}

	fn := strings.Split(rfn, "/")
	e.Str("func", fn[len(fn)-1:][0])
}
