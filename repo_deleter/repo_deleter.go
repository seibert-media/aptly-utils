package repo_deleter

import (
	"fmt"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type UnPublishRepo func(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error

type RepoDeleter interface {
	DeleteRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error
}

type repoDeleter struct {
	unPublishRepo              UnPublishRepo
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider
}

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider, unPublishRepo UnPublishRepo) *repoDeleter {
	r := new(repoDeleter)
	r.buildRequestAndExecute = buildRequestAndExecute
	r.httpRequestBuilderProvider = httpRequestBuilderProvider
	r.unPublishRepo = unPublishRepo
	return r
}

func (r *repoDeleter) DeleteRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error {
	logger.Debugf("DeleteRepo - repo: %s distribution: %s", repository, distribution)
	err := r.unPublishRepo(api, repository, distribution)
	if err != nil {
		//return err
	}
	return r.deleteRepo(api, repository)
}

func (r *repoDeleter) deleteRepo(api aptly_model.API, repository aptly_model.Repository) error {
	logger.Debugf("deleteRepo - repo: %s", repository)
	requestbuilder := r.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/repos/%s", api.APIUrl, repository))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("DELETE")
	return r.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
