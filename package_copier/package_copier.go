package package_copier

import (
	"fmt"
	"net/http"

	"github.com/bborbe/aptly_utils/model"
	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type PackageCopier interface {
	CopyPackage(api aptly_model.Api, sourceRepo aptly_model.Repository, targetRepo aptly_model.Repository, targetDistribution aptly_model.Distribution, packageName model.Package, version aptly_version.Version) error
}

type packageCopier struct {
	uploader                   aptly_package_uploader.PackageUploader
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
}

var logger = log.DefaultLogger

func New(uploader aptly_package_uploader.PackageUploader, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, executeRequest ExecuteRequest) *packageCopier {
	p := new(packageCopier)
	p.executeRequest = executeRequest
	p.uploader = uploader
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *packageCopier) CopyPackage(
	api aptly_model.Api,
	sourceRepo aptly_model.Repository,
	targetRepo aptly_model.Repository,
	targetDistribution aptly_model.Distribution,
	packageName model.Package,
	version aptly_version.Version,
) error {
	logger.Debugf("CopyPackage - sourceRepo: %s targetRepo: %s, targetDistribution: %s, package: %s_%s", sourceRepo, targetRepo, targetDistribution, packageName, version)
	url := fmt.Sprintf("%s/%s/pool/main/%s/%s/%s_%s.deb", api.RepoUrl, sourceRepo, packageName[0:1], packageName, packageName, version)
	logger.Debugf("download package url: %s", url)
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.ApiUsername), string(api.ApiPassword))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := c.executeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logger.Debugf("download package returncode: %d", resp.StatusCode)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("download package %s %s failed with status %d from url %s", packageName, version, resp.StatusCode, url)
	}
	logger.Debug("download completed, start upload")
	return c.uploader.UploadPackageByReader(api, targetRepo, targetDistribution, fmt.Sprintf("%s_%s.deb", packageName, version), resp.Body)
}
