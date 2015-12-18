package package_deleter

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageVersion(t *testing.T) {
	b := New()
	var i *PackageVersion
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
