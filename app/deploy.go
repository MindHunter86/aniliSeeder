package app

import (
	"sort"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
)

type deployObject struct {
}

// ?? func (*App) getDeloyObjectsFromAniDel()
func (*App) compareDelugeTorrentsWithAniTitles(dtrrs []*deluge.Torrent, attl []*anilibria.Title) []*deployObject {
	return nil
}

type deploy struct {
}

type deployType uint8

const (
	dplAnilibriaUpdates deployType = iota
	// dplAnilibriaSchedule
	// dplAnilibriaWatchingNow
)

func newDeploy() *deploy {
	return &deploy{}
}

func (m *deploy) run() error {
	return m.deploy(false)
}

func (m *deploy) dryRun() error {
	return m.deploy(true)
}

func (*deploy) deploy(isDryRun bool) (e error) {
	return
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
// !! There is no compring by torrent name!!!
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

	return
}

func (*deploy) sortTorrentListByLeechers(trrs []*anilibria.TitleTorrent) (_ []*anilibria.TitleTorrent) {
	sort.Slice(trrs, func(i, j int) bool {
		return trrs[i].Leechers < trrs[j].Leechers
	})

	return trrs
}

func (*deploy) balanceForWorkers(trrs []*anilibria.TitleTorrent) (_ map[string][]*anilibria.TitleTorrent, e error) {
	wrks := gSwarm.GetConnectedWorkers()
	blncr := make(chan string, len(wrks))

	var fspaces map[string]uint64
	for id := range wrks {
		if fspaces[id] == 0 {
			if fspaces[id], e = gSwarm.RequestFreeSpaceFromWorker(id); e != nil {
				gLog.Error().Err(e).Str("worker_id", id).Msg("got an error in free space request to the worker")
				continue
			}
		}

		blncr <- id
	}

	// for _, trr := range trrs {
	// 	if len(blncr) == 0 {
	// 		break
	// 	}

	// 	pid := <-blncr
	// }

	return
}
