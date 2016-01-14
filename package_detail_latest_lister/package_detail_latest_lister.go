package package_detail_latest_lister

import (
	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_package_detail "github.com/bborbe/aptly_utils/package_detail"
	aptly_package_name "github.com/bborbe/aptly_utils/package_name"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type ListPackageDetails func(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error)

type PackageDetailLatestLister interface {
	ListLatestPackageDetails(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error)
}

type packageDetailLatestLister struct {
	listPackageDetails ListPackageDetails
}

func New(listPackageDetails ListPackageDetails) *packageDetailLatestLister {
	p := new(packageDetailLatestLister)
	p.listPackageDetails = listPackageDetails
	return p
}

func (p *packageDetailLatestLister) ListLatestPackageDetails(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error) {
	logger.Debugf("ListPackageDetails")
	list, err := p.listPackageDetails(api, repository)
	if err != nil {
		return nil, err
	}
	return latest(list...), nil
}

func latest(list ...aptly_package_detail.PackageDetail) []aptly_package_detail.PackageDetail {
	latest := make(map[aptly_package_name.PackageName]aptly_version.Version)
	for _, e := range list {
		if val, ok := latest[e.PackageName]; !ok || aptly_version.Less(val, e.Version) {
			latest[e.PackageName] = e.Version
		}
	}
	var result []aptly_package_detail.PackageDetail
	for k, v := range latest {
		result = append(result, aptly_package_detail.New(k, v))
	}
	return result
}
