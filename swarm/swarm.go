package swarm

import (
	"sync"

	"github.com/MindHunter86/aniliSeeder/deluge"
)

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
	ForceReannounce(string) error
}

type SwarmWorker struct {
	Id      string
	Version string

	sync.RWMutex
	FreeSpace uint64

	ActiveTorrents []*deluge.Torrent
}

func (m *SwarmWorker) GetFreeSpace() (space uint64) {
	m.RLock()
	defer m.RUnlock()

	space = m.FreeSpace
	return
}

func (m *SwarmWorker) HasEnoughSpace(space uint64) bool {
	m.RLock()
	defer m.RUnlock()

	return m.FreeSpace > space // ? is it works ?
}

func (m *SwarmWorker) DecreaseFreeSpace(space uint64) bool {
	if !m.HasEnoughSpace(space) {
		return false
	}

	m.Lock()
	defer m.Unlock()

	m.FreeSpace = m.FreeSpace - space
	return true
}
