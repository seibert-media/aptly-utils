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
	DeletePackage(url string, user string, password string, repo string, name string, version string) error
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

func (p *packageDeleter) DeletePackage(url string, user string, password string, repo string, name string, version string) error {
	key, err := p.findKey(url, user, password, repo, name, version)
	if err != nil {
		return err
	}
	return p.deleteByKey(url, user, password, repo, key)
}

func (p *packageDeleter) findKey(url string, user string, password string, repo string, name string, version string) (string, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repo, name)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", url, repo))
	requestbuilder.AddBasicAuth(user, password)
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return "", err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("request failed: %s", (content))
	}

	var jsonStruct JsonStruct
	err = json.Unmarshal(content, &jsonStruct)
	if err != nil {
		return "", err
	}
	for _, info := range jsonStruct {
		if info["Package"] == name && info["Version"] == version {
			key := info["Key"]
			logger.Debugf("found key: %s", key)
			return key, nil
		}
	}
	return "", fmt.Errorf("package with version not found")
}

func (p *packageDeleter) deleteByKey(url string, user string, password string, repo string, key string) error {
	logger.Debugf("delete package with key: %s", key)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", url, repo))
	requestbuilder.AddBasicAuth(user, password)
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestContent, err := json.Marshal(map[string][]string{"PackageRefs": []string{key}})
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
