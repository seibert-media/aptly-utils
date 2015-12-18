package main

import (
	"testing"

	aptly_package_version "github.com/bborbe/aptly_utils/package_version"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_version := aptly_package_version.New()

	err = do(writer, package_version, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
