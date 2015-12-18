package package_versions

import (
	"fmt"

	"github.com/bborbe/log"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bborbe/aptly_utils/version"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)
type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PackageVersions interface {
	PackageVersions(url string, user string, password string, repo string, name string) ([]version.Version, error)
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

type JsonStruct []map[string]string

func (p *packageVersion) PackageVersions(url string, user string, password string, repo string, name string) ([]version.Version, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repo, name)
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

	var jsonStruct JsonStruct
	err = json.Unmarshal(content, &jsonStruct)
	if err != nil {
		return nil, err
	}
	var versions []version.Version
	for _, info := range jsonStruct {
		if info["Package"] == name {
			v := info["Version"]
			logger.Debugf("found version: %s", v)
			versions = append(versions, version.Version(v))
		}
	}
	return versions, nil
}
