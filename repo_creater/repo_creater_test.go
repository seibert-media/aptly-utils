package repo_creator

import (
	"testing"

	"github.com/bborbe/aptly_utils/model"
	. "github.com/bborbe/assert"
)

func TestImplementsRepoDeleter(t *testing.T) {
	b := New(nil, nil, nil)
	var i *RepoCreater
	if err := AssertThat(b, Implements(i).Message("check type")); err != nil {
		t.Fatal(err)
	}
}

func TestValidateArchitectures(t *testing.T) {
	if err := AssertThat(validateArchitectures(model.ArchitectureAMD64), NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(validateArchitectures(model.ArchitectureALL), NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(validateArchitectures(model.ArchitectureI386), NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(validateArchitectures(model.Architecture("asadf")), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
