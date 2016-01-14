package package_versions

import (
	aptly_api "github.com/bborbe/aptly_utils/api"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

type ListPackages func(
	api aptly_api.Api,
	repository aptly_repository.Repository) ([]map[string]string, error)

type PackageVersions interface {
	PackageVersions(
		api aptly_api.Api,
		repository aptly_repository.Repository,
		packageName package_name.PackageName) ([]aptly_version.Version, error)
}

type packageVersion struct {
	listPackages ListPackages
}

var logger = log.DefaultLogger

func New(listPackages ListPackages) *packageVersion {
	p := new(packageVersion)
	p.listPackages = listPackages
	return p
}

type JsonStruct []map[string]string

func (p *packageVersion) PackageVersions(
	api aptly_api.Api,
	repository aptly_repository.Repository,
	packageName package_name.PackageName) ([]aptly_version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repository, packageName)
	jsonStruct, err := p.listPackages(api, repository)
	if err != nil {
		return nil, err
	}
	var versions []aptly_version.Version
	for _, info := range jsonStruct {
		if info["Package"] == string(packageName) {
			v := info["Version"]
			logger.Debugf("found version: %s", v)
			versions = append(versions, aptly_version.Version(v))
		}
	}
	return versions, nil
}
