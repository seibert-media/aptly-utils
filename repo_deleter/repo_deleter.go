package repo_deleter

import (
	"fmt"

	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type UnPublishRepo func(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error

type RepoDeleter interface {
	DeleteRepo(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repo aptly_repository.Repository,
		distribution aptly_distribution.Distribution) error
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

func (c *repoDeleter) DeleteRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error {
	logger.Debugf("DeleteRepo - repo: %s distribution: %s", repo, distribution)
	err := c.unPublishRepo(apiUrl, apiUsername, apiPassword, repo, distribution)
	if err != nil {
		return err
	}
	return c.deleteRepo(apiUrl, apiUsername, apiPassword, repo)
}

func (p *repoDeleter) deleteRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository) error {
	logger.Debugf("deleteRepo - repo: %s", repo)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s", apiUrl, repo))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
