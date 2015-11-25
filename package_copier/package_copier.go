package package_copier

import (
	"fmt"
	"net/http"

	"github.com/bborbe/aptly-utils/package_uploader"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PackageCopier interface {
	CopyPackage(apiUrl string, apiUsername string, apiPassword string, sourceRepo string, targetRepo string, name, version string) error
}

type packageCopier struct {
	uploader                   package_uploader.PackageUploader
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	client                     *http.Client
}

var logger = log.DefaultLogger

func New(uploader package_uploader.PackageUploader, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, client *http.Client) *packageCopier {
	p := new(packageCopier)
	p.client = client
	p.uploader = uploader
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *packageCopier) CopyPackage(apiUrl string, apiUsername string, apiPassword string, sourceRepo string, targetRepo string, name, version string) error {
	logger.Debugf("CopyPackage - sourceRepo: %s targetRepo: %s, package: %s_%s", sourceRepo, targetRepo, name, version)
	url := fmt.Sprintf("%s/%s/pool/main/%s/%s/%s_%s.deb", apiUrl, sourceRepo, name[0:1], name, name, version)
	logger.Debugf("download package url: %s", url)
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(url)
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("download package %s_%s.deb failed", name, version)
	}
	return c.uploader.UploadPackageByReader(apiUrl, apiUsername, apiPassword, targetRepo, fmt.Sprintf("%s_%s.deb", name, version), resp.Body)
}
