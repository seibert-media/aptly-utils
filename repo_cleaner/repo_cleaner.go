package repo_deleter

import (
	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

type DeletePackagesByKey func(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, keys []aptly_model.Key) error

type ListPackages func(api aptly_model.API, repository aptly_model.Repository) ([]map[string]string, error)

type RepoCleaner interface {
	CleanRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error
}

type repoCleaner struct {
	deletePackagesByKey DeletePackagesByKey
	listPackages        ListPackages
}

func New(deletePackagesByKey DeletePackagesByKey, listPackages ListPackages) *repoCleaner {
	r := new(repoCleaner)
	r.deletePackagesByKey = deletePackagesByKey
	r.listPackages = listPackages
	return r
}

func (r *repoCleaner) CleanRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error {
	glog.V(2).Infof("clean repo: %s", repository)
	keys, err := r.findKeysToDelete(api, repository)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		glog.V(2).Infof("nothing to delete")
		return nil
	}
	return r.deletePackagesByKey(api, repository, distribution, keys)
}

func (r *repoCleaner) findKeysToDelete(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.Key, error) {
	glog.V(2).Infof("find keys to delete repo: %s", repository)
	packages, err := r.listPackages(api, repository)
	if err != nil {
		return nil, err
	}
	return packagesToKeys(packages), nil
}

func packagesToKeys(packages []map[string]string) []aptly_model.Key {
	latestVersions := make(map[string]map[string]string)
	var keys []aptly_model.Key
	for _, currentPackage := range packages {
		glog.V(2).Infof("handle package %s %s", currentPackage["Package"], currentPackage["Version"])
		name := currentPackage["Package"]
		if latestPackage, ok := latestVersions[name]; ok {
			var packageToDelete map[string]string
			glog.V(4).Infof("compare %s < %s", currentPackage["Version"], latestPackage["Version"])
			if aptly_version.LessThan(aptly_version.Version(currentPackage["Version"]), aptly_version.Version(latestPackage["Version"])) {
				packageToDelete = currentPackage
			} else {
				glog.V(4).Infof("set latest version %s %s", currentPackage["Package"], currentPackage["Version"])
				latestVersions[name] = currentPackage
				packageToDelete = latestPackage
			}
			keys = append(keys, aptly_model.Key(packageToDelete["Key"]))
			glog.V(2).Infof("mark package %s %s to delete", packageToDelete["Package"], packageToDelete["Version"])
		} else {
			latestVersions[name] = currentPackage
		}
	}
	for _, p := range latestVersions {
		glog.V(2).Infof("keep package %s %s", p["Package"], p["Version"])
	}
	return keys
}
