package package_versions

import (
	"fmt"

	"github.com/bborbe/log"

	"encoding/json"
	"io/ioutil"
	"net/http"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)
type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PackageLister interface {
	ListPackages(url string, user string, password string, repo string) ([]map[string]string, error)
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

func (p *packageVersion) ListPackages(url string, user string, password string, repo string) ([]map[string]string, error) {
	logger.Debugf("PackageVersions - repo: %s", repo)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", url, repo))
	requestbuilder.AddBasicAuth(user, password)
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
