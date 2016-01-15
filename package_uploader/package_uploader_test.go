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
