package package_uploader


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageUploader(t *testing.T) {
	b := New(nil)
	var i *PackageUploader
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}