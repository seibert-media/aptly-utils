package repo_creator

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PublishNewRepo func(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error

type RepoCreater interface {
	CreateRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error
}

type repoCreater struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	publishNewRepo             PublishNewRepo
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, publishNewRepo PublishNewRepo) *repoCreater {
	p := new(repoCreater)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.publishNewRepo = publishNewRepo
	return p
}

func (c *repoCreater) CreateRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error {
	if err := c.createRepo(apiUrl, apiUsername, apiPassword, repo); err != nil {
		//return err
	}
	if err := c.publishNewRepo(apiUrl, apiUsername, apiPassword, repo, distribution, architectures); err != nil {
		return err
	}
	return nil
}

func (c *repoCreater) createRepo(apiUrl string, apiUsername string, apiPassword string, repo string) error {
	logger.Debugf("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos", apiUrl))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]string{"Name": repo})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
