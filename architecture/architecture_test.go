package architecture

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestConsts(t *testing.T) {
	if err := AssertThat(string(ALL), Is("all")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(AMD64), Is("amd64")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(I386), Is("i386")); err != nil {
		t.Fatal(err)
	}
}

func TestDefault(t *testing.T) {
	if err := AssertThat(string(DEFAULT), Is("amd64")); err != nil {
		t.Fatal(err)
	}
}
