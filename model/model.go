package model

import (
	"bytes"
	aptly_version "github.com/bborbe/version"
)

const (
	DISTRIBUTION_DEFAULT = Distribution("default")
	ARCHITECTURE_ALL     = Architecture("all")
	ARCHITECTURE_I386    = Architecture("i386")
	ARCHITECTURE_AMD64   = Architecture("amd64")
	ARCHITECTURE_DEFAULT = ARCHITECTURE_AMD64
	PACKAGE_ALL          = Package("all")
)

type Api struct {
	RepoUrl     RepoUrl
	ApiUrl      ApiUrl
	ApiUsername ApiUsername
	ApiPassword ApiPassword
}

type Package string

type Architecture string

type Key string

type Distribution string

type RepoUrl string

type ApiUrl string

type ApiUsername string

type ApiPassword string

type Repository string

type PackageDetail struct {
	Package Package
	Version aptly_version.Version
}

func NewPackageDetail(packageName Package, version aptly_version.Version) PackageDetail {
	return PackageDetail{
		Package: packageName, Version: version}
}

func NewPackageDetailByString(packageName string, version string) PackageDetail {
	return NewPackageDetail(Package(packageName), aptly_version.Version(version))
}

func FromInfo(info map[string]string) PackageDetail {
	return NewPackageDetailByString(info["Package"], info["Version"])
}

func FromInfos(infos []map[string]string) []PackageDetail {
	var result []PackageDetail
	for _, info := range infos {
		result = append(result, FromInfo(info))
	}
	return result
}

func NewApi(
	repoUrl string,
	apiUrl string,
	apiUsername string,
	apiPassword string,
) Api {
	return Api{
		RepoUrl:     RepoUrl(repoUrl),
		ApiUrl:      ApiUrl(apiUrl),
		ApiUsername: ApiUsername(apiUsername),
		ApiPassword: ApiPassword(apiPassword),
	}
}

func JoinArchitectures(architectures []Architecture, sep string) string {
	b := bytes.NewBufferString("")
	first := true
	for _, a := range architectures {
		if first {
			first = false
		} else {
			b.WriteString(sep)
		}
		b.WriteString(string(a))
	}
	return string(b.Bytes())
}

func ParseArchitectures(architectures []string) []Architecture {
	var result []Architecture
	for _, name := range architectures {
		result = append(result, Architecture(name))
	}
	return result
}
