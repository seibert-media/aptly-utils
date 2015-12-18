package package_deleter

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageDeleter(t *testing.T) {
	b := New(nil, nil)
	var i *PackageDeleter
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
