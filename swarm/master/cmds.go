package master

import (
	"github.com/MindHunter86/aniliSeeder/swarm"
)

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
