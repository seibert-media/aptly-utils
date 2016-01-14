package distribution

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestDefault(t *testing.T) {
	if err := AssertThat(string(DEFAULT), Is("default")); err != nil {
		t.Fatal(err)
	}
}
