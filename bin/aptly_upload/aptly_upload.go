package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

const (
	PARAMETER_FILE              = "file"
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
	PARAMETER_REPO_URL          = "repo-url"
	PARAMETER_DISTRIBUTION      = "distribution"
)

var (
	logger             = log.DefaultLogger
	logLevelPtr        = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	filePtr            = flag.String(PARAMETER_FILE, "", "file")
	apiUrlPtr          = flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr         = flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr            = flag.String(PARAMETER_REPO, "", "repo")
	distributionPtr    = flag.String(PARAMETER_DISTRIBUTION, string(aptly_model.DISTRIBUTION_DEFAULT), "distribution")
	repoUrlPtr         = flag.String(PARAMETER_REPO_URL, "", "repo url")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	package_uploader := aptly_package_uploader.New(requestbuilder_executor, httpRequestBuilderProvider, repo_publisher.PublishRepo)

	if len(*repoUrlPtr) == 0 {
		*repoUrlPtr = *apiUrlPtr
	}

	err := do(
		package_uploader,
		*repoUrlPtr,
		*apiUrlPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*filePtr,
		*repoPtr,
		*distributionPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(
	package_uploader aptly_package_uploader.PackageUploader,
	repoUrl string,
	apiUrl string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	file string,
	repo string,
	distribution string,
) error {
	if len(apiPasswordfile) > 0 {
		content, err := ioutil.ReadFile(apiPasswordfile)
		if err != nil {
			return err
		}
		apiPassword = strings.TrimSpace(string(content))
	}
	if len(apiUrl) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_API_URL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_REPO)
	}
	if len(file) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_FILE)
	}
	logger.Debugf("upload file %s to repo %s dist %s on server %s", file, repo, distribution, apiUrl)
	return package_uploader.UploadPackageByFile(aptly_model.NewApi(repoUrl, apiUrl, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution), file)
}
