package repo_deleter

import "github.com/bborbe/log"

type DeletePackagesByKey func(url string, user string, password string, repo string, keys []string) error

type RepoCleaner interface {
	CleanRepo(url string, user string, password string, repo string) error
}

type repoCleaner struct {
	deletePackagesByKey DeletePackagesByKey
}

var logger = log.DefaultLogger

func New(deletePackagesByKey DeletePackagesByKey) *repoCleaner {
	r := new(repoCleaner)
	r.deletePackagesByKey = deletePackagesByKey
	return r
}

func (r *repoCleaner) CleanRepo(url string, user string, password string, repo string) error {
	logger.Debugf("clean repo: %s", repo)
	keys, err := r.findKeysToDelete(url, user, password, repo)
	if err != nil {
		return err
	}
	return r.deletePackagesByKey(url, user, password, repo, keys)
}

func (r *repoCleaner) findKeysToDelete(url string, user string, password string, repo string) ([]string, error) {
	logger.Debugf("find keys to delete repo: %s", repo)
	return nil, nil
}
