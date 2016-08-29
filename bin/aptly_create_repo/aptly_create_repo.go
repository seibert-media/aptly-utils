package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

const (
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepo            = "repo"
	parameterDistribution    = "distribution"
	parameterArchitecture    = "architecture"
	parameterRepoURL         = "repo-url"
)

var (
	apiURLPtr          = flag.String(parameterAPIURL, "", "url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	repoPtr            = flag.String(parameterRepo, "", "repo")
	distributionPtr    = flag.String(parameterDistribution, string(aptly_model.DistribuionDefault), "distribution")
	architecturePtr    = flag.String(parameterArchitecture, string(aptly_model.ArchitectureDefault), "architecture")
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	requestbuilder := http_requestbuilder.NewHTTPRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	repo_creater := aptly_repo_creater.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishNewRepo)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	err := do(
		repo_creater,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
		*distributionPtr,
		strings.Split(*architecturePtr, ","),
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	repo_creater aptly_repo_creater.RepoCreater,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	repo string,
	distribution string,
	architectures []string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v apiPasswordfile: %v repo: %v distribution: %v architectures: %v", repoURL, apiURL, apiUsername, apiPassword, apiPasswordfile, repo, distribution, architectures)
	if len(apiPasswordfile) > 0 {
		content, err := ioutil.ReadFile(apiPasswordfile)
		if err != nil {
			return err
		}
		apiPassword = strings.TrimSpace(string(content))
	}
	if len(apiURL) == 0 {
		return fmt.Errorf("parameter %s missing", parameterAPIURL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", parameterRepo)
	}
	return repo_creater.CreateRepo(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution), aptly_model.ParseArchitectures(architectures))
}
