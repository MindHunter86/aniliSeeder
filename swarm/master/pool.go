package master

import (
	"sync"

	"github.com/hashicorp/yamux"
)

type workerPool struct {
	pool sync.Pool

	sync.RWMutex
	workers map[string]*worker
}

func newWorkerPool() *workerPool {
	return &workerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &worker{}
			},
		},
		workers: make(map[string]*worker),
	}
}

func (m *workerPool) newWorker(msess *yamux.Session) (wrk *worker, e error) {
	wrk = newWorker(msess)

	if e = wrk.connect(); e != nil {
		return
	}

	if m.isWorkerExists(wrk.getId()) {
		gLog.Warn().Str("worker_id", wrk.getId()).Msg("the worker is already exists in pool; seems there were connection errors")
		gLog.Debug().Str("worker_id", wrk.getId()).Msg("rewriting an exist record with the new worker")
	}

	m.Lock()
	m.workers[wrk.getId()] = wrk
	m.Unlock()

	return
}

func (m *workerPool) isWorkerExists(wid string) bool {
	m.RLock()
	m.RUnlock()

	return m.workers[wid] != nil
}

func (m *workerPool) dropWorker(wid string) (e error) {
	m.RLock()
	w := m.workers[wid]
	m.RUnlock()

	w.disconnect()

	m.Lock()
	m.workers[wid] = nil
	m.Unlock()

	return
}
