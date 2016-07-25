package model

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestDefaultDistribution(t *testing.T) {
	if err := AssertThat(string(DISTRIBUTION_DEFAULT), Is("default")); err != nil {
		t.Fatal(err)
	}
}

func TestArchitectureConsts(t *testing.T) {
	if err := AssertThat(string(ARCHITECTURE_ALL), Is("all")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(ARCHITECTURE_AMD64), Is("amd64")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(ARCHITECTURE_I386), Is("i386")); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultArchitecture(t *testing.T) {
	if err := AssertThat(string(ARCHITECTURE_DEFAULT), Is("amd64")); err != nil {
		t.Fatal(err)
	}
}
