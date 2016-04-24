package package_detail

import (
	aptly_package_name "github.com/bborbe/aptly_utils/package_name"
	aptly_version "github.com/bborbe/version"
)

type PackageDetail struct {
	PackageName aptly_package_name.PackageName
	Version     aptly_version.Version
}

func New(packageName aptly_package_name.PackageName, version aptly_version.Version) PackageDetail {
	return PackageDetail{
		PackageName: packageName, Version: version}
}

func NewByString(packageName string, version string) PackageDetail {
	return New(aptly_package_name.PackageName(packageName), aptly_version.Version(version))
}

func FromInfo(info map[string]string) PackageDetail {
	return NewByString(info["Package"], info["Version"])
}

func FromInfos(infos []map[string]string) []PackageDetail {
	var result []PackageDetail
	for _, info := range infos {
		result = append(result, FromInfo(info))
	}
	return result
}
