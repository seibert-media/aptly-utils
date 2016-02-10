package main

import (
	"testing"

	aptly_repo_cleaner "github.com/bborbe/aptly_utils/repo_cleaner"

	"bytes"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")
	repo_cleaner := aptly_repo_cleaner.New(nil, nil)
	err = do(writer, repo_cleaner, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
