package swarm

type Minion struct{}

func NewMinion() *Minion {
	return &Minion{}
}

func (*Minion) Bootstrap() error {
	return nil
}
