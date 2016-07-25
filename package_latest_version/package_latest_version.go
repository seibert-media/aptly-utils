package package_latest_version

import (
	"fmt"
	"sort"

	"github.com/bborbe/aptly_utils/model"
	aptly_model "github.com/bborbe/aptly_utils/model"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

var logger = log.DefaultLogger

type PackageVersions func(api aptly_model.Api, repository aptly_model.Repository, packageName model.Package) ([]aptly_version.Version, error)

type PackageLatestVersion interface {
	PackageLatestVersion(api aptly_model.Api, repository aptly_model.Repository, packageName model.Package) (*aptly_version.Version, error)
}

type packageLatestVersion struct {
	packageVersions PackageVersions
}

func New(packageVersions PackageVersions) *packageLatestVersion {
	p := new(packageLatestVersion)
	p.packageVersions = packageVersions
	return p
}

func (p *packageLatestVersion) PackageLatestVersion(api aptly_model.Api, repository aptly_model.Repository, packageName model.Package) (*aptly_version.Version, error) {
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
