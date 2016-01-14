package package_detail

import (
	aptly_package_name "github.com/bborbe/aptly_utils/package_name"
	aptly_version "github.com/bborbe/aptly_utils/version"
)

type PackageDetail struct {
	PackageName aptly_package_name.PackageName
	Version     aptly_version.Version
}

func FromInfo(info map[string]string) PackageDetail {
	return PackageDetail{
		PackageName: aptly_package_name.PackageName(info["Package"]),
		Version:     aptly_version.Version(info["Version"]),
	}
}
