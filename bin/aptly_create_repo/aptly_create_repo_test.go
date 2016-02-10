package main

import (
	"testing"

	"bytes"

	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")

	repo_creator := aptly_repo_creater.New(nil, nil, nil)

	err = do(writer, repo_creator, "", "", "", "", "", "", nil)
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
