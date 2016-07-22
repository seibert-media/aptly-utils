package repo_publisher

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

type RepoPublisher interface {
	PublishNewRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, architectures []aptly_architecture.Architecture) error
	PublishRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error
	UnPublishRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error
}

type repoPublisher struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

type publishJson struct {
	ForceOverwrite bool
	Distribution   aptly_distribution.Distribution
	SourceKind     string
	Sources        []map[string]aptly_repository.Repository
	Architectures  []aptly_architecture.Architecture
}

var logger = log.DefaultLogger

func New(
	buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor,
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider,
) *repoPublisher {
	p := new(repoPublisher)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *repoPublisher) PublishNewRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, architectures []aptly_architecture.Architecture) error {
	logger.Debugf("publishRepo - repo: %s dist: %s arch: %s", repository, distribution, aptly_architecture.Join(architectures, ","))
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s", api.Url, repository))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJson{
		ForceOverwrite: true, Distribution: distribution, SourceKind: "local", Sources: []map[string]aptly_repository.Repository{{"Name": repository}}, Architectures: architectures})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *repoPublisher) PublishRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error {
	logger.Debugf("publishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", api.Url, repository, distribution))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("PUT")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]bool{"ForceOverwrite": true})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *repoPublisher) UnPublishRepo(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error {
	logger.Debugf("unPublishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", api.Url, repository, distribution))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
