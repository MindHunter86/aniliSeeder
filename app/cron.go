package app

import (
	"sync"
	"time"
)

type cron struct {
	ticker *time.Ticker

	wg sync.WaitGroup

	mu    sync.RWMutex
	tasks cronTask
}
type cronTask uint8

const (
	cronTaskDeployUpdates cronTask = 1 << iota
	cronTaskReannounce
)

func (m *cronTask) toggle(lock sync.Locker, task cronTask) {
	lock.Lock()
	defer lock.Unlock()

	*m = *m ^ task
}

func (m *cronTask) isEnabled(lock sync.Locker, task cronTask) (ct cronTask) {
	lock.Lock()
	defer lock.Unlock()

	ct = *m & task // copy value and unlock main cronTask
	return ct
}

func (m *cronTask) getTasks(lock sync.Locker) (tasks uint8) {
	lock.Lock()
	defer lock.Unlock()

	tasks = uint8(*m)
	return
}

func newCron() *cron {
	return &cron{}
}

func (m *cron) run() {
	if gCli.Bool("cron-disable") {
		gLog.Warn().Msg("cron-disable flag detected; cron services will be deleted from bootstrap process")
		return
	}

	gLog.Debug().Time("time_now", time.Now()).Msg("starting cron subservice...")
	m.ticker = time.NewTicker(time.Minute)

loop:
	for {
		select {
		case <-gCtx.Done():
			gLog.Warn().Msg("context done() has been caught; stopping cron subservice...")
			break loop
		case <-m.ticker.C:
			m.runCronTasks()
		}
	}

	m.ticker.Stop()
	gLog.Debug().Msg("waiting for goroutines with cron jobs...")
	m.wg.Wait()
}

//	triggers
//	1 min - collect stats and push to graphite (?)
//
//	5 mins - check torrents announces
//		if have failed announces - search titles and update torrents
//
//	10 mins - check drydeploy report
//		if report has changes - try to deploy new titles
//			- check free space, try to deploy without deletions
//			if have not free space
//				- collect all torrnets with their VKscore
//				- delete torrnets with small VKscore and get free space
//				- try to deploy again...
//
//	60 mins - check for updates
//		if updates found
//			- try to update
//			- push notification to telegram

func (m *cron) runCronTasks() {
	gLog.Debug().Uint8("cron_tasks", m.tasks.getTasks(m.mu.RLocker())).Msg("mask before switch")

	tm := time.Now()

	// < 1 min jobs here >
	gLog.Debug().Msg("running 1min jobs...")

	if tm.Minute()%5 == 0 {
		// < 5 min jobs here >
		gLog.Debug().Msg("running 5min jobs...")

		// reannounce
		m.reannounce()

		// !! temporary shit
		// https://github.com/MindHunter86/aniliSeeder/issues/73
		time.Sleep(30 * time.Second)

		// redeploy
		m.redeploy()
	}

	if tm.Minute()%60 == 0 {
		// < 1 hour jobs here >
		gLog.Debug().Msg("running 60min jobs...")
	}

	gLog.Debug().Uint8("cron_tasks", m.tasks.getTasks(m.mu.RLocker())).Msg("mask after switch")
}

func (m *cron) redeploy() {
	if m.tasks.isEnabled(m.mu.RLocker(), cronTaskDeployUpdates) != 0 {
		gLog.Debug().Msg("deploy updates is locked now; skipping job...")
		return
	}

	m.tasks.toggle(m.mu.RLocker(), cronTaskDeployUpdates)
	gLog.Debug().Msg("deploy updates is not locked; running job...")

	m.wg.Add(1)
	go func(done func()) {
		if _, e := newDeploy().run(); e != nil && e != errNothingDeploy && e != errNothingAssigned {
			gLog.Error().Err(e).Msg("got an error in cron deployUpdates job")
		}

		//! TODO
		//! if errNothingAssigned
		//!   try to clean weak VKscore torrents

		m.tasks.toggle(m.mu.RLocker(), cronTaskDeployUpdates)
		done()
	}(m.wg.Done)
}

func (m *cron) reannounce() {
	if m.tasks.isEnabled(m.mu.RLocker(), cronTaskReannounce) != 0 {
		gLog.Debug().Msg("reannounces is locked now; skipping job...")
		return
	}

	m.tasks.toggle(m.mu.RLocker(), cronTaskReannounce)
	gLog.Debug().Msg("reannounces is not locked; running job...")

	m.wg.Add(1)
	go func(done func()) {
		var e error
		wrks := gSwarm.GetConnectedWorkers()

		for wid := range wrks {
			if e = gSwarm.ForceReannounce(wid); e != nil {
				gLog.Error().Err(e).Msg("got an error in force reannounce request")
			}
		}

		// wait for reannounces && force redeploy
		if _, e = newDeploy().deployFailedAnnounces(); e != nil && e != errNoFailures {
			gLog.Error().Err(e).Msg("got an error in deploy failed_announces request")
		}

		m.tasks.toggle(m.mu.RLocker(), cronTaskReannounce)
		done()
	}(m.wg.Done)
}
