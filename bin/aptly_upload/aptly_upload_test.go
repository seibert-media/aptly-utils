package main

import (
	"testing"

	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_uploader := aptly_package_uploader.New(nil, nil, nil)

	err = do(writer, package_uploader, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
