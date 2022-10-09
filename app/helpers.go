package app

import (
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

func (*deploy) getAnilibriaTorrents() (e error) {
	return
}

func (*deploy) getWorkersTorrents() (e error) {
	return
}

func (*deploy) sortTorrentListByLeechers() {

}

func (*deploy) compareTorrentLists() {

}

func (*deploy) balanceForWorkers() {

}
