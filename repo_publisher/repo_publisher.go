package repo_publisher

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_architecture "github.com/bborbe/aptly_utils/architecture"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type RepoPublisher interface {
	PublishNewRepo(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repo aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		architectures []aptly_architecture.Architecture) error
	PublishRepo(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repo aptly_repository.Repository,
		distribution aptly_distribution.Distribution) error
	UnPublishRepo(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repo aptly_repository.Repository,
		distribution aptly_distribution.Distribution) error
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

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *repoPublisher {
	p := new(repoPublisher)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *repoPublisher) PublishNewRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	architectures []aptly_architecture.Architecture) error {
	logger.Debugf("publishRepo - repo: %s arch: %s", repo, aptly_architecture.Join(architectures, ","))
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s", apiUrl, repo))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJson{
		ForceOverwrite: true,
		Distribution:   distribution,
		SourceKind:     "local",
		Sources:        []map[string]aptly_repository.Repository{{"Name": repo}},
		Architectures:  architectures,
	})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *repoPublisher) PublishRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error {
	logger.Debugf("publishRepo - repo: %s distribution: %s", repo, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", apiUrl, repo, distribution))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("PUT")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]bool{"ForceOverwrite": true})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *repoPublisher) UnPublishRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error {
	logger.Debugf("unPublishRepo - repo: %s distribution: %s", repo, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", apiUrl, repo, distribution))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
