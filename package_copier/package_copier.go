package package_copier

import (
	"fmt"
	"net/http"

	"github.com/seibert-media/aptly-utils/model"
	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_package_uploader "github.com/seibert-media/aptly-utils/package_uploader"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type PackageCopier interface {
	CopyPackage(api aptly_model.API, sourceRepo aptly_model.Repository, targetRepo aptly_model.Repository, targetDistribution aptly_model.Distribution, packageName model.Package, version aptly_version.Version) error
}

type packageCopier struct {
	uploader                   aptly_package_uploader.PackageUploader
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider
	executeRequest             ExecuteRequest
}

func New(uploader aptly_package_uploader.PackageUploader, httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider, executeRequest ExecuteRequest) *packageCopier {
	p := new(packageCopier)
	p.executeRequest = executeRequest
	p.uploader = uploader
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *packageCopier) CopyPackage(
	api aptly_model.API,
	sourceRepo aptly_model.Repository,
	targetRepo aptly_model.Repository,
	targetDistribution aptly_model.Distribution,
	packageName model.Package,
	version aptly_version.Version,
) error {
	glog.V(2).Infof("CopyPackage - sourceRepo: %s targetRepo: %s, targetDistribution: %s, package: %s_%s", sourceRepo, targetRepo, targetDistribution, packageName, version)
	url := fmt.Sprintf("%s/%s/pool/main/%s/%s/%s_%s.deb", api.RepoURL, sourceRepo, packageName[0:1], packageName, packageName, version)
	glog.V(2).Infof("download package url: %s", url)
	requestbuilder := c.httpRequestBuilderProvider.NewHTTPRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := c.executeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	glog.V(2).Infof("download package returncode: %d", resp.StatusCode)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("download package %s %s failed with status %d from url %s", packageName, version, resp.StatusCode, url)
	}
	glog.V(2).Info("download completed, start upload")
	return c.uploader.UploadPackageByReader(api, targetRepo, targetDistribution, fmt.Sprintf("%s_%s.deb", packageName, version), resp.Body)
}
