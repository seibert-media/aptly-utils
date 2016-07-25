package package_versions

import (
	"github.com/bborbe/aptly_utils/model"
	aptly_model "github.com/bborbe/aptly_utils/model"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

type ListPackageDetails func(api aptly_model.Api, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error)

type PackageVersions interface {
	PackageVersions(api aptly_model.Api, repository aptly_model.Repository, packageName model.Package) ([]aptly_version.Version, error)
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

func (p *packageVersion) PackageVersions(api aptly_model.Api, repository aptly_model.Repository, packageName model.Package) ([]aptly_version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repository, packageName)
	packageDetails, err := p.listPackageDetails(api, repository)
	if err != nil {
		return nil, err
	}
	var versions []aptly_version.Version
	for _, packageDetail := range packageDetails {
		if packageDetail.Package == packageName {
			v := packageDetail.Version
			logger.Debugf("found version: %s", v)
			versions = append(versions, aptly_version.Version(v))
		}
	}
	return versions, nil
}
