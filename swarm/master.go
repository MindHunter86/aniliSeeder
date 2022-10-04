package swarm

type Master struct{}

func NewMaster() *Master {
	return &Master{}
}

func (*Master) Bootstrap() error {
	return nil
}
