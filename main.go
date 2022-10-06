package main

import (
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	application "github.com/MindHunter86/aniliSeeder/app"
	appcli "github.com/MindHunter86/aniliSeeder/cli"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var version = "devel" // -ldflags="-X 'main.version=X.X.X'"

func main() {
	// logger
	log := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger().Hook(SeverityHook{})
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// application
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
		&cli.DurationFlag{
			Name:  "http-client-timeout",
			Usage: "Internal HTTP client connection `TIMEOUT` (format: 1000ms, 1s)",
			Value: 3 * time.Second,
		},
		&cli.BoolFlag{
			Name:  "http-client-insecure",
			Usage: "Flag for TLS certificate verification disabling",
		},
		&cli.DurationFlag{
			Name:  "http-tcp-timeout",
			Usage: "",
			Value: 1 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-tls-handshake-timeout",
			Usage: "",
			Value: 1 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-idle-timeout",
			Usage: "",
			Value: 300 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http-keepalive-timeout",
			Usage: "",
			Value: 300 * time.Second,
		},
		&cli.IntFlag{
			Name:  "http-max-idle-conns",
			Usage: "",
			Value: 100,
		},
		&cli.BoolFlag{
			Name:  "http-debug",
			Usage: "",
		},

		&cli.StringFlag{
			Name:  "socket-path",
			Usage: "",
			Value: "aniliSeeder.sock",
		},

		&cli.IntFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   5,
			Usage:   "Verbose `LEVEL` (value from 5(debug) to 0(panic) and -1 for log disabling(quite mode))",
		},
		&cli.BoolFlag{
			Name:    "quite",
			Aliases: []string{"q"},
			Usage:   "Flag is equivalent to verbose -1",
		},

		// swarm settings
		&cli.BoolFlag{
			Name:  "swarm-is-master",
			Usage: "",
		},
		&cli.StringFlag{
			Name:  "swarm-master-listen",
			Usage: "",
			Value: "localhost:8081",
		},
		&cli.StringFlag{
			Name:  "swarm-master-addr",
			Usage: "",
			Value: "localhost:8081",
		},
		&cli.StringFlag{
			Name:  "swarm-custom-ca-path",
			Usage: "",
			Value: "",
		},
		&cli.StringFlag{
			Name:    "swarm-master-secret",
			Usage:   "",
			Value:   "randomsecretkey",
			EnvVars: []string{"SWARM_MASTER_SECRETKEY"},
		},

		// gRPC settings
		&cli.BoolFlag{
			Name:  "grpc-insecure",
			Usage: "",
		},
		&cli.DurationFlag{
			Name:  "grpc-connect-timeout",
			Usage: "for worker",
			Value: 3 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "grpc-ping-interval",
			Usage: "0 for disabling",
			Value: time.Second,
		},
		&cli.DurationFlag{
			Name:  "grpc-ping-timeout",
			Usage: "",
			Value: 300 * time.Millisecond,
		},

		// http2 settings
		&cli.DurationFlag{
			Name:  "http2-ping-time",
			Usage: "for worker",
			Value: 3 * time.Second,
		},
		&cli.DurationFlag{
			Name:  "http2-ping-timeout",
			Usage: "for worker",
			Value: time.Second,
		},
		&cli.DurationFlag{
			Name:  "http2-conn-max-age",
			Usage: "for master; 0 for disable",
			Value: 600 * time.Second,
		},

		// queue settings
		// ...

		// anilibria settings
		&cli.StringFlag{
			Name:  "anilibria-api-baseurl",
			Usage: "",
			Value: "https://api.anilibria.tv/v2",
		},
		&cli.StringFlag{
			Name:  "anilibria-baseurl",
			Usage: "",
			Value: "https://www.anilibria.tv",
		},
		&cli.StringFlag{
			Name:    "anilibria-login-username",
			Usage:   "login",
			EnvVars: []string{"ANILIBRIA_LOGIN", "ANILIBRIA_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "anilibria-login-password",
			Usage:   "password",
			EnvVars: []string{"ANILIBRIA_PASSWORD"},
		},

		// deluge settings
		&cli.StringFlag{
			Name:  "deluge-addr",
			Usage: "",
			Value: "127.0.0.1:58846",
		},
		&cli.StringFlag{
			Name:    "deluge-username",
			Usage:   "",
			Value:   "localclient",
			EnvVars: []string{"DELUGE_LOGIN", "DELUGE_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "deluge-password",
			Usage:   "",
			Value:   "",
			EnvVars: []string{"DELUGE_PASSWORD"},
		},

		// common settings
		&cli.StringFlag{
			Name:  "torrentfiles-dir",
			Usage: "",
			Value: "./data",
		},
		&cli.IntFlag{
			Name:  "torrents-vkscore-line",
			Usage: "",
			Value: 25,
		},
		&cli.UintFlag{
			Name:  "disk-minimal-avaliable",
			Usage: "In MB",
			Value: 128,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Debug().Msg("ready...")
		log.Debug().Strs("args", os.Args).Msg("")

		// TODO
		// if c.Int("verbose") < -1 || c.Int("verbose") > 5 {
		// 	log.Fatal().Msg("There is invalid data in verbose option. Option supports values for -1 to 5")
		// }

		// zerolog.SetGlobalLevel(zerolog.Level(int8((c.Int("verbose") - 5) * -1)))
		// if c.Int("verbose") == -1 || c.Bool("quite") {
		// 	zerolog.SetGlobalLevel(zerolog.Disabled)
		// }

		// return p2p.NewP2PClient(&log).Bootstrap()

		// ====================

		// api, err := anilibria.NewApiClient(c, &log)
		// if err != nil {
		// 	return err
		// }

		// if _, err = api.GetTitleSchedule(); err != nil {
		// 	return err
		// }

		// dClient, err := deluge.NewClient(c, &log)
		// if err != nil {
		// 	return err
		// }

		// return dClient.GetTorrentsStatus()
		// ===========

		return os.ErrInvalid
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "serve",
			Usage: "",
			Action: func(c *cli.Context) error {
				a := application.NewApp(c, &log)
				return a.Bootstrap()
			},
		},
		&cli.Command{
			Name:  "test",
			Usage: "",
			Action: func(c *cli.Context) error {
				return appcli.TestDial(c, "")
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if e := app.Run(os.Args); e != nil {
		log.Fatal().Err(e).Msg("")
	}
}

type SeverityHook struct{}

func (SeverityHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
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
