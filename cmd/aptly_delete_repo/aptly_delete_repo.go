package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_repo_deleter "github.com/seibert-media/aptly-utils/repo_deleter"
	aptly_repo_publisher "github.com/seibert-media/aptly-utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/seibert-media/aptly-utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/io/util"
	"github.com/golang/glog"
)

const (
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepo            = "repo"
	parameterDistribution    = "distribution"
	parameterRepoURL         = "repo-url"
)

var (
	apiURLPtr          = flag.String(parameterAPIURL, "", "url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	repoPtr            = flag.String(parameterRepo, "", "repo")
	distributionPtr    = flag.String(parameterDistribution, string(aptly_model.DistribuionDefault), "distribution")
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	requestbuilderExecutor := aptly_requestbuilder_executor.New(httpClient.Do)
	repoPublisher := aptly_repo_publisher.New(requestbuilderExecutor, http_requestbuilder.NewHTTPRequestBuilderProvider())
	repoDeleter := aptly_repo_deleter.New(requestbuilderExecutor, http_requestbuilder.NewHTTPRequestBuilderProvider(), repoPublisher.UnPublishRepo)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	err := do(
		repoDeleter,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
		*distributionPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	repo_deleter aptly_repo_deleter.RepoDeleter,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	repo string,
	distribution string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v apiPasswordfile: %v repo: %v distribution: %v", repoURL, apiURL, apiUsername, apiPassword, apiPasswordfile, repo, distribution)
	if len(apiPasswordfile) > 0 {
		apiPasswordfile, err := util.NormalizePath(apiPasswordfile)
		if err != nil {
			return fmt.Errorf("normalize path %s failed: %v", apiPasswordfile, err)
		}
		content, err := ioutil.ReadFile(apiPasswordfile)
		if err != nil {
			return fmt.Errorf("read password from file failed: %v", err)
		}
		apiPassword = strings.TrimSpace(string(content))
	}
	if len(apiURL) == 0 {
		return fmt.Errorf("parameter %s missing", parameterAPIURL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", parameterRepo)
	}
	return repo_deleter.DeleteRepo(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution))
}
