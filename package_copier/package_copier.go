package package_copier
import (
	"github.com/bborbe/log"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/aptly/requestbuilder_executor"
)

type PackageCopier interface {
	CopyPackage(apiUrl string, apiUsername string, apiPassword string, sourceRepo string, targetRepo string, pkg string) error
}

type packageCopier struct {
	buildRequestAndExecute     requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *packageCopier {
	p := new(packageCopier)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}


func (c *packageCopier ) CopyPackage(apiUrl string, apiUsername string, apiPassword string, sourceRepo string, targetRepo string, pkg string) error {
	logger.Debugf("CopyPackage")
	return nil
}
