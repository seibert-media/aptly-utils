package package_versions

import (
	"github.com/seibert-media/aptly-utils/model"
	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

type ListPackageDetails func(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error)

type PackageVersions interface {
	PackageVersions(api aptly_model.API, repository aptly_model.Repository, packageName model.Package) ([]aptly_version.Version, error)
}

type packageVersion struct {
	listPackageDetails ListPackageDetails
}

func New(listPackages ListPackageDetails) *packageVersion {
	p := new(packageVersion)
	p.listPackageDetails = listPackages
	return p
}

func (p *packageVersion) PackageVersions(api aptly_model.API, repository aptly_model.Repository, packageName model.Package) ([]aptly_version.Version, error) {
	glog.V(2).Infof("PackageVersions - repo: %s package: %s", repository, packageName)
	packageDetails, err := p.listPackageDetails(api, repository)
	if err != nil {
		return nil, err
	}
	var versions []aptly_version.Version
	for _, packageDetail := range packageDetails {
		if packageDetail.Package == packageName {
			v := packageDetail.Version
			glog.V(2).Infof("found version: %s", v)
			versions = append(versions, aptly_version.Version(v))
		}
	}
	return versions, nil
}
