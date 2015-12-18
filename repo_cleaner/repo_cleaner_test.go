package repo_deleter

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoCleaner(t *testing.T) {
	b := New()
	var i *RepoCleaner
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
