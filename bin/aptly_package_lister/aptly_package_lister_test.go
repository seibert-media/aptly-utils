package main

import (
	"testing"

	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()
	package_lister := aptly_package_lister.New(nil, nil)
	err = do(writer, package_lister, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
