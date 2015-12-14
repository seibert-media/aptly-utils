package main

import (
	"testing"

	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client "github.com/bborbe/http/client"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	client := http_client.GetClientWithoutProxy()
	requestbuilder_executor := aptly_requestbuilder_executor.New(client)
	requestbuilder := http_requestbuilder.NewHttpRequestBuilderProvider()
	package_uploader := aptly_package_uploader.New(requestbuilder_executor, requestbuilder)
	package_copier := aptly_package_copier.New(package_uploader, requestbuilder, client)

	err = do(writer, package_copier, "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
