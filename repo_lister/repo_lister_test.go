package repo_lister

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoLister(t *testing.T) {
	b := New(nil, nil)
	var i *RepoLister
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
