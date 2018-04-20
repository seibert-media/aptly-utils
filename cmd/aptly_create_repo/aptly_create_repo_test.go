package main

import (
	"testing"

	aptly_repo_creater "github.com/seibert-media/aptly-utils/repo_creater"
	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	repo_creator := aptly_repo_creater.New(nil, nil, nil)
	err = do(repo_creator, "", "", "", "", "", "", "", nil)
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
