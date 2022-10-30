package app

import (
	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/MindHunter86/aniliSeeder/utils"
)

func (m *deploy) run() (map[string][]*anilibria.TitleTorrent, error) {
	return m.deploy(false)
}

func (m *deploy) dryRun() (map[string][]*anilibria.TitleTorrent, error) {
	return m.deploy(true)
}

func (m *deploy) deploy(isDryRun bool) (_ map[string][]*anilibria.TitleTorrent, e error) {
	var titles []*anilibria.TitleTorrent
	if titles, e = m.getAnilibriaUpdatesTorrents(); e != nil {
		return
	}

	var torrents []*deluge.Torrent
	if torrents, e = m.getWorkersTorrents(); e != nil {
		return
	}

	titleUpdates := m.compareUpdateListWithTorrents(titles, torrents)
	if len(titleUpdates) == 0 {
		return nil, errNothingDeploy
	}

	sortedUpdates := m.sortTorrentListByLeechers(titleUpdates)

	var assignedTitles = make(map[string][]*anilibria.TitleTorrent)
	if assignedTitles, e = m.balanceForWorkers(sortedUpdates); e != nil {
		return
	}

	if len(assignedTitles) == 0 {
		return nil, errNothingAssigned
	}

	if !isDryRun {
		m.sendDeployCommand(assignedTitles)
	}

	return assignedTitles, e
}

func (*deploy) getAnilibriaUpdatesTorrents() (trrs []*anilibria.TitleTorrent, e error) {
	var ttls []*anilibria.Title
	if ttls, e = gAniApi.GetLastUpdates(); e != nil {
		return
	}

	for _, ttl := range ttls {
		for _, trr := range ttl.Torrents.List {
			if tsize := utils.GetMBytesFromBytes(trr.TotalSize); tsize > int64(gCli.Uint64("anilibria-max-torrent-size")) {
				gLog.Info().Str("title_name", ttl.Names.En).Str("torrent_hash", trr.GetShortHash()).Int64("torrent_size_mb", tsize).
					Int("title_id", ttl.Id).Uint64("download_limit", gCli.Uint64("anilibria-max-torrent-size")).
					Msg("skipping a torrent because the torrent is larger than the download limit...")
				continue
			}

			trrs = append(trrs, trr)
		}

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
		gLog.Debug().Str("torrent_hash", trr.GetShortHash()).Int64("torrent_size", trr.TotalSize).Msg("there is a new item for the further deploy")
	}

	return
}
