package repo_deleter

import (
	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_key "github.com/bborbe/aptly_utils/key"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_version "github.com/bborbe/version"
	"github.com/bborbe/log"
)

type DeletePackagesByKey func(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, keys []aptly_key.Key) error

type ListPackages func(api aptly_api.Api, repository aptly_repository.Repository) ([]map[string]string, error)

type RepoCleaner interface {
	CleanRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error
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

func (r *repoCleaner) CleanRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error {
	logger.Debugf("clean repo: %s", repository)
	keys, err := r.findKeysToDelete(api, repository)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		logger.Debugf("nothing to delete")
		return nil
	}
	return r.deletePackagesByKey(api, repository, distribution, keys)
}

func (r *repoCleaner) findKeysToDelete(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_key.Key, error) {
	logger.Debugf("find keys to delete repo: %s", repository)
	packages, err := r.listPackages(api, repository)
	if err != nil {
		return nil, err
	}
	return packagesToKeys(packages), nil
}

func packagesToKeys(packages []map[string]string) []aptly_key.Key {
	latestVersions := make(map[string]map[string]string)
	var keys []aptly_key.Key
	for _, currentPackage := range packages {
		logger.Debugf("handle package %s %s", currentPackage["Package"], currentPackage["Version"])
		name := currentPackage["Package"]
		if latestPackage, ok := latestVersions[name]; ok {
			var packageToDelete map[string]string
			logger.Tracef("compare %s < %s", currentPackage["Version"], latestPackage["Version"])
			if aptly_version.LessThan(aptly_version.Version(currentPackage["Version"]), aptly_version.Version(latestPackage["Version"])) {
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
