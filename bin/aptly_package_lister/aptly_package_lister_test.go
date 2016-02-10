package main

import (
	"testing"

	"bytes"

	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")
	package_lister := aptly_package_lister.New(nil, nil)
	err = do(writer, package_lister, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
