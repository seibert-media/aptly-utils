package main

import (
	"testing"

	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_package_latest_version "github.com/bborbe/aptly_utils/package_latest_version"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_copier := aptly_package_copier.New(nil, nil, nil)
	packageLastestVersion := aptly_package_latest_version.New(nil)
	packageLister := aptly_package_lister.New(nil, nil)

	err = do(writer, package_copier, packageLastestVersion, packageLister, "", "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
