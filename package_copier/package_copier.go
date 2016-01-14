package package_copier

import (
	"fmt"
	"net/http"

	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PackageCopier interface {
	CopyPackage(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		sourceRepo aptly_repository.Repository,
		targetRepo aptly_repository.Repository,
		targetDistribution aptly_distribution.Distribution,
		packageName package_name.PackageName,
		version aptly_version.Version) error
}

type packageCopier struct {
	uploader                   aptly_package_uploader.PackageUploader
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	client                     *http.Client
}

var logger = log.DefaultLogger

func New(uploader aptly_package_uploader.PackageUploader, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, client *http.Client) *packageCopier {
	p := new(packageCopier)
	p.client = client
	p.uploader = uploader
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *packageCopier) CopyPackage(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	sourceRepo aptly_repository.Repository,
	targetRepo aptly_repository.Repository,
	targetDistribution aptly_distribution.Distribution,
	packageName package_name.PackageName,
	version aptly_version.Version) error {
	logger.Debugf("CopyPackage - sourceRepo: %s targetRepo: %s, package: %s_%s", sourceRepo, targetRepo, packageName, version)
	requestUrl := fmt.Sprintf("%s/%s/pool/main/%s/%s/%s_%s.deb", url, sourceRepo, packageName[0:1], packageName, packageName, version)
	logger.Debugf("download package url: %s", requestUrl)
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(requestUrl)
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("download package %s %s failed", packageName, version)
	}
	return c.uploader.UploadPackageByReader(url, user, password, targetRepo, targetDistribution, packageName, resp.Body)
}
