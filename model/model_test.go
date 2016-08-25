package model

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestDefaultDistribution(t *testing.T) {
	if err := AssertThat(string(DistribuionDefault), Is("default")); err != nil {
		t.Fatal(err)
	}
}

func TestArchitectureConsts(t *testing.T) {
	if err := AssertThat(string(ArchitectureALL), Is("all")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(ArchitectureAMD64), Is("amd64")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(ArchitectureI386), Is("i386")); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultArchitecture(t *testing.T) {
	if err := AssertThat(string(ArchitectureDefault), Is("amd64")); err != nil {
		t.Fatal(err)
	}
}
