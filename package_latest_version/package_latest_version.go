package package_latest_version

import (
	"fmt"
	"sort"

	aptly_api "github.com/bborbe/aptly_utils/api"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

var logger = log.DefaultLogger

type PackageVersions func(api aptly_api.Api, repository aptly_repository.Repository, packageName package_name.PackageName) ([]aptly_version.Version, error)

type PackageLatestVersion interface {
	PackageLatestVersion(api aptly_api.Api, repository aptly_repository.Repository, packageName package_name.PackageName) (*aptly_version.Version, error)
}

type packageLatestVersion struct {
	packageVersions PackageVersions
}

func New(packageVersions PackageVersions) *packageLatestVersion {
	p := new(packageLatestVersion)
	p.packageVersions = packageVersions
	return p
}

func (p *packageLatestVersion) PackageLatestVersion(api aptly_api.Api, repository aptly_repository.Repository, packageName package_name.PackageName) (*aptly_version.Version, error) {
	logger.Debugf("PackageLatestVersion")
	var err error
	var versions []aptly_version.Version
	if versions, err = p.packageVersions(api, repository, packageName); err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("package %s not found", packageName)
	}
	sort.Sort(aptly_version.VersionByName(versions))
	latestVersion := versions[len(versions)-1]
	logger.Debugf("found latest version %v for package %s", latestVersion, packageName)
	return &latestVersion, nil
}
