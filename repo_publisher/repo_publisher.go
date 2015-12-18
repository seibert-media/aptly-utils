package repo_publisher

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type RepoPublisher interface {
	PublishNewRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error
}

type repoPublisher struct {
	buildRequestAndExecute     requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

type publishJson struct {
	ForceOverwrite bool
	Distribution   string
	SourceKind     string
	Sources        []map[string]string
	Architectures  []string
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *repoPublisher {
	p := new(repoPublisher)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *repoPublisher) PublishNewRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error {
	logger.Debugf("publishRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s", apiUrl, repo))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJson{
		ForceOverwrite: true,
		Distribution:   distribution,
		SourceKind:     "local",
		Sources:        []map[string]string{{"Name": repo}},
		Architectures:  architectures,
	})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
