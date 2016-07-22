package repo_details

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoDetails(t *testing.T) {
	b := New(nil, nil)
	var i *RepoDetails
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
