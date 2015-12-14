package main

import (
	"testing"

	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	requestbuilder_executor := aptly_requestbuilder_executor.New(client.GetClientWithoutProxy())
	repo_creator := aptly_repo_creater.New(requestbuilder_executor, requestbuilder.NewHttpRequestBuilderProvider())

	err = do(writer, repo_creator, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
