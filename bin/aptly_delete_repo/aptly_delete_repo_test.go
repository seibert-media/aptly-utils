package main

import (
	"testing"

	aptly_repo_deleter "github.com/bborbe/aptly_utils/repo_deleter"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	repo_deleter := aptly_repo_deleter.New(nil, nil, nil)
	err = do(repo_deleter, "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
