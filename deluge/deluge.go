package deluge

import (
	"net"

	delugeclient "github.com/MindHunter86/go-libdeluge"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type TrackerStatus uint8

const (
	TrackerStatusOK TrackerStatus = iota
	TrackerStatusUnknown
	TrackerStatusTimedOut
	TrackerStatusSent
	TrackerStatusNotRegistered
)

type Client struct {
	deluge *delugeclient.ClientV2
}

var (
	gCli *cli.Context
	gLog *zerolog.Logger
)

func NewClient(c *cli.Context, log *zerolog.Logger) (*Client, error) {
	gCli, gLog = c, log

	var e error
	var addr *net.TCPAddr

	if addr, e = net.ResolveTCPAddr("tcp", gCli.String("deluge-addr")); e != nil {
		return nil, e
	}

	deluge := delugeclient.NewV2(delugeclient.Settings{
		Hostname: addr.IP.String(),
		Port:     uint(addr.Port),
		Login:    gCli.String("deluge-username"),
		Password: gCli.String("deluge-password"),
	})

	if e = deluge.Connect(); e != nil {
		return nil, e
	}

	return &Client{
		deluge: deluge,
	}, e
}
