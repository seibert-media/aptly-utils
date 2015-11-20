package repo_creator


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoDeleter(t *testing.T) {
	b := New(nil, nil)
	var i *RepoDeleter
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}