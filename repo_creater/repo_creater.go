package repo_creator

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_architecture "github.com/bborbe/aptly_utils/architecture"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PublishNewRepo func(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, architectures []aptly_architecture.Architecture) error

type RepoCreater interface {
	CreateRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, architectures []aptly_architecture.Architecture) error
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

func (c *repoCreater) CreateRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, architectures []aptly_architecture.Architecture) error {
	if err := c.createRepo(api, repository); err != nil {
		//return err
	}
	if err := c.publishNewRepo(api, repository, distribution, architectures); err != nil {
		return err
	}
	return nil
}

func (c *repoCreater) createRepo(api aptly_api.Api, repository aptly_repository.Repository) error {
	logger.Debugf("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos", api.Url))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]aptly_repository.Repository{"Name": repository})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
