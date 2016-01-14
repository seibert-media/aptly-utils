package package_latest_version

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageLatestVersion(t *testing.T) {
	b := New(nil)
	var i *PackageLatestVersion
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
