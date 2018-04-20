package version

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bborbe/stringutil"
)

type VersionByName []Version

func (v VersionByName) Len() int {
	return len(v)
}
func (v VersionByName) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v VersionByName) Less(i, j int) bool {
	return LessThan(v[i], v[j])
}

func LessThan(a Version, b Version) bool {
	result, err := lessArray(a, b)
	if err != nil {
		return stringutil.StringLess(string(a), string(b))
	}
	return result
}

func GreaterThan(a Version, b Version) bool {
	return LessThan(b, a)
}

func LessEqual(a Version, b Version) bool {
	return !GreaterThan(a, b)
}

func GreaterEqual(a Version, b Version) bool {
	return !LessThan(a, b)
}

func Equal(a Version, b Version) bool {
	return !LessThan(a, b) && !GreaterThan(a, b)
}

func lessArray(a Version, b Version) (bool, error) {
	ai, err := toIntArray(a)
	if err != nil {
		return false, err
	}
	bi, err := toIntArray(b)
	if err != nil {
		return false, err
	}
	if len(ai) != len(bi) {
		return false, fmt.Errorf("can not compare by array")
	}
	for i := 0; i < len(ai); i++ {
		if ai[i] < bi[i] {
			return true, nil
		}
		if ai[i] > bi[i] {
			return false, nil
		}
	}
	return false, nil
}

func toIntArray(v Version) ([]int, error) {
	parts := regSplit(string(v), "[^0-9]+")
	return stringArrayToIntArray(parts)
}

func stringArrayToIntArray(strings []string) ([]int, error) {
	result := make([]int, len(strings))
	for i := 0; i < len(strings); i++ {
		num, err := strconv.Atoi(strings[i])
		if err != nil {
			return nil, err
		}
		result[i] = num
	}
	return result, nil
}

func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}
