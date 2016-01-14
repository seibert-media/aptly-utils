package package_name

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestFromFileNameWithoutSlash(t *testing.T) {
	name := FromFileName("foo.deb")
	if err := AssertThat(string(name), Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}

func TestFromFileNameWithSlash(t *testing.T) {
	name := FromFileName("asdf/foo.deb")
	if err := AssertThat(string(name), Is("foo.deb")); err != nil {
		t.Fatal(err)
	}
}
