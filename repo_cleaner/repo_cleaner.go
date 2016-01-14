package repo_deleter

import (
	aptly_key "github.com/bborbe/aptly_utils/key"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

type DeletePackagesByKey func(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	keys []aptly_key.Key) error

type ListPackages func(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository) ([]map[string]string, error)

type RepoCleaner interface {
	CleanRepo(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository) error
}

type repoCleaner struct {
	deletePackagesByKey DeletePackagesByKey
	listPackages        ListPackages
}

var logger = log.DefaultLogger

func New(deletePackagesByKey DeletePackagesByKey, listPackages ListPackages) *repoCleaner {
	r := new(repoCleaner)
	r.deletePackagesByKey = deletePackagesByKey
	r.listPackages = listPackages
	return r
}

func (r *repoCleaner) CleanRepo(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository) error {
	logger.Debugf("clean repo: %s", repository)
	keys, err := r.findKeysToDelete(url, user, password, repository)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		logger.Debugf("nothing to delete")
		return nil
	}
	return r.deletePackagesByKey(url, user, password, repository, keys)
}

func (r *repoCleaner) findKeysToDelete(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository) ([]aptly_key.Key, error) {
	logger.Debugf("find keys to delete repo: %s", repository)
	packages, err := r.listPackages(url, user, password, repository)
	if err != nil {
		return nil, err
	}
	return packagesToKeys(packages), nil
}

func packagesToKeys(packages []map[string]string) []aptly_key.Key {
	latestVersions := make(map[string]map[string]string)
	var keys []aptly_key.Key
	for _, currentPackage := range packages {
		name := currentPackage["Package"]
		if latestPackage, ok := latestVersions[name]; ok {
			var packageToDelete map[string]string
			logger.Tracef("compare %s < %s", currentPackage["Version"], latestPackage["Version"])
			if aptly_version.Less(aptly_version.Version(currentPackage["Version"]), aptly_version.Version(latestPackage["Version"])) {
				packageToDelete = currentPackage
			} else {
				logger.Tracef("set latest version %s %s", currentPackage["Package"], currentPackage["Version"])
				latestVersions[name] = currentPackage
				packageToDelete = latestPackage
			}
			keys = append(keys, aptly_key.Key(packageToDelete["Key"]))
			logger.Debugf("mark package %s %s to delete", packageToDelete["Package"], packageToDelete["Version"])
		} else {
			latestVersions[name] = currentPackage
		}
	}
	for _, p := range latestVersions {
		logger.Debugf("keep package %s %s", p["Package"], p["Version"])
	}
	return keys
}
