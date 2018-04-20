package repo_creator

import (
	"bytes"
	"encoding/json"
	"fmt"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_requestbuilder_executor "github.com/seibert-media/aptly-utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

type PublishNewRepo func(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error

type RepoCreater interface {
	CreateRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error
}

type repoCreater struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider
	publishNewRepo             PublishNewRepo
}

func New(
	buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor,
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider,
	publishNewRepo PublishNewRepo,
) *repoCreater {
	p := new(repoCreater)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.publishNewRepo = publishNewRepo
	return p
}

func (c *repoCreater) CreateRepo(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, architectures []aptly_model.Architecture) error {
	if err := validateArchitectures(architectures...); err != nil {
		return err
	}
	if err := c.createRepo(api, repository); err != nil {
		//return err
	}
	if err := c.publishNewRepo(api, repository, distribution, architectures); err != nil {
		return err
	}
	return nil
}

func validateArchitectures(architectures ...aptly_model.Architecture) error {
	for _, architecture := range architectures {
		allowed := false
		for _, allowedArchitecture := range aptly_model.AllowedArchitectures {
			if architecture == allowedArchitecture {
				allowed = true
			}
		}
		if !allowed {
			return fmt.Errorf("not support architecture: %v", architecture)
		}
	}
	return nil
}

func (c *repoCreater) createRepo(api aptly_model.API, repository aptly_model.Repository) error {
	glog.V(2).Infof("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/repos", api.APIUrl))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]aptly_model.Repository{"Name": repository})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
