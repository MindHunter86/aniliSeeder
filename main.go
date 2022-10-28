package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	application "github.com/MindHunter86/aniliSeeder/app"
	appcli "github.com/MindHunter86/aniliSeeder/cli"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var version = "devel" // -ldflags="-X 'main.version=X.X.X'"

func main() {
	// debug
	// pprof.WriteHeapProfile("mem.pprof")
	// defer profile.Start(profile.MemProfileHeap, profile.ProfilePath(".")).Stop()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	f, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	defer f.Close()

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
		// common flags
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

		// http client settings
		&cli.BoolFlag{
			Name:  "http-client-insecure",
			Usage: "Flag for TLS certificate verification disabling",
		},
		&cli.DurationFlag{
			Name:  "http-client-timeout",
			Usage: "Internal HTTP client connection `TIMEOUT` (format: 1000ms, 1s)",
			Value: 3 * time.Second,
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

		// syslog settings
		&cli.StringFlag{
			Name:  "syslog-addr",
			Usage: "",
			Value: "10.10.11.1:33517",
		},
		&cli.StringFlag{
			Name:  "syslog-proto",
			Usage: "",
			Value: "tcp",
		},
		&cli.StringFlag{
			Name:  "syslog-tag",
			Usage: "",
			Value: "aniliseeder",
		},

		// swarm settings
		&cli.BoolFlag{
			Name:    "is-master",
			Usage:   "",
			EnvVars: []string{"IS_MASTER"},
		},
		&cli.StringFlag{
			Name:  "master-addr",
			Usage: "",
			Value: "localhost:8081",
		},
		&cli.StringFlag{
			Name:    "master-secret",
			Usage:   "",
			Value:   "randomsecretkey",
			EnvVars: []string{"SWARM_MASTER_SECRETKEY"},
		},
		&cli.DurationFlag{
			Name:  "master-mon-interval",
			Usage: "master workers monitoring checks interval; 0 - for disabling",
			Value: 3 * time.Second,
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
			Name:  "grpc-request-timeout",
			Usage: "",
			Value: time.Second,
		},
		&cli.BoolFlag{
			Name:  "grpc-disable-reconnect",
			Usage: "",
		},
		&cli.DurationFlag{
			Name:  "grpc-ping-reconnect-hold",
			Usage: "time for grpc reconnection process",
			Value: 5 * time.Second,
		},
		&cli.IntFlag{
			Name:  "grpc-reconnect-tries",
			Usage: "",
			Value: 10,
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

		// anilibria settings
		&cli.StringFlag{
			Name:  "anilibria-baseurl",
			Usage: "",
			Value: "https://www.anilibria.tv",
		},
		&cli.StringFlag{
			Name:  "anilibria-api-baseurl",
			Usage: "",
			Value: "https://api.anilibria.tv/v2",
		},
		&cli.StringFlag{
			Name:    "anilibria-login-username",
			Usage:   "login",
			EnvVars: []string{"ANILIBRIA_USERNAME"},
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
			EnvVars: []string{"DELUGE_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "deluge-password",
			Usage:   "",
			Value:   "",
			EnvVars: []string{"DELUGE_PASSWORD"},
		},
		&cli.StringFlag{
			Name:  "deluge-data-path",
			Usage: "directory for space monitoring",
			Value: "./data",
		},
		&cli.StringFlag{
			Name:  "deluge-torrents-path",
			Usage: "download directory for .torrent files",
			Value: "./data",
		},
		&cli.Uint64Flag{
			Name:  "deluge-disk-minimal",
			Usage: "in MB; ",
			Value: 128,
		},

		// legacy settings
		&cli.IntFlag{
			Name:  "cmd-vkscore-warn",
			Usage: "all torrents below this value will be marked as inefficient",
			Value: 25,
		},

		// deploy settings
		&cli.BoolFlag{
			Name:  "deploy-ignore-errors",
			Usage: "",
		},

		// cron settings
		&cli.BoolFlag{
			Name:  "cron-disable",
			Usage: "",
		},

		// master cli settings
		&cli.StringFlag{
			Name:  "socket-path",
			Usage: "",
			Value: "aniliSeeder.sock",
		},
	}

	app.Action = func(c *cli.Context) error {
		// log.Debug().Msg("ready...")
		// log.Debug().Strs("args", os.Args).Msg("")

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
			Action: func(c *cli.Context) (e error) {
				if c.String("syslog-addr") != "" {
					if runtime.GOOS == "windows" {
						log.Error().Msg("sorry, but syslog is not worked for windows; golang does not support syslog for win systems")
						return os.ErrProcessDone
					}

					log.Debug().Msg("connecting to the syslog server...")

					var nlog *zerolog.Logger
					if nlog, e = utils.SetUpSyslogWriter(c); e != nil {
						return
					}

					log.Info().Msg("connection to the syslog server has been established; reset log driver ...")

					log = *nlog
					log = log.Hook(SeverityHook{})
					log.Info().Msg("zerolog has been reinited; starting application ...")
				}

				a := application.NewApp(c, &log)
				return a.Bootstrap()
			},
		},
		&cli.Command{
			Name:  "cli",
			Usage: "",
			Action: func(c *cli.Context) error {
				return appcli.TestDial(c, "")
			},
		},
		&cli.Command{
			Name:  "test",
			Usage: "",
			Action: func(c *cli.Context) error {
				aniApi, e := anilibria.NewApiClient(c, &log)
				if e != nil {
					return e
				}

				titles, e := aniApi.SearchTitlesByName("Urusei Yatsura 2022")
				for _, title := range titles {
					log.Debug().Str("title_name", title.Names.Ru).Msg("")
				}

				return e
			},
		},
	}

	// sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if e := app.Run(os.Args); e != nil {
		log.Fatal().Err(e).Msg("")
	}
}

type SeverityHook struct{}

func (SeverityHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level != zerolog.DebugLevel && version != "devel" {
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
