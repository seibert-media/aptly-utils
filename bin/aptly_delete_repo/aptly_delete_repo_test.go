package main

import (
	"testing"

	aptly_repo_deleter "github.com/bborbe/aptly_utils/repo_deleter"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	repo_deleter := aptly_repo_deleter.New(nil, nil, nil)

	err = do(writer, repo_deleter, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
