package package_lister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type NewHTTPRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PackageLister interface {
	ListPackages(api aptly_model.API, repository aptly_model.Repository) ([]map[string]string, error)
}

type packageVersion struct {
	executeRequest        ExecuteRequest
	newHTTPRequestBuilder NewHTTPRequestBuilder
}

func New(executeRequest ExecuteRequest, newHTTPRequestBuilder NewHTTPRequestBuilder) *packageVersion {
	p := new(packageVersion)
	p.newHTTPRequestBuilder = newHTTPRequestBuilder
	p.executeRequest = executeRequest
	return p
}

func (p *packageVersion) ListPackages(api aptly_model.API, repository aptly_model.Repository) ([]map[string]string, error) {
	glog.V(2).Infof("ListPackages - repo: %s", repository)
	url := fmt.Sprintf("%s/api/repos/%s/packages?format=details", api.APIUrl, repository)
	requestbuilder := p.newHTTPRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.Build()
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
		return nil, fmt.Errorf("request to %s failed with status %d", url, resp.StatusCode)
	}
	var packages []map[string]string
	err = json.Unmarshal(content, &packages)
	if err != nil {
		return nil, err
	}
	return packages, nil
}
