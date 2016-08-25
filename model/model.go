package model

import (
	"bytes"

	aptly_version "github.com/bborbe/version"
)

const (
	DistribuionDefault  = Distribution("default")
	ArchitectureALL     = Architecture("all")
	ArchitectureI386    = Architecture("i386")
	ArchitectureAMD64   = Architecture("amd64")
	ArchitectureDefault = ArchitectureAMD64
	PackageAll          = Package("all")
)

type API struct {
	RepoURL     RepoURL
	APIUrl      APIUrl
	APIUsername APIUsername
	APIPassword APIPassword
}

type Package string

type Architecture string

type Key string

type Distribution string

type RepoURL string

type APIUrl string

type APIUsername string

type APIPassword string

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

func NewAPI(
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
) API {
	return API{
		RepoURL:     RepoURL(repoURL),
		APIUrl:      APIUrl(apiURL),
		APIUsername: APIUsername(apiUsername),
		APIPassword: APIPassword(apiPassword),
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
