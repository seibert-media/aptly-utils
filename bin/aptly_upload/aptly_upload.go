package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
	"github.com/bborbe/io/util"
)

const (
	parameterFile            = "file"
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepo            = "repo"
	parameterRepoURL         = "repo-url"
	parameterDistribution    = "distribution"
)

var (
	filePtr            = flag.String(parameterFile, "", "file")
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
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	package_uploader := aptly_package_uploader.New(requestbuilder_executor, httpRequestBuilderProvider, repo_publisher.PublishRepo)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	err := do(
		package_uploader,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*filePtr,
		*repoPtr,
		*distributionPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	package_uploader aptly_package_uploader.PackageUploader,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	file string,
	repo string,
	distribution string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v apiPasswordfile: %v file: %v repo: %v distribution: %v", repoURL, apiURL, apiUsername, apiPassword, apiPasswordfile, file, repo, distribution)
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
	if len(file) == 0 {
		return fmt.Errorf("parameter %s missing", parameterFile)
	}
	glog.V(2).Infof("upload file %s to repo %s dist %s on server %s", file, repo, distribution, apiURL)
	return package_uploader.UploadPackageByFile(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution), file)
}
