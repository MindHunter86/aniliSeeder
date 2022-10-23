package app

import (
	"errors"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
)

var (
	errDplWorkerUnavailable = errors.New("")
	errDplFailedAnnounces   = errors.New("")
)

type failedTitle struct {
	workerId   string
	oldTorrent *deluge.Torrent
	aniTorrent *anilibria.TitleTorrent

	sizeChanges int64

	noDeploy bool
}

type workerTorrents struct {
	wid      string
	torrents []*deluge.Torrent
}

func (m *deploy) dryDeployFailedAnnounces() error {
	return m.deployFailedAnnounces(true)
}

func (m *deploy) deployFailedAnnounces(dryrun ...bool) (e error) {
	// wtorrents, ok := m.getWorkersTorrentsV2()
	// if !ok {
	// 	return errors.New("could not continue the delpoy process because one of workers errors")
	// }

	//

	//
	return
}

func (*deploy) getWorkersTorrentsV2() (_ []*workerTorrents, ok bool) {
	var wts []*workerTorrents
	var err error

	ok = true
	for id := range gSwarm.GetConnectedWorkers() {
		var wt *workerTorrents
		if wt.torrents, err = gSwarm.RequestTorrentsFromWorker(id); err != nil {
			gLog.Error().Str("worker_id", id).Err(err).Msg("could not get torrents from the given worker id; skipping...")
			ok = false
			continue
		}

		wts = append(wts, wt)
	}

	return wts, ok
}

func (*deploy) searchFailedTitles(wtorrents []*workerTorrents) (_ []*failedTitle, e error) {
	var ftitles []*failedTitle

	for _, worker := range wtorrents {
		for _, trr := range worker.torrents {

			// skip the torrents with OK announces
			if trr.IsTrackerOk() {
				continue
			}

			// get titltes from anilibria API for failed torrents
			gLog.Debug().Str("torrents_hash", trr.GetShortHash()).Msg("torrent marked as failed")

			var titles []*anilibria.Title
			if titles, e = gAniApi.SearchTitlesByName(trr.Name); e != nil {
				return
			}

			if ltitles := len(titles); ltitles != 1 {
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.Name).Int("titles_count", ltitles).
					Msg("got a problem in searching failed titles; there are none, two or more titles in the result; manual search required")

				// TODO
				// ?? Telegram Alert... github.com/MindHunter86/aniliSeeder/issues/55
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.Name).
					Msg("failed torrent will be skipped and not deleted")
				continue
			}

			// trying to find failed torrent by quality
			var found bool
			for _, anitrr := range titles[1].Torrents.List {
				if trr.GetTorrentQuality() != anitrr.Quality.String {
					continue
				}

				gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.Name).
					Str("new_torrent_hash", anitrr.GetShortHash()).Msg("the title's torrent replacement has been found")

				ftitles = append(ftitles, &failedTitle{
					workerId:   worker.wid,
					oldTorrent: trr,
					aniTorrent: anitrr,
				})
				found = true
			}

			if !found {
				gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.Name).
					Msg("there is a proble in searching title's torrent by quality string; manual search required")
			}
		}
	}

	return ftitles, e
}

func (*deploy) sortTitlesByLeechers(ftitles []*failedTitle) {

}

func (*deploy) isSpaceEnoughForUpdate(ftitles []*failedTitle) (ok bool) {
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
			gLog.Warn().Str("torrent_hash", ftitle.oldTorrent.GetShortHash()).Str("torrnet_name", ftitle.oldTorrent.GetName()).
				Int64("torrent_size_changes", ftitle.sizeChanges).Msg("could not deploy the torrents because of insufficient space")
			ftitle.noDeploy = true
		}
	}

	return
}

func (*deploy) dropFailedTorrent(ftitles []*failedTitle) error {
	// for _, ftitle := range ftitles {

	// }

	return nil
}
