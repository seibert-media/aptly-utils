package package_copier

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageCopier(t *testing.T) {
	b := New(nil, nil)
	var i *PackageCopier
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
