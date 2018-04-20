package package_uploader

import (
	"bytes"
	"net/http"
	"testing"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_requestbuilder_executor "github.com/seibert-media/aptly-utils/requestbuilder_executor"
	. "github.com/bborbe/assert"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

func TestImplementsPackageUploader(t *testing.T) {
	b := New(nil, nil, nil)
	var i *PackageUploader
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
func TestFromFileNameWithoutSlash(t *testing.T) {
	name := FromFileName("foo.deb")
	if err := AssertThat(name, Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}

func TestFromFileNameWithSlash(t *testing.T) {
	name := FromFileName("asdf/foo.deb")
	if err := AssertThat(name, Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile(t *testing.T) {
	filecontent := "hello"
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	counter := 0
	requestbuilder_executor := aptly_requestbuilder_executor.New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		if err := AssertThat(req.ContentLength, Is(int64(245))); err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: 200,
		}, nil
	})
	uploader := New(requestbuilder_executor, httpRequestBuilderProvider, nil)
	reader := bytes.NewBufferString(filecontent)
	err := uploader.uploadFile(aptly_model.API{}, "filename", reader)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
