package app

import (
	"sort"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/utils"
	"github.com/rs/zerolog"
)

// type failedTitle struct {
// 	workerId   string
// 	oldTorrent *deluge.Torrent
// 	aniTorrent *anilibria.TitleTorrent

// 	sizeChanges int64

// 	noDeploy     bool
// 	isDuplicated bool
// }

// func (m *deploy) dryDeployFailedAnnounces() ([]*failedTitle, error) {
// 	return m.deployFailedAnnounces(true)
// }

func (m *deploy) deployFailedAnnounces(dryrun ...bool) (ftitles []*deploymentObject, e error) {
	wtorrents, ok := m.getWorkersTorrentsV2()
	if !ok {
		return nil, errFailedWorker
	}

	if len(wtorrents) == 0 {
		return nil, errNoWorkers
	}

	if ftitles, e = m.searchFailedTitles(wtorrents); e != nil {
		return
	}

	if len(ftitles) == 0 {
		gLog.Debug().Msg("there is nothing for redeployment")
		return nil, errNoFailures
	}

	m.sortTitlesByLeechers(ftitles)
	m.searchForDuplicates(ftitles, wtorrents)

	ok = m.isSpaceEnoughForUpdate(ftitles)
	if !ok && !gCli.Bool("deploy-ignore-errors") {
		return nil, errInsufficientSpace
	}

	// redeploy ...

	// panic avoid
	dryrun = append(dryrun, true)
	if !dryrun[0] {
		ok = m.dropFailedTorrent(ftitles)
		if !ok && !gCli.Bool("deploy-ignore-errors") {
			return nil, errFailedDeletions
		}

		m.deployAssignedTorrents(ftitles)
	}

	// TODO
	// return failed titles with redeploy status (OK\nonOK)
	return ftitles, e
}

func (*deploy) getWorkersTorrentsV2() (_ []*workerTorrents, ok bool) {
	var wts []*workerTorrents
	var err error

	ok = true
	for id := range gSwarm.GetConnectedWorkers() {
		wt := &workerTorrents{
			wid: id,
		}

		if wt.torrents, err = gSwarm.RequestTorrentsFromWorker(id); err != nil {
			gLog.Error().Str("worker_id", id).Err(err).Msg("could not get torrents from the given worker id; skipping...")
			ok = false
			continue
		}

		wts = append(wts, wt)
	}

	return wts, ok
}

func (*deploy) searchFailedTitles(wtorrents []*workerTorrents) (_ []*deploymentObject, e error) {
	var ftitles []*deploymentObject

	for _, worker := range wtorrents {
		for _, trr := range worker.torrents {

			// skip torrents with OK and WARN status in announce
			switch trr.GetTrackerStatus() {
			case deluge.TrackerStatusNotRegistered:
				gLog.Trace().Str("torrent_hash", trr.GetShortHash()).Str("torrent_status", trr.GetTrackerRawError()).
					Msg("the torrent has NotRegistered status; processing it...")
			case deluge.TrackerStatusOK:
				gLog.Trace().Str("torrent_hash", trr.GetShortHash()).Str("torrent_status", trr.GetTrackerRawError()).
					Msg("the torrent has OK status; skipping...")
				continue
			default:
				gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Str("torrent_status", trr.GetTrackerRawError()).
					Msg("there is torrent with warning announce status; skipping...")
				continue
			}

			// get titltes from anilibria API for failed torrents
			gLog.Debug().Str("torrents_hash", trr.GetShortHash()).Msg("torrent marked as failed")

			var titles []*anilibria.Title
			if titles, e = gAniApi.SearchTitlesByName(trr.GetName()); e != nil {
				return
			}

			if ltitles := len(titles); ltitles != 1 {
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.GetName()).Int("titles_count", ltitles).
					Msg("got a problem in searching failed titles; there are none, two or more titles in the result; manual search required")

				// TODO
				// ?? Telegram Alert... github.com/MindHunter86/aniliSeeder/issues/55
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.GetName()).
					Msg("failed torrent will be skipped and not deleted")
				continue
			}

			// trying to find failed torrent by quality
			var found bool
			for _, anitrr := range titles[0].Torrents.List {
				if trr.GetQuality() != anitrr.Quality.String {
					gLog.Debug().Str("title_name", trr.GetName()).Str("torrent_found_quality", anitrr.Quality.String).
						Msg("anilibria quality found but skipped")
					continue
				}

				gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.GetName()).
					Str("new_torrent_hash", anitrr.GetShortHash()).Msg("the title's torrent replacement has been found")

				ftitles = append(ftitles, &deploymentObject{
					workerId:   worker.wid,
					oldTorrent: trr,
					aniTorrent: anitrr,
				})
				found = true
			}

			if !found {
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.GetName()).Str("torrent_quality", trr.GetQuality()).
					Msg("there is a problem in searching title's torrent by quality string; manual search required")
			}
		}
	}

	return ftitles, e
}

func (*deploy) sortTitlesByLeechers(ftitles []*deploymentObject) {
	sort.Slice(ftitles, func(i, j int) bool {
		return ftitles[i].aniTorrent.Leechers > ftitles[j].aniTorrent.Leechers
	})

	// debug
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		for _, ftitle := range ftitles {
			gLog.Debug().Str("torrent_hash", ftitle.oldTorrent.GetShortHash()).Str("torrent_hash_new", ftitle.aniTorrent.GetShortHash()).
				Int64("torrent_size_mb", utils.GetMBytesFromBytes(ftitle.aniTorrent.TotalSize)).
				Int("torrent_leechers", ftitle.aniTorrent.Leechers).Msg("sorted slice debug")
		}
	}
}

func (*deploy) isSpaceEnoughForUpdate(ftitles []*deploymentObject) (ok bool) {
	var fspaces = make(map[string]uint64)

	ok = true
	for _, ftitle := range ftitles {
		// cache worker's free space for another titles
		if fspaces[ftitle.workerId] == 0 {
			var e error
			if fspaces[ftitle.workerId], e = gSwarm.RequestFreeSpaceFromWorker(ftitle.workerId); e != nil {
				gLog.Warn().Err(e).Msg("got an error in get free space for worker operation in deploy process")
				ok = false
			}

			gLog.Debug().Str("worker_id", ftitle.workerId).Uint64("worker_free_space", fspaces[ftitle.workerId]).
				Msg("title's worker free space debug")
		}

		// title's deploy status definition
		ftitle.sizeChanges = ftitle.aniTorrent.TotalSize - ftitle.oldTorrent.TotalSize
		if fspaces[ftitle.workerId]-uint64(ftitle.sizeChanges) <= 0 {
			gLog.Warn().Str("torrent_hash", ftitle.oldTorrent.GetShortHash()).Str("torrent_name", ftitle.oldTorrent.GetName()).
				Int64("torrent_size_changes", ftitle.sizeChanges).Msg("could not deploy the torrents because of insufficient space")
			ftitle.noDeploy = true
		}

		gLog.Debug().Str("torrent_name", ftitle.oldTorrent.GetName()).Int64("torrent_size_changes", ftitle.sizeChanges).
			Msg("torrent ready for deletion")
	}

	return
}

func (*deploy) dropFailedTorrent(ftitles []*deploymentObject) (ok bool) {
	ok = true

	for _, ftitle := range ftitles {
		if ftitle.noDeploy {
			gLog.Debug().Str("torrent_name", ftitle.oldTorrent.GetName()).
				Msg("the torrent marked as noDeploy, so delition process will be skipped too")
			continue
		}

		dbytes, tbytes, err := gSwarm.RemoveTorrent(ftitle.workerId, ftitle.oldTorrent.Hash, ftitle.oldTorrent.GetName())
		if err != nil {
			gLog.Error().Err(err).Str("torrent_hash", ftitle.aniTorrent.GetShortHash()).Str("torrent_name", ftitle.oldTorrent.GetName()).
				Msg("got an error in torrent removing process; skipping the torrent...")
			ok = false
			continue
		}

		if dbytes != 0 {
			gLog.Warn().Str("torrent_hash", ftitle.aniTorrent.GetShortHash()).Str("torrent_name", ftitle.oldTorrent.GetName()).
				Msg("an internal error has occurred, operator intervention is required")
			ok = false
			continue
		}

		gLog.Info().Str("torrent_hash", ftitle.aniTorrent.GetShortHash()).Str("torrent_name", ftitle.oldTorrent.GetName()).
			Msg("the torrent has been removed; it is now ready to be redeployed")

		gLog.Debug().Str("worker_id", ftitle.workerId).Uint64("worker_fspace", tbytes).
			Msg("worker free space debug after torrent deleting")
	}

	return ok
}

func (*deploy) searchForDuplicates(ftitles []*deploymentObject, wtorrents []*workerTorrents) {
	for _, ftitle := range ftitles {
		if ftitle.noDeploy {
			continue
		}

		for _, worker := range wtorrents {
			for _, trr := range worker.torrents {
				if trr.Hash == ftitle.aniTorrent.Hash {
					ftitle.isDuplicated = true
					gLog.Debug().Str("torrent_name", ftitle.oldTorrent.GetName()).
						Msg("found duplication in searched titles; the torrent marked as 'duplicated' and will be skipped")
				}
			}
		}
	}
}
