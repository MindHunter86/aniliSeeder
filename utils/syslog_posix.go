//go:build !windows

package utils

import (
	"log/syslog"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func SetUpSyslogWriter(c *cli.Context) (_ *zerolog.Logger, e error) {
	var slog *syslog.Writer
	if slog, e = syslog.Dial(c.String("syslog-proto"), c.String("syslog-addr"), syslog.LOG_INFO, c.String("syslog-tag")); e != nil {
		return nil, e
	}

	nlog := zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		slog,
	)).With().Timestamp().Logger()

	return &nlog, e
}
