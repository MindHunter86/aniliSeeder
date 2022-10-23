package app

import (
	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
)

func (m *deploy) dryDeployFailedAnnounces() error {
	return m.deployFailedAnnounces(true)
}

func (*deploy) deployFailedAnnounces(dryrun ...bool) (e error) {
	return
}

// func (*deploy) getWorkersTorrents() {
// }

func (*deploy) getFailedAnnounces(trrs []*deluge.Torrent) []*deluge.Torrent {
	var ftorrents []*deluge.Torrent

	for _, trr := range trrs {
		if trr.IsTrackerOk() {
			continue
		}

		gLog.Debug().Str("torrents_hash", trr.Hash[0:9]).Msg("found torrents with failed announces")
		ftorrents = append(ftorrents, trr)
	}

	return ftorrents
}

func (*deploy) searchFailedTitles(trrs []*deluge.Torrent) (_ []*anilibria.TitleTorrent, e error) {
	var anitorrents []*anilibria.TitleTorrent

	for _, trr := range trrs {
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

			anitorrents = append(anitorrents, anitrr)
			found = true
		}

		if !found {
			gLog.Warn().Str("torrent_hash", trr.GetShortHash()).Str("title_name", trr.Name).
				Msg("there is a proble in searching title's torrent by quality string; manual search required")
		}
	}

	return anitorrents, e
}

func (*deploy) isSpaceEnough() {

}

func (*deploy) dropFailedTorrent() {

}
