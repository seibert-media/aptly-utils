package repo_deleter

import (
	"sort"
	"testing"

	aptly_model "github.com/bborbe/aptly_utils/model"
	. "github.com/bborbe/assert"
)

func TestImplementsRepoCleaner(t *testing.T) {
	b := New(nil, nil)
	var i *RepoCleaner
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPackagesToKeys(t *testing.T) {
	keys := packagesToKeys(nil)
	if err := AssertThat(len(keys), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestPackagesToKeysOnePackage(t *testing.T) {
	keys := packagesToKeys([]map[string]string{map[string]string{"Key": "keyA", "Package": "packageA", "Version": "1.2.3"}})
	if err := AssertThat(len(keys), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestPackagesToKeysTwoDifferentPackages(t *testing.T) {
	keys := packagesToKeys([]map[string]string{map[string]string{"Key": "keyA", "Package": "packageA", "Version": "1.2.3"}, map[string]string{"Key": "keyB", "Package": "packageB", "Version": "1.2.3"}})
	if err := AssertThat(len(keys), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestPackagesToKeysTwoDifferentVersions(t *testing.T) {
	keys := packagesToKeys([]map[string]string{map[string]string{"Key": "keyA", "Package": "packageA", "Version": "1.2.3"}, map[string]string{"Key": "keyB", "Package": "packageA", "Version": "1.2.2"}})
	if err := AssertThat(len(keys), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(keys[0]), Is("keyB")); err != nil {
		t.Fatal(err)
	}
}

func TestPackagesToKeysThreeDifferentVersions(t *testing.T) {
	keys := packagesToKeys([]map[string]string{map[string]string{"Key": "keyA", "Package": "packageA", "Version": "1.2.3"}, map[string]string{"Key": "keyB", "Package": "packageA", "Version": "1.2.2"}, map[string]string{"Key": "keyA", "Package": "packageA", "Version": "1.2.4"}})
	if err := AssertThat(len(keys), Is(2)); err != nil {
		t.Fatal(err)
	}
	sort.Sort(aptly_model.KeySlice(keys))
	if err := AssertThat(string(keys[0]), Is("keyA")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(keys[1]), Is("keyB")); err != nil {
		t.Fatal(err)
	}
}
