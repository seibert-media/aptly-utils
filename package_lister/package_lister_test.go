package package_versions

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageLister(t *testing.T) {
	b := New(nil, nil)
	var i *PackageLister
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
