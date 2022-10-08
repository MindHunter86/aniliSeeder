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

	m.Lock()
	m.workers[wrk.getId()] = wrk
	m.Unlock()

	return
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
