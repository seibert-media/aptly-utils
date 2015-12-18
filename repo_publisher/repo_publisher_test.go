package repo_publisher

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRepoPublisher(t *testing.T) {
	b := New(nil, nil)
	var i *RepoPublisher
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
