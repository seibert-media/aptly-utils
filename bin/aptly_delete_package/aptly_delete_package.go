package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/bborbe/aptly_utils/model"
	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_deleter "github.com/bborbe/aptly_utils/package_deleter"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

const (
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO_URL          = "repo-url"
	PARAMETER_REPO              = "repo"
	PARAMETER_NAME              = "name"
	PARAMETER_VERSION           = "version"
	PARAMETER_DISTRIBUTION      = "distribution"
)

var (
	logger             = log.DefaultLogger
	logLevelPtr        = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr          = flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr         = flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr            = flag.String(PARAMETER_REPO, "", "repo")
	namePtr            = flag.String(PARAMETER_NAME, "", "name")
	versionPtr         = flag.String(PARAMETER_VERSION, "", "version")
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
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	package_deleter := aptly_package_deleter.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder, repo_publisher.PublishRepo)

	if len(*repoUrlPtr) == 0 {
		*repoUrlPtr = *apiUrlPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		package_deleter,
		*repoUrlPtr,
		*apiUrlPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
		*distributionPtr,
		*namePtr,
		*versionPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(
	writer io.Writer,
	package_deleter aptly_package_deleter.PackageDeleter,
	repoUrl string,
	apiUrl string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	repo string,
	distribution string,
	name string,
	version string,
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
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_NAME)
	}
	if len(version) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_VERSION)
	}
	return package_deleter.DeletePackageByNameAndVersion(aptly_model.NewApi(repoUrl, apiUrl, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Distribution(distribution), model.Package(name), aptly_version.Version(version))
}
