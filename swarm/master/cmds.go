package master

import (
	"errors"

	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/swarm"
)

var (
	errWorkerNotFound = errors.New("there is not worker with such id")
)

func (*Master) IsMaster() bool {
	return true
}

func (m *Master) GetConnectedWorkers() (_ map[string]*swarm.SwarmWorker) {
	var wrks = make(map[string]*swarm.SwarmWorker)

	for _, id := range m.workerPool.getWorkerIds() {
		wrk := m.workerPool.getWorker(id)

		wrks[id] = &swarm.SwarmWorker{
			Id:             wrk.id,
			Version:        wrk.version,
			FreeSpace:      wrk.wdFreeSpace,
			ActiveTorrents: wrk.trrs,
		}
	}

	return wrks
}

func (m *Master) RequestTorrentsFromWorker(wid string) ([]*deluge.Torrent, error) {
	if !m.workerPool.isWorkerExists(wid) {
		return nil, errWorkerNotFound
	}

	wrk := m.workerPool.getWorker(wid)
	return wrk.getTorrents()
}

func (m *Master) RequestFreeSpaceFromWorker(wid string) (uint64, error) {
	if !m.workerPool.isWorkerExists(wid) {
		return 0, errWorkerNotFound
	}

	wrk := m.workerPool.getWorker(wid)
	return wrk.getFreeSpace()
}

func (m *Master) SaveTorrentFile(wid string, fname string, fbytes *[]byte) (int64, error) {
	if !m.workerPool.isWorkerExists(wid) {
		return 0, errWorkerNotFound
	}

	wrk := m.workerPool.getWorker(wid)
	return wrk.saveTorrentFile(fname, fbytes)
}

func (m *Master) RemoveTorrent(wid string, name string, hash string, wdata ...bool) (uint64, uint64, error) {
	if !m.workerPool.isWorkerExists(wid) {
		return 0, 0, errWorkerNotFound
	}

	// panic avoid
	wdata = append(wdata, false)
	return m.workerPool.getWorker(wid).deleteTorrent(name, hash, wdata[0])
}

func (m *Master) ForceReannounce(wid string) error {
	if !m.workerPool.isWorkerExists(wid) {
		return errWorkerNotFound
	}

	return m.workerPool.getWorker(wid).forceReannounce()
}
