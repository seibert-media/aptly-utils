package main

import (
	"testing"

	aptly_package_uploader "github.com/bborbe/aptly/package_uploader"
	aptly_requestbuilder_executor "github.com/bborbe/aptly/requestbuilder_executor"
	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	requestbuilder_executor := aptly_requestbuilder_executor.New(client.GetClientWithoutProxy())
	package_uploader := aptly_package_uploader.New(requestbuilder_executor, requestbuilder.NewHttpRequestBuilderProvider())

	err = do(writer, package_uploader, "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
