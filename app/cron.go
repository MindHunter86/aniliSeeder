package app

import "time"

type cron struct{}

func newCron() *cron {
	return &cron{}
}

func (*cron) run() {
	//

	ticker := time.NewTimer(time.Second)

loop:
	for {
		select {
		case <-gCli.Done():
			gLog.Warn().Msg("context done() has been caught; stopping cron subservice...")
			break loop
		case <-ticker.C:
			// tm := time.Now()
		}
	}
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
