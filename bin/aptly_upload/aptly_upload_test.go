package main

import (
	"testing"

	"net/http"

	aptly_package_uploader "github.com/bborbe/aptly/package_uploader"
	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()
	package_uploader := aptly_package_uploader.New(func() *http.Client {
		return client.GetClientWithoutProxy()
	}, requestbuilder.NewHttpRequestBuilderProvider().NewHttpRequestBuilder)

	err = do(writer, package_uploader, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
