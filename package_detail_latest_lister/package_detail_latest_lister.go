package package_detail_latest_lister

import (
	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

type ListPackageDetails func(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error)

type PackageDetailLatestLister interface {
	ListLatestPackageDetails(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error)
}

type packageDetailLatestLister struct {
	listPackageDetails ListPackageDetails
}

func New(listPackageDetails ListPackageDetails) *packageDetailLatestLister {
	p := new(packageDetailLatestLister)
	p.listPackageDetails = listPackageDetails
	return p
}

func (p *packageDetailLatestLister) ListLatestPackageDetails(api aptly_model.API, repository aptly_model.Repository) ([]aptly_model.PackageDetail, error) {
	glog.V(2).Infof("ListPackageDetails")
	list, err := p.listPackageDetails(api, repository)
	if err != nil {
		return nil, err
	}
	return latest(list...), nil
}

func latest(list ...aptly_model.PackageDetail) []aptly_model.PackageDetail {
	latest := make(map[aptly_model.Package]aptly_version.Version)
	for _, e := range list {
		if val, ok := latest[e.Package]; !ok || aptly_version.LessThan(val, e.Version) {
			latest[e.Package] = e.Version
		}
	}
	var result []aptly_model.PackageDetail
	for k, v := range latest {
		result = append(result, aptly_model.NewPackageDetail(k, v))
	}
	return result
}
