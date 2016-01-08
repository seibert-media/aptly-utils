package repo_deleter

import (
	"fmt"

	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type UnPublishRepo func(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string) error

type RepoDeleter interface {
	DeleteRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string) error
}

type repoDeleter struct {
	unPublishRepo              UnPublishRepo
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, unPublishRepo UnPublishRepo) *repoDeleter {
	p := new(repoDeleter)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.unPublishRepo = unPublishRepo
	return p
}

func (c *repoDeleter) DeleteRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string) error {
	logger.Debugf("DeleteRepo - repo: %s distribution: %s", repo, distribution)
	err := c.unPublishRepo(apiUrl, apiUsername, apiPassword, repo, distribution)
	if err != nil {
		return err
	}
	return c.deleteRepo(apiUrl, apiUsername, apiPassword, repo)
}

func (p *repoDeleter) deleteRepo(apiUrl string, apiUsername string, apiPassword string, repo string) error {
	logger.Debugf("deleteRepo - repo: %s", repo)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s", apiUrl, repo))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
