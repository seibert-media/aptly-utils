package package_detail_lister

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageDetailLister(t *testing.T) {
	b := New(nil)
	var i *PackageDetailLister
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
