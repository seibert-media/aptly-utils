package package_deleter

type PackageDeleter interface {
	DeletePackage() error
}

type packageDeleter struct {

}

func New() *packageDeleter {
	return new(packageDeleter)
}

func (c *packageDeleter ) DeletePackage() error {
	return nil
}
