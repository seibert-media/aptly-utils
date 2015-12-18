package package_deleter

import (
	"fmt"

	"github.com/bborbe/log"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"bytes"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)
type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PackageDeleter interface {
	DeletePackageByNameAndVersion(url string, user string, password string, repo string, name string, version string) error
	DeletePackagesByKey(url string, user string, password string, repo string, key []string) error
}

type packageDeleter struct {
	executeRequest        ExecuteRequest
	newHttpRequestBuilder NewHttpRequestBuilder
}

var logger = log.DefaultLogger

func New(executeRequest ExecuteRequest, newHttpRequestBuilder NewHttpRequestBuilder) *packageDeleter {
	p := new(packageDeleter)
	p.executeRequest = executeRequest
	p.newHttpRequestBuilder = newHttpRequestBuilder
	return p
}

type JsonStruct []map[string]string

func (p *packageDeleter) DeletePackageByNameAndVersion(url string, user string, password string, repo string, name string, version string) error {
	keys, err := p.findKeys(url, user, password, repo, name, version)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("no package found")
	}
	return p.DeletePackagesByKey(url, user, password, repo, keys)
}

func (p *packageDeleter) findKeys(url string, user string, password string, repo string, name string, version string) ([]string, error) {
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
	var keys []string
	for _, info := range jsonStruct {
		if info["Package"] == name && info["Version"] == version {
			key := info["Key"]
			logger.Debugf("found key: %s", key)
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (p *packageDeleter) DeletePackagesByKey(url string, user string, password string, repo string, keys []string) error {
	logger.Debugf("delete package with keys: %v", keys)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", url, repo))
	requestbuilder.AddBasicAuth(user, password)
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestContent, err := json.Marshal(map[string][]string{"PackageRefs": keys})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(requestContent))
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request failed: %s", (content))
	}
	return nil
}
