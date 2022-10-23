package swarm

import (
	"github.com/MindHunter86/aniliSeeder/deluge"
)

type SwarmWorker struct {
	Id        string
	Version   string
	FreeSpace uint64

	ActiveTorrents []*deluge.Torrent
}

type Swarm interface {
	// common methods
	IsMaster() bool
	Bootstrap() error

	// master methods
	GetConnectedWorkers() map[string]*SwarmWorker
	RequestTorrentsFromWorker(string) ([]*deluge.Torrent, error)
	RequestFreeSpaceFromWorker(string) (uint64, error)
	SaveTorrentFile(string, string, *[]byte) (int64, error)
	RemoveTorrent(string, string, string, ...bool) (uint64, uint64, error)
}
