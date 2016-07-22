package repo_details

import (
	"encoding/json"
	"fmt"
	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
	"io/ioutil"
	"net/http"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type RepoDetails interface {
	Repos(api aptly_api.Api, repository aptly_repository.Repository) (map[string]string, error)
}

type repoDetails struct {
	executeRequest        ExecuteRequest
	newHttpRequestBuilder NewHttpRequestBuilder
}

var logger = log.DefaultLogger

func New(executeRequest ExecuteRequest, newHttpRequestBuilder NewHttpRequestBuilder) *repoDetails {
	p := new(repoDetails)
	p.newHttpRequestBuilder = newHttpRequestBuilder
	p.executeRequest = executeRequest
	return p
}

func (p *repoDetails) Repos(api aptly_api.Api, repository aptly_repository.Repository) (map[string]string, error) {
	logger.Debugf("list repos")
	url := fmt.Sprintf("%s/api/repos/%s", api.Url, repository)
	requestbuilder := p.newHttpRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
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
	var repo map[string]string
	err = json.Unmarshal(content, &repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
