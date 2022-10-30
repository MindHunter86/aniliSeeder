package app

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/rs/zerolog"
)

var (
	errInsufficientSpace = errors.New("could not continue the deploy process because of insufficient space for some torrents")
	errFailedDeletions   = errors.New("could not continue the deploy process because of unsuccessful deletions")
	errFailedWorker      = errors.New("could not continue the delpoy process because one of workers errors")
	errNoFailures        = errors.New("there is nothing to redeploy; all torrents with OK announces")
	errNoWorkers         = errors.New("there is nothing to redeploy; all workers are unavailable")

	errNothingDeploy   = errors.New("there is nothing to deploy")
	errNothingAssigned = errors.New("found some updates but there is now assigned titles")
)

type deploy struct{}

// type deployType uint8

// const (
// 	dplAnilibriaUpdates deployType = iota
// 	// dplAnilibriaSchedule
// 	// dplAnilibriaWatchingNow
// )

func newDeploy() *deploy {
	return &deploy{}
}

func (*deploy) getWorkersTorrents() (trrs []*deluge.Torrent, e error) {
	for id := range gSwarm.GetConnectedWorkers() {
		var wtrrs []*deluge.Torrent
		if wtrrs, e = gSwarm.RequestTorrentsFromWorker(id); e != nil {
			gLog.Error().Str("worker_id", id).Err(e).Msg("could not get torrents from the given worker id; skipping...")
			continue
		}

		trrs = append(trrs, wtrrs...)
	}

	return
}

func (*deploy) fixTorrentFileName(fname, quality, series string) (_ string, e error) {
	tname, _, ok := strings.Cut(fname, "AniLibria.TV")
	if !ok {
		return "", errors.New("there are troubles with fixing torrent name")
	}

	return tname + "AniLibria.TV" + " [" + quality + "][" + series + "]" + ".torrent", nil
}

func (*deploy) sortTorrentListByLeechers(trrs []*anilibria.TitleTorrent) (_ []*anilibria.TitleTorrent) {
	sort.Slice(trrs, func(i, j int) bool {
		return trrs[i].Leechers > trrs[j].Leechers
	})

	// debug
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		for _, trr := range trrs {
			gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Int64("torrent_size_mb", trr.TotalSize/1024/1024).
				Int("torrent_leechers", trr.Leechers).Msg("sorted slice debug")
		}
	}

	return trrs
}

func (*deploy) balanceForWorkers(trrs []*anilibria.TitleTorrent) (_ map[string][]*anilibria.TitleTorrent, e error) {
	wrks := gSwarm.GetConnectedWorkers()
	var fspaces = make(map[string]uint64)

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

	if len(wrks) == 0 {
		return nil, errors.New("there is no avaliable workers for the balancing process")
	}
	blncr := make(chan string, len(wrks))

	for id := range wrks {
		if fspaces[id] == 0 {
			if fspaces[id], e = gSwarm.RequestFreeSpaceFromWorker(id); e != nil {
				gLog.Error().Err(e).Str("worker_id", id).Msg("got an error in free space request to the worker")
				continue
			}
		}

		gLog.Debug().Str("worker_id", id).Msg("collecting the worker for further balancing")
		blncr <- id
	}

	if len(blncr) == 0 {
		return nil, errors.New("there is no workers with free space for the balancing process")
	}

	var wtitles = make(map[string][]*anilibria.TitleTorrent)

loop:
	for {

		w := <-blncr
		var assigned bool

		for id, trr := range trrs {
			if trr == nil {
				continue
			}

			if uint64(trr.TotalSize) > fspaces[w] {
				gLog.Info().Str("worker_id", w).Str("torrent_hash", trr.GetShortHash()).Int64("fspace", int64(fspaces[w])).Int64("tspace", trr.TotalSize).
					Msg("skipping torrents because of insufficient disk space on the worker")
				continue
			}

			// decrease the workers free space
			gLog.Debug().Uint64("old_fspace", fspaces[w]).Int64("torrent_size", trr.TotalSize).
				Uint64("new_fspace", fspaces[w]-uint64(trr.TotalSize)).Msg("decreasing the worker's free space")
			fspaces[w] = fspaces[w] - uint64(trr.TotalSize)

			// assigning the torrent to the worker
			var atrr = new(anilibria.TitleTorrent)
			*atrr = *trr
			wtitles[w] = append(wtitles[w], atrr)

			// remove the title from a slice
			trrs[id] = nil

			assigned = true
			gLog.Debug().Str("worker_id", w).Str("torrent_hash", trr.GetShortHash()).Msg("the torrent has been assigned")

			break
		}

		if assigned {
			gLog.Debug().Str("worker_id", w).Msg("put the worker into balancer chan")
			blncr <- w
		}

		if len(blncr) != 0 {
			gLog.Debug().Int("balance_queue", len(blncr)).Msg("found workers in the balancing chan")
			continue
		}

		gLog.Debug().Msg("there is no avaliable workers for the balancing process")
		break loop
	}

	return wtitles, e
}

func (m *deploy) sendDeployCommand(deployTasks map[string][]*anilibria.TitleTorrent) {
	var e error

	for wid, trrs := range deployTasks {
		gLog.Debug().Str("worker_id", wid).Msg("starting deploy process for the worker...")

		for _, trr := range trrs {
			name, fbytes, err := gAniApi.GetTitleTorrentFile(strconv.Itoa(trr.TorrentId))
			if err != nil {
				gLog.Error().Err(e).Msg("got an error in receiving the deploy file form the anilibria")
				break
			}

			gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Str("old_torrent_name", name).Msg("fixing torrent name...")
			if name, e = m.fixTorrentFileName(name, trr.Quality.String, trr.Series.String); e != nil {
				gLog.Error().Err(e).Msg("got an error in fixing torrent name")
				break
			}

			gLog.Debug().Str("torrent_name", name).Str("torrent_hash", trr.GetShortHash()).
				Msg("sendind deploy request to the worker...")

			var wbytes int64
			if wbytes, e = gSwarm.SaveTorrentFile(wid, name, fbytes); e != nil {
				gLog.Error().Err(e).Msg("got an error while processing the deploy request")
				continue
			}

			gLog.Info().Str("worker_id", wid).Int64("written_bytes", wbytes).
				Msg("the torrent file has been sended to the worker")
		}

		gLog.Debug().Str("worker_id", wid).Msg("deploy process for the worker has been finished")
	}
}
