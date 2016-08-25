package package_detail_lister

import (
	aptly_model "github.com/bborbe/aptly_utils/model"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type ListPackages func(api aptly_model.API, repository aptly_model.Repository) ([]map[string]string, error)

type PackageDetailLister interface {
	ListPackageDetails(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error)
}

type packageDetailLister struct {
	listPackages ListPackages
}

func New(listPackages ListPackages) *packageDetailLister {
	p := new(packageDetailLister)
	p.listPackages = listPackages
	return p
}

func (p *packageDetailLister) ListPackageDetails(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error) {
	logger.Debugf("ListPackageDetails")
	infos, err := p.listPackages(api, repository)
	if err != nil {
		return nil, err
	}
	return aptly_model.FromInfos(infos), nil
}
