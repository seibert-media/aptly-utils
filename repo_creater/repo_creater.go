package repo_creator
import (
	"github.com/bborbe/log"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/bborbe/aptly/requestbuilder_executor"
	"github.com/bborbe/aptly/defaults"
)

type RepoDeleter interface {
	DeleteRepo(apiUrl string, apiUsername string, apiPassword string, repo string) error
}

type repoDeleter struct {
	buildRequestAndExecute     requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}


var logger = log.DefaultLogger

func New(buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *repoDeleter {
	p := new(repoDeleter)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (c *repoDeleter ) DeleteRepo(apiUrl string, apiUsername string, apiPassword string, repo string) error {
	if err := c.createRepo(apiUrl, apiUsername, apiPassword, repo); err != nil {
		return err
	}
	if err := c.publishRepo(apiUrl, apiUsername, apiPassword, repo, defaults.DEFAULT_DISTRIBUTION, []string{defaults.DEFAULT_ARCHITECTURE}); err != nil {
		return err
	}
	return nil
}

func (c *repoDeleter ) createRepo(apiUrl string, apiUsername string, apiPassword string, repo string) error {
	logger.Debugf("createRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos", apiUrl))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]string{"Name": repo})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

type publishJson struct {
	ForceOverwrite bool
	Distribution   string
	SourceKind     string
	Sources        []map[string]string
	Architectures  []string
}

func (c *repoDeleter ) publishRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string, architectures []string) error {
	logger.Debugf("publishRepo")
	requestbuilder := c.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s", apiUrl, repo))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(publishJson{
		ForceOverwrite:true,
		Distribution: distribution,
		SourceKind: "local",
		Sources:   []map[string]string{{"Name":repo}},
		Architectures: architectures,
	})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return c.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}