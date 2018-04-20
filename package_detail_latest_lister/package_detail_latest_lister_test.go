package package_detail_latest_lister

import (
	"testing"

	aptly_model "github.com/seibert-media/aptly-utils/model"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageDetailLatestLister(t *testing.T) {
	b := New(nil)
	var i *PackageDetailLatestLister
	if err := AssertThat(b, Implements(i).Message("check type")); err != nil {
		t.Fatal(err)
	}
}

func TestLatestOne(t *testing.T) {
	result := latest(aptly_model.NewPackageDetailByString("abc", "1.2.3"))
	if err := AssertThat(len(result), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Package), Is("abc")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Version), Is("1.2.3")); err != nil {
		t.Fatal(err)
	}
}

func TestLatestTwoGreater(t *testing.T) {
	result := latest(aptly_model.NewPackageDetailByString("abc", "1.2.3"), aptly_model.NewPackageDetailByString("abc", "1.2.2"))
	if err := AssertThat(len(result), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Package), Is("abc")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Version), Is("1.2.3")); err != nil {
		t.Fatal(err)
	}
}

func TestLatestTwoLess(t *testing.T) {
	result := latest(aptly_model.NewPackageDetailByString("abc", "1.2.2"), aptly_model.NewPackageDetailByString("abc", "1.2.3"))
	if err := AssertThat(len(result), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Package), Is("abc")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(result[0].Version), Is("1.2.3")); err != nil {
		t.Fatal(err)
	}
}
