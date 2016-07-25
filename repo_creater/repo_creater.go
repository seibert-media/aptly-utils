package repo_creator

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PublishNewRepo func(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error

type RepoCreater interface {
	CreateRepo(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error
}

type repoCreater struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	publishNewRepo             PublishNewRepo
}

var logger = log.DefaultLogger

func New(
	buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor,
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider,
	publishNewRepo PublishNewRepo,
) *repoCreater {
	p := new(repoCreater)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.publishNewRepo = publishNewRepo
	return p
}

func (c *repoCreater) CreateRepo(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error {
	if err := c.createRepo(api, repository); err != nil {
		//return err
	}
	if err := c.publishNewRepo(api, repository, distribution, architectures); err != nil {
		return err
	}
	return nil
}

func (c *repoCreater) createRepo(api aptly_model.Api, repository aptly_model.Repository) error {
	logger.Debugf("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos", api.ApiUrl))
	requestbuilder.AddBasicAuth(string(api.ApiUsername), string(api.ApiPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]aptly_model.Repository{"Name": repository})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
