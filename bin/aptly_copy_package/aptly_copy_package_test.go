package main

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
	io_mock "github.com/bborbe/io/mock"
	aptly_package_copier "github.com/bborbe/aptly/package_copier"
	aptly_requestbuilder_executor "github.com/bborbe/aptly/requestbuilder_executor"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	requestbuilder_executor := aptly_requestbuilder_executor.New(client.GetClientWithoutProxy())
	package_copier := aptly_package_copier.New(requestbuilder_executor, requestbuilder.NewHttpRequestBuilderProvider())

	err = do(writer, package_copier, "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
