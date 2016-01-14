package package_versions

import (
	"github.com/bborbe/aptly_utils/package_name"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

type ListPackages func(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository) ([]map[string]string, error)

type PackageVersions interface {
	PackageVersions(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
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
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository,
	packageName package_name.PackageName) ([]aptly_version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repository, packageName)
	jsonStruct, err := p.listPackages(apiUrl, apiUsername, apiPassword, repository)
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
