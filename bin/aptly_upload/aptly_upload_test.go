package main

import (
	"testing"

	"net/http"
	. "github.com/bborbe/assert"
	"github.com/bborbe/http/client"
	io_mock "github.com/bborbe/io/mock"
	aptly_package_uploader "github.com/bborbe/aptly/package_uploader"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()
	package_uploader := aptly_package_uploader.New(func() *http.Client {
		return client.GetClientWithoutProxy()
	})

	err = do(writer, package_uploader, "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
