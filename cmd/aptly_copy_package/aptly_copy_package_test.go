package main

import (
	"testing"

	aptly_package_copier "github.com/seibert-media/aptly-utils/package_copier"
	aptly_model_lister "github.com/seibert-media/aptly-utils/package_detail_latest_lister"
	aptly_package_latest_version "github.com/seibert-media/aptly-utils/package_latest_version"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	package_copier := aptly_package_copier.New(nil, nil, nil)
	packageLastestVersion := aptly_package_latest_version.New(nil)
	packageDetailLister := aptly_model_lister.New(nil)

	err = do(package_copier, packageLastestVersion, packageDetailLister, "", "", "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
