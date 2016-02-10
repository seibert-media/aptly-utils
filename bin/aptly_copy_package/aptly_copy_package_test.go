package main

import (
	"testing"

	"bytes"

	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_package_detail_lister "github.com/bborbe/aptly_utils/package_detail_latest_lister"
	aptly_package_latest_version "github.com/bborbe/aptly_utils/package_latest_version"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")

	package_copier := aptly_package_copier.New(nil, nil, nil)
	packageLastestVersion := aptly_package_latest_version.New(nil)
	packageDetailLister := aptly_package_detail_lister.New(nil)

	err = do(writer, package_copier, packageLastestVersion, packageDetailLister, "", "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
