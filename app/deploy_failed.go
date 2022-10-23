package app

func (m *deploy) dryDeployFailedAnnounces() error {
	return m.deployFailedAnnounces(true)
}

func (*deploy) deployFailedAnnounces(dryrun ...bool) (e error) {
	return
}

// func (*deploy) getWorkersTorrents() {
// }

func (*deploy) getFailedAnnounces() {
}

func (*deploy) searchFailedTorrents() {

}

func (*deploy) isSpaceEnough() {

}

func (*deploy) dropFailedTorrent() {

}
