package package_uploader

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageUploader(t *testing.T) {
	b := New(nil, nil, nil)
	var i *PackageUploader
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtractNameOfFileWithoutSlash(t *testing.T) {
	name := extractPkgOfFile("foo.deb")
	if err := AssertThat(name, Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}

func TestExtractNameOfFileWithSlash(t *testing.T) {
	name := extractPkgOfFile("asdf/foo.deb")
	if err := AssertThat(name, Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}
