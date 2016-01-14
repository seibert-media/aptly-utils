package package_name

import "strings"

type PackageName string

func FromFileName(path string) PackageName {
	slashPos := strings.LastIndex(path, "/")
	if slashPos != -1 {
		return PackageName(path[slashPos+1:])
	}
	return PackageName(path)
}
