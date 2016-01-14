package package_versions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)
type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PackageLister interface {
	ListPackages(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repository aptly_repository.Repository) ([]map[string]string, error)
}

type packageVersion struct {
	executeRequest        ExecuteRequest
	newHttpRequestBuilder NewHttpRequestBuilder
}

var logger = log.DefaultLogger

func New(executeRequest ExecuteRequest, newHttpRequestBuilder NewHttpRequestBuilder) *packageVersion {
	p := new(packageVersion)
	p.newHttpRequestBuilder = newHttpRequestBuilder
	p.executeRequest = executeRequest
	return p
}

func (p *packageVersion) ListPackages(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository) ([]map[string]string, error) {
	logger.Debugf("PackageVersions - repo: %s", repository)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", apiUrl, repository))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return nil, err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed: %s", (content))
	}

	var packages []map[string]string
	err = json.Unmarshal(content, &packages)
	if err != nil {
		return nil, err
	}
	return packages, nil
}
