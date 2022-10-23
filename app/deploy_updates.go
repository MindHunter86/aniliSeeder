package app

import (
	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
)

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

		if !found {
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
