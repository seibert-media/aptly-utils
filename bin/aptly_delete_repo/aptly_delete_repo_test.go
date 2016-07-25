package main

import (
	"testing"

	aptly_repo_deleter "github.com/bborbe/aptly_utils/repo_deleter"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")

	repo_deleter := aptly_repo_deleter.New(nil, nil, nil)

	err = do(writer, repo_deleter, "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
