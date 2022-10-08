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
	IsMaster() bool
	Bootstrap() error
	GetConnectedWorkers() map[string]*SwarmWorker
}
