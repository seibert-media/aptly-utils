package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"strings"

	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_package_deleter "github.com/bborbe/aptly_utils/package_deleter"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_repo_cleaner "github.com/bborbe/aptly_utils/repo_cleaner"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
	http_client_builder "github.com/bborbe/http/client_builder"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_API_URL = "url"
	PARAMETER_API_USER = "username"
	PARAMETER_API_PASSWORD = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO = "repo"
	PARAMETER_DISTRIBUTION = "distribution"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	urlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	passwordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	passwordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	distributionPtr := flag.String(PARAMETER_DISTRIBUTION, string(aptly_distribution.DEFAULT), "distribution")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	packageLister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	repoPublisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	packageDeleter := aptly_package_deleter.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder, repoPublisher.PublishRepo)
	repoCleaner := aptly_repo_cleaner.New(packageDeleter.DeletePackagesByKey, packageLister.ListPackages)

	writer := os.Stdout
	err := do(writer, repoCleaner, *urlPtr, *apiUserPtr, *passwordPtr, *passwordFilePtr, *repoPtr, *distributionPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repo_cleaner aptly_repo_cleaner.RepoCleaner, url string, user string, password string, passwordfile string, repo string, distribution string) error {
	if len(passwordfile) > 0 {
		content, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			return err
		}
		password = strings.TrimSpace(string(content))
	}
	if len(url) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_API_URL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_REPO)
	}
	return repo_cleaner.CleanRepo(aptly_api.New(url, user, password), aptly_repository.Repository(repo), aptly_distribution.Distribution(distribution))
}
