package repo_deleter


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoDeleter(t *testing.T) {
	b := New()
	var i *RepoDeleter
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}