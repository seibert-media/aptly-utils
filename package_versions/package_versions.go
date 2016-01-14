package package_versions

import (
	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_package_detail "github.com/bborbe/aptly_utils/package_detail"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

type ListPackageDetails func(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error)

type PackageVersions interface {
	PackageVersions(api aptly_api.Api, repository aptly_repository.Repository, packageName package_name.PackageName) ([]aptly_version.Version, error)
}

type packageVersion struct {
	listPackageDetails ListPackageDetails
}

var logger = log.DefaultLogger

func New(listPackages ListPackageDetails) *packageVersion {
	p := new(packageVersion)
	p.listPackageDetails = listPackages
	return p
}

func (p *packageVersion) PackageVersions(api aptly_api.Api, repository aptly_repository.Repository, packageName package_name.PackageName) ([]aptly_version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repository, packageName)
	packageDetails, err := p.listPackageDetails(api, repository)
	if err != nil {
		return nil, err
	}
	var versions []aptly_version.Version
	for _, packageDetail := range packageDetails {
		if packageDetail.PackageName == packageName {
			v := packageDetail.Version
			logger.Debugf("found version: %s", v)
			versions = append(versions, aptly_version.Version(v))
		}
	}
	return versions, nil
}
