package app

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
)

var (
	errInsufficientSpace = errors.New("could not continue the deploy process because of insufficient space for some torrents")
	errFailedDeletions   = errors.New("could not continue the deploy process because of unsuccessful deletions")
	errFailedWorker      = errors.New("could not continue the delpoy process because one of workers errors")
	errNoFailures        = errors.New("there is nothing to redeploy; all torrents with OK announces")
	errNoWorkers         = errors.New("there is nothing to redeploy; all workers are unavailable")
	errNoTitles          = errors.New("there are no titles in the api response")

	errEmptyWorkerQueue = errors.New("no workers in balancing queue; seems that all workers haven't free space")

	errNothingDeploy   = errors.New("there is nothing to deploy")
	errNothingAssigned = errors.New("found some updates but there is now assigned titles")

	errSearchFailed = errors.New("there is abnormal name for searching title")
	errTitlesCount  = errors.New("abnormal titles length in the api response")
)

type (
	deploymentObject struct {
		workerId string

		aniTorrent *anilibria.TitleTorrent
		oldTorrent *deluge.Torrent

		sizeChanges            int64
		noDeploy, isDuplicated bool
	}
	workerTorrents struct {
		wid      string
		torrents []*deluge.Torrent
	}
	deploy struct{}

	deployType uint8
)

const (
	deployTypeAniUpdates deployType = iota
	deployTypeAniChanges
	deployTypeFailedAnnounces

	// TODO
	// deployTypeAniWatchers
	// deployTypeScheduler
)

func newWorkerTorrents(wid string) *workerTorrents {
	return &workerTorrents{
		wid: wid,
	}
}

func newDeploy() *deploy {
	return &deploy{}
}

// TODO
// func (*deploy) getAnilibriaScheduleTorrents() (e error) {
// 	return
// }

// TODO
// func (*deploy) getAnilibriaWatchingNowTorrents() (e error) {
// 	return
// }

func (*deploy) getAnilibriaTorrents(dtype deployType, payload ...interface{}) (trrs []*anilibria.TitleTorrent, e error) {
	var titles []*anilibria.Title

	switch dtype {
	case deployTypeAniUpdates:
		if titles, e = gAniApi.GetLastUpdates(); e != nil {
			return
		}
	case deployTypeAniChanges:
		if titles, e = gAniApi.GetLastChanges(); e != nil {
			return
		}
	case deployTypeFailedAnnounces:
		// panic avoid
		payload = append(payload, nil)

		if payload[0] == nil {
			return nil, errSearchFailed
		}

		if titles, e = gAniApi.SearchTitlesByName(payload[0].(*deluge.Torrent).GetName()); e != nil {
			return
		}

		if ltitles := len(titles); ltitles != 1 {
			gLog.Warn().Str("title_name", payload[0].(*deluge.Torrent).GetName()).Int("titles_count", ltitles).
				Msg("got a problem in searching failed titles; there are none, two or more titles in the result; manual search required")
			return nil, errTitlesCount
		}
	}

	if len(titles) == 0 {
		return nil, errNoTitles
	}

	for _, title := range titles {
		for _, trr := range title.Torrents.List {
			if tsize := utils.GetMBytesFromBytes(trr.TotalSize); tsize > int64(gCli.Uint64("anilibria-max-torrent-size")) {
				gLog.Info().Str("title_name", title.Names.En).Str("torrent_hash", trr.GetShortHash()).Int64("torrent_size_mb", tsize).
					Int("title_id", title.Id).Uint64("download_limit", gCli.Uint64("anilibria-max-torrent-size")).
					Msg("skipping a torrent because the torrent is larger than the download limit...")
				continue
			}

			trrs = append(trrs, trr)
		}
	}
	return
}

func (*deploy) getWorkersTorrents() (wts []*workerTorrents, _ bool) {
	ok, err := true, error(nil)

	for id := range gSwarm.GetConnectedWorkers() {
		wt := newWorkerTorrents(id)

		if wt.torrents, err = gSwarm.RequestTorrentsFromWorker(id); err != nil {
			ok = false

			gLog.Error().Str("worker_id", id).Err(err).Msg("could not get torrents from the given worker id; skipping...")
			continue
		}

		wts = append(wts, wt)
	}

	return wts, ok
}

// sometimes anilibria api forget to set the correct quality; on HEVC torrents it respond with [WEBRip 1080p]
func (*deploy) fixTorrentFileName(fname, quality, series string) (_ string, e error) {
	tname, _, ok := strings.Cut(fname, "AniLibria.TV")
	if !ok {
		return "", errors.New("there are troubles with fixing torrent name")
	}

	return tname + "AniLibria.TV" + " [" + quality + "][" + series + "]" + ".torrent", nil
}

func (*deploy) sortTorrentsByLeechers(aobjects []*deploymentObject) {
	sort.Slice(aobjects, func(i, j int) bool {
		return aobjects[i].aniTorrent.Leechers > aobjects[j].aniTorrent.Leechers
	})

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		for _, obj := range aobjects {
			gLog.Debug().Str("torrent_hash", obj.aniTorrent.GetShortHash()).
				Int64("torrent_size_mb", utils.GetMBytesFromBytes(obj.aniTorrent.TotalSize)).
				Int("torrent_leechers", obj.aniTorrent.Leechers).Msg("sorted slice debug")
		}
	}
}

func (*deploy) assignTorrentsToWorkers(aobjects []*deploymentObject) (ok bool, e error) {
	// add fake workers
	// for i := 0; i < 3; i++ {
	// 	id, fspace := uuid.NewV4().String(), 21474836480 //20gb

	// 	wrks[id] = &swarm.SwarmWorker{
	// 		Id:             id,
	// 		FreeSpace:      uint64(fspace),
	// 		ActiveTorrents: []*deluge.Torrent{},
	// 	}
	// 	fspaces[id] = uint64(fspace)

	// 	gLog.Debug().Str("worker_id", id).Msg("added fake worker")
	// }

	workers := gSwarm.GetConnectedWorkers()
	// no len(workers) check because it was previously

	queue := make(chan string, len(workers))

	// check workers for free space
	for id, worker := range workers {
		if _, e = gSwarm.RequestFreeSpaceFromWorker(id); e != nil {
			gLog.Error().Err(e).Str("worker_id", id).
				Msg("got an error in get free space from worker request; the worker will be removed from further deployment...")
			continue
		}

		// if worker has space leak or worker has not free space - remove from further balancing
		if worker.GetFreeSpace() == 0 {
			gLog.Warn().Str("worker_id", id).Uint64("free_space_kb", worker.GetFreeSpace()/1024).
				Msg("there is no free space on the worker; the worker will be removed from further deployment...")
			continue
		}

		// put worker in "balancer" (queue channel)
		gLog.Debug().Str("worker_id", id).Msg("collecting the worker for further balancing")
		queue <- id
	}

	// shit happens...
	if len(queue) == 0 {
		return false, errEmptyWorkerQueue
	}

	// balancing titles
	var trr *anilibria.TitleTorrent
	for _, aobject := range aobjects {
		// if no workers in the queue...
		if len(queue) == 0 {
			gLog.Debug().Msg("no free workers for further balancing; cancellation of the assignment process...")
		}

		trr = aobject.aniTorrent

		for wcount := len(queue); wcount != 0; wcount-- {
			// get worker id from balancer (queue)
			wid, ok := <-queue
			if !ok {
				gLog.Debug().Int("wcount", wcount).Msg("could not get worker id from the queue; ok is false; cancellation of the assignment...")
				break
			}

			// check worker free space for title deployment
			if !workers[wid].HasEnoughSpace(uint64(trr.TotalSize)) {
				gLog.Debug().Str("worker_id", wid).Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
					Int64("worker_size_mb", utils.GetMBytesFromBytes(int64(workers[wid].GetFreeSpace()))).
					Int64("torrent_size_mb", utils.GetMBytesFromBytes(trr.TotalSize)).
					Msg("the worker hasn't free space for the title's torrent; skipping the worker...")

				// put the worker at the end of the queue and try to assign to the next worker
				queue <- wid
				continue
			}

			// decrease the worker's free space for next torrents
			if !workers[wid].DecreaseFreeSpace(uint64(trr.TotalSize)) {
				gLog.Error().Str("worker_id", wid).Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
					Msg("got non ok in descrease free space for worker request; the worker will be removed from further deployment...")
				continue
			}

			// assign a deployment object (the torrent) to a worker
			aobject.workerId = wid
			gLog.Debug().Str("worker_id", wid).Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
				Msg("the torrent has been assigned to the worker")

			// the ok flag helps us detect pending deployments
			ok = true

			// put the worker at the end of the queue and try to assign to the next worker
			queue <- wid
			break
		}

		// set noDeploy flag for further sendDeployCommand func()
		if aobject.workerId == "" {
			gLog.Info().Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
				Msg("there is no worker assignment found; the worker will be removed from further deployment...")
			aobject.noDeploy = true
		}
	}

	return
}

func (m *deploy) deployAssignedTorrents(aobjects []*deploymentObject) {
	for _, aobject := range aobjects {
		trr := aobject.aniTorrent

		if aobject.noDeploy || aobject.isDuplicated {
			gLog.Debug().Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
				Msg("skipping torrent because of noDeploy or isDuplicated flag detection")
		}

		gLog.Debug().Str("worker_id", aobject.workerId).Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
			Msg("starting deploy process for the torrent...")

		name, fbytes, err := gAniApi.GetTitleTorrentFile(strconv.Itoa(trr.TorrentId))
		if err != nil {
			gLog.Error().Err(err).Msg("got an error in receiving the deploy file form the anilibria")
			break
		}

		// fix quality in torrentfile name (see func comments)
		gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Str("old_torrent_name", name).Msg("fixing torrent name...")
		if name, err = m.fixTorrentFileName(name, trr.Quality.String, trr.Series.String); err != nil {
			gLog.Error().Err(err).Msg("got an error in fixing torrent name")
			break
		}

		gLog.Debug().Str("torrent_name", name).Str("torrent_hash", trr.GetShortHash()).
			Msg("sending deploy request to the worker...")

		var wbytes int64
		if wbytes, err = gSwarm.SaveTorrentFile(aobject.workerId, name, fbytes); err != nil {
			gLog.Error().Err(err).Msg("got an error while processing the deploy request")
			continue
		}

		gLog.Info().Str("worker_id", aobject.workerId).Str("torrent_name", trr.GetName()).Str("torrent_hash", trr.GetShortHash()).
			Int64("written_bytes", wbytes).Msg("the torrent file has been sended to the worker")
	}
}
