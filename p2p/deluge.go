package p2p

import (
	"encoding/json"
	"net"
	"os"

	delugeclient "github.com/MindHunter86/go-libdeluge"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
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

func (m *Client) GetTorrentsStatus() (e error) {
	var trrs map[string]*delugeclient.TorrentStatus
	if trrs, e = m.deluge.TorrentsStatus(delugeclient.StateUnspecified, nil); e != nil {
		return
	}

	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "  ")

	if e = je.Encode(trrs); e != nil {
		return
	}

	return
}
