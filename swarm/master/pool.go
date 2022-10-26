package master

import (
	"sync"
	"time"

	"github.com/hashicorp/yamux"
)

type workerPool struct {
	sync.RWMutex
	workers map[string]*worker
}

func newWorkerPool() *workerPool {
	return &workerPool{
		workers: make(map[string]*worker),
	}
}

func (m *workerPool) newWorker(msess *yamux.Session, mid string) (wrk *worker, e error) {
	wrk = newWorker(msess, mid)

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
	defer m.RUnlock()

	return m.workers[wid] != nil
}

func (m *workerPool) getWorkerIds() []string {
	m.RLock()
	defer m.RUnlock()

	var ids []string
	for id := range m.workers {
		if m.workers[id] == nil {
			gLog.Warn().Str("worker_id", id).Msg("abnormal worker detected in pool")
			continue
		}

		ids = append(ids, id)
	}

	return ids
}

func (m *workerPool) getWorker(wid string) *worker {
	m.RLock()
	defer m.RUnlock()

	return m.workers[wid]
}

func (m *workerPool) findDeadWorkers() {
	var wrks = make(map[string]*worker)

	m.RLock()
	wrks = m.workers
	m.RUnlock()

	for wid, wrk := range wrks {
		if wrk == nil {
			continue
		}

		if e := wrk.isMuxSessionAlive(); e != nil {
			gLog.Error().Err(e).Msg("got an error in mux session validataion; removing worker from pool...")
			m.dropWorker(wid)
		}

		gLog.Trace().Str("worker_id", wid).Msg("worker is alive")
	}
}

func (m *worker) isMuxSessionAlive() (e error) {
	var du time.Duration
	if du, e = m.msess.Ping(); e != nil {
		return
	}

	gLog.Trace().Str("worker_id", m.id).Dur("worker_mux_ping", du).Msg("")
	return
}

func (m *workerPool) dropWorker(wid string) {
	m.Lock()
	w := m.workers[wid]
	delete(m.workers, wid)
	m.Unlock()

	w.gconn.Close()
	w.msess.Close()
}
