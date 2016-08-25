package repo_publisher

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type RepoPublisher interface {
	PublishNewRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error
	PublishRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error
	UnPublishRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error
}

type repoPublisher struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider
}

type publishJSON struct {
	ForceOverwrite bool
	Distribution   aptly_model.Distribution
	SourceKind     string
	Sources        []map[string]aptly_model.Repository
	Architectures  []aptly_model.Architecture
}

var logger = log.DefaultLogger

func New(
	buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor,
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider,
) *repoPublisher {
	r := new(repoPublisher)
	r.buildRequestAndExecute = buildRequestAndExecute
	r.httpRequestBuilderProvider = httpRequestBuilderProvider
	return r
}

func (r *repoPublisher) PublishNewRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error {
	logger.Debugf("PublishNewRepo - repo: %s dist: %s arch: %s", repository, distribution, aptly_model.JoinArchitectures(architectures, ","))
	requestbuilder := r.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/publish/%s", api.APIUrl, repository))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJSON{
		ForceOverwrite: true, Distribution: distribution, SourceKind: "local", Sources: []map[string]aptly_model.Repository{{"Name": repository}}, Architectures: architectures})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	logger.Debugf("PublishNewRepo ...")
	if err := r.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder); err != nil {
		return err
	}
	logger.Debugf("PublishNewRepo finished")
	return nil
}

func (r *repoPublisher) PublishRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error {
	logger.Debugf("PublishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := r.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", api.APIUrl, repository, distribution))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("PUT")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]bool{"ForceOverwrite": true})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	logger.Debugf("PublishRepo ...")
	if err := r.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder); err != nil {
		return err
	}
	logger.Debugf("PublishRepo finished")
	return nil
}

func (r *repoPublisher) UnPublishRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error {
	logger.Debugf("UnPublishRepo - repo: %s distribution: %s", repository, distribution)
	requestbuilder := r.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", api.APIUrl, repository, distribution))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("DELETE")
	logger.Debugf("UnPublishRepo ...")
	if err := r.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder); err != nil {
		return err
	}
	logger.Debugf("UnPublishRepo finished")
	return nil
}
