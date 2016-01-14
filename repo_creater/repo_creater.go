package repo_creator

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

type PublishNewRepo func(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	architectures []aptly_architecture.Architecture) error

type RepoCreater interface {
	CreateRepo(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repo aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		architectures []aptly_architecture.Architecture) error
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

func (c *repoCreater) CreateRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	architectures []aptly_architecture.Architecture) error {
	if err := c.createRepo(apiUrl, apiUsername, apiPassword, repo); err != nil {
		//return err
	}
	if err := c.publishNewRepo(apiUrl, apiUsername, apiPassword, repo, distribution, architectures); err != nil {
		return err
	}
	return nil
}

func (c *repoCreater) createRepo(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repo aptly_repository.Repository) error {
	logger.Debugf("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos", apiUrl))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]aptly_repository.Repository{"Name": repo})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
