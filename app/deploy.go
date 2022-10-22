package app

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
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

func (m *deploy) run() (map[string][]anilibria.TitleTorrent, error) {
	return m.deploy(false)
}

func (m *deploy) dryRun() (map[string][]anilibria.TitleTorrent, error) {
	return m.deploy(true)
}

func (m *deploy) deploy(isDryRun bool) (_ map[string][]anilibria.TitleTorrent, e error) {
	var titles []*anilibria.TitleTorrent
	if titles, e = m.getAnilibriaUpdatesTorrents(); e != nil {
		return
	}

	var torrents []*deluge.Torrent
	if torrents, e = m.getWorkersTorrents(); e != nil {
		return
	}

	titleUpdates := m.compareUpdateListWithTorrents(titles, torrents)

	sortedUpdates := m.sortTorrentListByLeechers(titleUpdates)

	var assignedTitles = make(map[string][]anilibria.TitleTorrent)
	if assignedTitles, e = m.balanceForWorkers(sortedUpdates); e != nil {
		return
	}

	if len(assignedTitles) == 0 {
		return nil, errors.New("there is nothing to deploy")
	}

	if !isDryRun {
		m.sendDeployCommand(assignedTitles)
	}

	return assignedTitles, e
}

func (m *deploy) sendDeployCommand(deployTasks map[string][]anilibria.TitleTorrent) {
	var e error

	for wid, trrs := range deployTasks {
		gLog.Debug().Str("worker_id", wid).Msg("starting deploy process for the worker...")

		for _, trr := range trrs {
			name, fbytes, err := gAniApi.GetTitleTorrentFile(strconv.Itoa(trr.TorrentId))
			if err != nil {
				gLog.Error().Err(e).Msg("got an error in receiving the deploy file form the anilibria")
				break
			}

			gLog.Debug().Str("torrent_hash", trr.Hash[0:9]).Str("old_torrent_name", name).Msg("fixing torrent name...")
			if name, e = m.fixTorrentFileName(name, trr.Quality.String, trr.Series.String); e != nil {
				gLog.Error().Err(e).Msg("got an error in fixing torrent name")
				break
			}

			gLog.Debug().Str("torrent_name", name).Str("torrent_hash", trr.Hash[0:9]).
				Msg("sendind deploy request to the worker...")

			var wbytes int64
			if wbytes, e = gSwarm.SaveTorrentFile(wid, name, fbytes); e != nil {
				gLog.Error().Err(e).Msg("got an error while processing the deploy request")
				continue
			}

			gLog.Info().Str("worker_ud", wid).Int64("written_bytes", wbytes).
				Msg("the torrent file has been sended to the worker")
		}

		gLog.Debug().Str("worker_id", wid).Msg("deploy proccess for the worker has been finished")
	}
}

func (*deploy) fixTorrentFileName(fname, quality, series string) (_ string, e error) {
	tname, _, ok := strings.Cut(fname, "AniLibria.TV")
	if !ok {
		return "", errors.New("there are troubles with fixing torrent name")
	}

	return tname + "AniLibria.TV" + " [" + quality + "][" + series + "]" + ".torrent", nil
}

//

func (*deploy) getAnilibriaUpdatesTorrents() (trrs []*anilibria.TitleTorrent, e error) {
	var ttls []*anilibria.Title
	if ttls, e = gAniApi.GetLastUpdates(); e != nil {
		return
	}

	for _, ttl := range ttls {
		trrs = append(trrs, ttl.Torrents.List...)
	}

	return
}

// func (*deploy) getAnilibriaScheduleTorrents() (e error) {
// 	return
// }

// func (*deploy) getAnilibriaWatchingNowTorrents() (e error) {
// 	return
// }

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

// TODO optimize
// !! WARNING
// !! There is no comparing by torrent name!!!
// !! Some torrents may be removed from the anilibria announces
// !! and updated their hashes because of title update
// !! Must be implemented shortly
func (*deploy) compareUpdateListWithTorrents(atrrs []*anilibria.TitleTorrent, wtrrs []*deluge.Torrent) (mtrrs []*anilibria.TitleTorrent) {
	for _, atrr := range atrrs {
		found := false

		for _, wtrr := range wtrrs {
			if wtrr.Hash != atrr.Hash {
				continue
			}

			found = true
			break
		}

		if found != true {
			gLog.Debug().Str("hash", atrr.Hash).Msg("torrent compare process: missed hash found")
			mtrrs = append(mtrrs, atrr)
			continue
		}

		// !! Check by name and series ...
		// TODO
	}

	// debug
	for _, trr := range mtrrs {
		gLog.Debug().Str("torrent_hash", trr.Hash[0:9]).Int64("torrent_size", trr.TotalSize).Msg("there is a new item for the further deploy")
	}

	return
}

func (*deploy) sortTorrentListByLeechers(trrs []*anilibria.TitleTorrent) (_ []*anilibria.TitleTorrent) {
	sort.Slice(trrs, func(i, j int) bool {
		return trrs[i].Leechers > trrs[j].Leechers
	})

	// debug
	for _, trr := range trrs {
		gLog.Debug().Str("torrent_hash", trr.Hash[0:9]).Int64("torrnet_size_mb", trr.TotalSize/1024/1024).
			Int("torrent_leechers", trr.Leechers).Msg("sorted slice debug")
	}

	return trrs
}

func (*deploy) balanceForWorkers(trrs []*anilibria.TitleTorrent) (_ map[string][]anilibria.TitleTorrent, e error) {
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
		return nil, errors.New("there is no avaliable workers for the balancing proccess")
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
		return nil, errors.New("there is no workers with free space for the balancing proccess")
	}

	var wtitles = make(map[string][]anilibria.TitleTorrent)

loop:
	for {

		w := <-blncr
		var assigned bool

		for id, trr := range trrs {
			if trr == nil {
				continue
			}

			if uint64(trr.TotalSize) > fspaces[w] {
				gLog.Info().Str("worker_id", w).Str("torrent_hash", trr.Hash[0:9]).Int64("fspace", int64(fspaces[w])).Int64("tspace", trr.TotalSize).
					Msg("skipping torrents because of insufficient disk space on the worker")
				continue
			}

			// decrease the workers free space
			gLog.Debug().Uint64("old_fspace", fspaces[w]).Int64("torrent_size", trr.TotalSize).
				Uint64("new_fspace", fspaces[w]-uint64(trr.TotalSize)).Msg("decreasing the worker's free space")
			fspaces[w] = fspaces[w] - uint64(trr.TotalSize)

			// assigning the torrent to the worker
			wtitles[w] = append(wtitles[w], *trr)

			// remove the title from a slice
			trrs[id] = nil

			assigned = true
			gLog.Debug().Str("worker_id", w).Str("torrent_hash", trr.Hash[0:9]).Msg("the torrent has been assigned")

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

		gLog.Debug().Msg("there is no avaliable workers for the balancing proccess")
		break loop
	}

	return wtitles, e
}