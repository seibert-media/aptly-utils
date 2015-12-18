package main

import (
	"testing"

	aptly_package_versions "github.com/bborbe/aptly_utils/package_versions"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_versions := aptly_package_versions.New(nil)

	err = do(writer, package_versions, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
