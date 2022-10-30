package app

import (
	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
)

func (m *deploy) deployFromAniApi(dtype deployType, dryrun ...bool) (aobjects []*deploymentObject, e error) {
	wtorrents, ok := m.getWorkersTorrents()
	if !ok {
		return nil, errFailedWorker
	} else if len(wtorrents) == 0 {
		return nil, errNoWorkers
	}

	var titles []*anilibria.TitleTorrent
	if titles, e = m.getAnilibriaTorrents(dtype); e != nil {
		return
	}

	var dtorrents []*deluge.Torrent
	for _, worker := range wtorrents {
		dtorrents = append(dtorrents, worker.torrents...)
	}

	if aobjects = m.getTorrentsDifference(titles, dtorrents); len(aobjects) == 0 {
		return nil, errNothingDeploy
	}
	m.sortTorrentsByLeechers(aobjects)

	if ok, e = m.assignTorrentsToWorkers(aobjects); e != nil {
		return
	} else if !ok {
		return nil, errNothingAssigned
	}

	// panic avoid
	dryrun = append(dryrun, true)
	if !dryrun[0] {
		m.assignTorrentsToWorkers(aobjects)
	}

	return
}

func (*deploy) getTorrentsDifference(atorrents []*anilibria.TitleTorrent, dtorrents []*deluge.Torrent) (aobjects []*deploymentObject) {
	for _, atrr := range atorrents {
		var found bool

		for _, dtrr := range dtorrents {
			// stop searching if title's torrent has been found in the deluge trrs list
			if atrr.Hash == dtrr.Hash {
				found = true
				break
			}

			// !! WARNING
			// !! There is no comparing by torrent name!!!
			// !! Some torrents may be removed from the anilibria announces
			// !! and updated their hashes because of title update
			// !! Must be implemented shortly

			// TODO
			// !! Check by name and series ...

			// !! CHECK METADATA.NAME
			// !! CHECK METADATA.NAME
			// !! CHECK METADATA.NAME
		}

		// process next title's torrent if current torrent has been found in the deluge trrs list
		if found {
			continue
		}

		// create & append deploy object if title's torrent has NOT been found
		gLog.Debug().Str("title_name", atrr.GetName()).Str("torrent_hash", atrr.GetShortHash()).Int64("torrent_size", atrr.TotalSize).
			Msg("there is deployment candidate found")

		aobjects = append(aobjects, &deploymentObject{
			aniTorrent: atrr,
		})
	}

	return aobjects
}
