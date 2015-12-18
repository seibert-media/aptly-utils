package version

import (
	"testing"

	"sort"

	. "github.com/bborbe/assert"
)

func TestSortByNameNumber(t *testing.T) {
	versions := []Version{"1", "3", "2"}
	sort.Sort(VersionByName(versions))

	if err := AssertThat(versions[0], Is(Version("1"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[1], Is(Version("2"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[2], Is(Version("3"))); err != nil {
		t.Fatal(err)
	}
}

func TestSortByNameLetter(t *testing.T) {
	versions := []Version{"a", "c", "b"}
	sort.Sort(VersionByName(versions))

	if err := AssertThat(versions[0], Is(Version("a"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[1], Is(Version("b"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[2], Is(Version("c"))); err != nil {
		t.Fatal(err)
	}
}

func TestSortByNameNumberLength(t *testing.T) {
	versions := []Version{"3", "11", "2"}
	sort.Sort(VersionByName(versions))

	if err := AssertThat(versions[0], Is(Version("2"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[1], Is(Version("3"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[2], Is(Version("11"))); err != nil {
		t.Fatal(err)
	}
}

func TestSortByNameNumberWithDot(t *testing.T) {
	versions := []Version{"1.0.1", "1.0.3", "1.0.2"}
	sort.Sort(VersionByName(versions))

	if err := AssertThat(versions[0], Is(Version("1.0.1"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[1], Is(Version("1.0.2"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[2], Is(Version("1.0.3"))); err != nil {
		t.Fatal(err)
	}
}

func TestSortByNameComplex(t *testing.T) {
	versions := []Version{"1.0.1-b80", "1.0.1-b102", "1.0.1-b99"}
	sort.Sort(VersionByName(versions))

	if err := AssertThat(versions[0], Is(Version("1.0.1-b80"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[1], Is(Version("1.0.1-b99"))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(versions[2], Is(Version("1.0.1-b102"))); err != nil {
		t.Fatal(err)
	}
}
