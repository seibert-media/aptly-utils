package main

import (
	"testing"

	aptly_package_deleter "github.com/seibert-media/aptly-utils/package_deleter"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	package_deleter := aptly_package_deleter.New(nil, nil, nil)
	err = do(package_deleter, "", "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
