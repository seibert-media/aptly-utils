package repo_deleter

type RepoDeleter interface {
	DeleteRepo() error
}

type repoDeleter struct {

}

func New() *repoDeleter {
	return new(repoDeleter)
}

func (c *repoDeleter ) DeleteRepo() error {
	return nil
}
