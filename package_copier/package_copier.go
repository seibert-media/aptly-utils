package package_copier

type PackageCopier interface {
	CopyPackage() error
}

type packageCopier struct {

}

func New() *packageCopier {
	return new(packageCopier)
}

func (c *packageCopier ) CopyPackage() error {
	return nil
}
