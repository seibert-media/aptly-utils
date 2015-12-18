package package_deleter

type PackageVersion interface {
	PackageVersion(url string, user string, password string, repo string, name string) (string, error)
}

type packageVersion struct {
}

func New() *packageVersion {
	return new(packageVersion)
}

func (c *packageVersion) PackageVersion(url string, user string, password string, repo string, name string) (string, error) {
	return "1.2.3", nil
}
