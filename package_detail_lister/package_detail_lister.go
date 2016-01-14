package package_detail_lister

import (
	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_package_detail "github.com/bborbe/aptly_utils/package_detail"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type ListPackages func(api aptly_api.Api, repository aptly_repository.Repository) ([]map[string]string, error)

type PackageDetailLister interface {
	ListPackageDetails(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error)
}

type packageDetailLister struct {
	listPackages ListPackages
}

func New(listPackages ListPackages) *packageDetailLister {
	p := new(packageDetailLister)
	p.listPackages = listPackages
	return p
}

func (p *packageDetailLister) ListPackageDetails(api aptly_api.Api, repository aptly_repository.Repository) ([]aptly_package_detail.PackageDetail, error) {
	logger.Debugf("ListPackageDetails")
	infos, err := p.listPackages(api, repository)
	if err != nil {
		return nil, err
	}
	return aptly_package_detail.FromInfos(infos), nil
}
