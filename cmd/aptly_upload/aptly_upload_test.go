package main

import (
	"testing"

	aptly_package_uploader "github.com/seibert-media/aptly-utils/package_uploader"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	package_uploader := aptly_package_uploader.New(nil, nil, nil)
	err := do(package_uploader, "", "", "", "", "", "", "", "")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
