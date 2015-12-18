package repo_deleter

type RepoCleaner interface {
	CleanRepo() error
}

type repoCleaner struct {
}

func New() *repoCleaner {
	return new(repoCleaner)
}

func (c *repoCleaner) CleanRepo() error {
	return nil
}
