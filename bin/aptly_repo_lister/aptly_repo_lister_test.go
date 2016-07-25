package main

import (
	"testing"

	"bytes"

	aptly_repo_lister "github.com/bborbe/aptly_utils/repo_lister"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")
	repo_lister := aptly_repo_lister.New(nil, nil)
	err = do(writer, repo_lister, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
