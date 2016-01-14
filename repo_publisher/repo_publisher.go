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
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		architectures []aptly_architecture.Architecture) error
	PublishRepo(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
		distribution aptly_distribution.Distribution) error
	UnPublishRepo(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
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
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	architectures []aptly_architecture.Architecture) error {
	logger.Debugf("publishRepo - repo: %s arch: %s", repository, aptly_architecture.Join(architectures, ","))
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s", url, repository))
	requestbuilder.AddBasicAuth(string(user), string(password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJson{
		ForceOverwrite: true,
		Distribution:   distribution,
		SourceKind:     "local",
		Sources:        []map[string]aptly_repository.Repository{{"Name": repository}},
		Architectures:  architectures,
	})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *repoPublisher) PublishRepo(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error {
	logger.Debugf("publishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", url, repository, distribution))
	requestbuilder.AddBasicAuth(string(user), string(password))
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
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error {
	logger.Debugf("unPublishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", url, repository, distribution))
	requestbuilder.AddBasicAuth(string(user), string(password))
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
