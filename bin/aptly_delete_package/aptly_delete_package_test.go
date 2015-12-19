package main

import (
	"testing"

	aptly_package_deleter "github.com/bborbe/aptly_utils/package_deleter"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_deleter := aptly_package_deleter.New(nil, nil, nil)

	err = do(writer, package_deleter, "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
