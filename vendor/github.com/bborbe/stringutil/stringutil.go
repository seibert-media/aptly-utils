package stringutil

import (
	"fmt"
	"strings"
)

func StringAfter(content string, find string) string {
	match, err := StringAfterError(content, find)
	if err != nil {
		return ""
	}
	return match
}

func StringAfterError(content string, find string) (string, error) {
	if len(find) == 0 {
		return content, nil
	}
	pos := strings.Index(content, find)
	if pos == -1 {
		return "", fmt.Errorf("can't find '%s' in content", find)
	}
	return content[pos+len(find):], nil
}

func StringBefore(content string, find string) string {
	match, err := StringBeforeError(content, find)
	if err != nil {
		return ""
	}
	return match
}

func StringBeforeError(content string, find string) (string, error) {
	if len(find) == 0 {
		return content, nil
	}
	pos := strings.Index(content, find)
	if pos == -1 {
		return "", fmt.Errorf("can't find '%s' in content", find)
	}
	return content[:pos], nil
}

func Trim(content string) string {
	runes := []rune(content)
	if len(runes) > 0 {
		if trimableChar(runes[0]) {
			return Trim(content[1:])
		}
		if trimableChar(runes[len(content)-1]) {
			return Trim(content[:len(content)-1])
		}
	}
	return content
}

func trimableChar(c rune) bool {
	return c == ' ' || c == '\n' || c == '\r'
}

func StringLess(a, b string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return len(a) < len(b)
}
