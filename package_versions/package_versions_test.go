package package_versions

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageVersions(t *testing.T) {
	b := New(nil, nil)
	var i *PackageVersions
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
