package app

import (
	"sync"
	"time"
)

type cron struct {
	ticker *time.Ticker

	wg sync.WaitGroup

	mu    sync.Mutex
	tasks cronTask
}
type cronTask uint8

const (
	cronTaskDeployUpdates cronTask = iota
)

func newCron() *cron {
	return &cron{}
}

func (m *cron) run() {
	gLog.Debug().Time("time_now", time.Now()).Msg("starting cron subservice...")

	m.wg.Add(1)
	defer m.wg.Done()

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

func (*cron) stop() {
	return
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

//

func (m *cron) runCronTasks() {
	gLog.Debug().Uint8("cron_tasks", uint8(m.tasks)).Msg("mask before switch")

	tm := time.Now()

	switch {
	// 1 min
	case tm.Minute()%1 == 0:
		gLog.Debug().Msg("running 1min jobs...")

	// 5 mins
	// !! BUG
	// !! Not worked
	case tm.Minute()%5 == 0:
		gLog.Debug().Msg("running 5min jobs...")

		m.mu.Lock()
		m.tasks = m.tasks ^ cronTaskDeployUpdates
		m.mu.Unlock()

		m.wg.Add(1)
		go func(done func()) {
			defer done()
			if err := m.deployUpdates(); err != nil {
				gLog.Error().Err(err).Msg("got an error in cron deployUpdates job")
			}
		}(m.wg.Done)

	// 60 mins
	case tm.Minute() == 0:
		gLog.Debug().Msg("running 60min jobs...")

	default:
		gLog.Debug().Msg("no jobs for running now")
		return
	}

	gLog.Debug().Uint8("cron_tasks", uint8(m.tasks)).Msg("mask after switch")
}

func (*cron) checkTorrentsAnnounces() (e error) {
	return
}

func (m *cron) deployUpdates() (e error) {
	gLog.Info().Msg("starting check deployment status cronjob...")

	if e = m.checkDryDeployReport(); e != nil && e != errNothingDeploy {
		return
	}

	if e == errNothingDeploy {
		gLog.Info().Msg("there is nothing to deploy; deploy stopped")
		return nil
	}

	// !!
	// TODO try to clean weak VKscore torrents
	// !!

	gLog.Info().Msg("dry deploy says that anilibria has some updates; deploying titles...")
	return m.sendDeployCommand()
}

func (*cron) checkDryDeployReport() (e error) {
	_, e = newDeploy().dryRun()
	return
}

func (*cron) sendDeployCommand() (e error) {
	_, e = newDeploy().run()
	return
}
