package package_versions

import (
	"github.com/bborbe/log"

	aptly_version "github.com/bborbe/aptly_utils/version"
)

type ListPackages func(url string, user string, password string, repo string) ([]map[string]string, error)

type PackageVersions interface {
	PackageVersions(url string, user string, password string, repo string, name string) ([]aptly_version.Version, error)
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

func (p *packageVersion) PackageVersions(url string, user string, password string, repo string, name string) ([]aptly_version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repo, name)
	jsonStruct, err := p.listPackages(url, user, password, repo)
	if err != nil {
		return nil, err
	}
	var versions []aptly_version.Version
	for _, info := range jsonStruct {
		if info["Package"] == name {
			v := info["Version"]
			logger.Debugf("found version: %s", v)
			versions = append(versions, aptly_version.Version(v))
		}
	}
	return versions, nil
}
