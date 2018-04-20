package main

import (
	"testing"

	aptly_repo_cleaner "github.com/seibert-media/aptly-utils/repo_cleaner"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	var err error
	repo_cleaner := aptly_repo_cleaner.New(nil, nil)
	err = do(repo_cleaner, "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
