package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"

	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_deleter "github.com/bborbe/aptly_utils/package_deleter"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_repo_cleaner "github.com/bborbe/aptly_utils/repo_cleaner"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
	"github.com/bborbe/io/util"
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
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
	repoPtr            = flag.String(parameterRepo, "", "repo")
	distributionPtr    = flag.String(parameterDistribution, string(aptly_model.DistribuionDefault), "distribution")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	packageLister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	repoPublisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	packageDeleter := aptly_package_deleter.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder, repoPublisher.PublishRepo)
	repoCleaner := aptly_repo_cleaner.New(packageDeleter.DeletePackagesByKey, packageLister.ListPackages)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	err := do(
		repoCleaner,
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
	repo_cleaner aptly_repo_cleaner.RepoCleaner,
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
	err := repo_cleaner.CleanRepo(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution))
	if err != nil {
		return err
	}
	return nil
}
