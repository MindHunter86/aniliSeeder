package p2p

import (
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/anacrolix/torrent"
)

var (
	gLogger *zerolog.Logger
)

type P2PClient struct {
	config *torrent.ClientConfig
}

func NewP2PClient(l *zerolog.Logger) *P2PClient {
	gLogger = l
	return &P2PClient{}
}

func (m *P2PClient) Bootstrap() error {
	config := torrent.NewDefaultClientConfig()
	config.SetListenAddr("0.0.0.0:59149")
	config.ListenPort = 59149
	config.Debug = true
	config.AlwaysWantConns = true
	config.EstablishedConnsPerTorrent = 1000
	config.HalfOpenConnsPerTorrent = 100
	config.AcceptPeerConnections = true
	config.DisableAggressiveUpload = false
	config.DisableAcceptRateLimiting = true
	config.KeepAliveTimeout = 600 * time.Second
	config.TotalHalfOpenConns = 250
	config.PublicIp4 = net.ParseIP("212.30.191.141")
	config.UpnpID = "golang"

	config.DisableIPv6 = true
	config.DataDir = "/media/vkomissarov/VKHDDGames/vkomissarov/deluge"
	config.DisableWebseeds = true
	config.DisableWebtorrent = true
	config.Seed = true
	config.NoUpload = false

	client, err := torrent.NewClient(config)
	if err != nil {
		return err
	}
	defer client.Close()

	tFiles, err := ioutil.ReadDir("/media/vkomissarov/VKHDDGames/vkomissarov/deluge-files")
	if err != nil {
		return err
	}

	var tChan = make(chan *torrent.Torrent, 128)

	for _, tFile := range tFiles {
		if !tFile.IsDir() {
			tr, e := client.AddTorrentFromFile("/media/vkomissarov/VKHDDGames/vkomissarov/deluge-files/" + tFile.Name())
			if e != nil {
				return e
			}

			<-tr.GotInfo()
			tr.AllowDataUpload()
			tChan <- tr
		}
	}

	// tr.DownloadAll()
	client.WaitAll()
	log.Println("Done!")
	time.Sleep(6000000 * time.Second)

	return nil
}
