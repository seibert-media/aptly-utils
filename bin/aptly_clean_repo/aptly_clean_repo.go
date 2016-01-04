package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	aptly_package_deleter "github.com/bborbe/aptly_utils/package_deleter"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_repo_cleaner "github.com/bborbe/aptly_utils/repo_cleaner"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client "github.com/bborbe/http/client"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
"strings"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	package_lister := aptly_package_lister.New(client.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)
	requestbuilder_executor := aptly_requestbuilder_executor.New(client)
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, httpRequestBuilderProvider)
	package_deleter := aptly_package_deleter.New(client.Do, httpRequestBuilderProvider.NewHttpRequestBuilder, repo_publisher.PublishRepo)
	repo_cleaner := aptly_repo_cleaner.New(package_deleter.DeletePackagesByKey, package_lister.ListPackages)

	writer := os.Stdout
	err := do(writer, repo_cleaner, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *repoPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repo_cleaner aptly_repo_cleaner.RepoCleaner, url string, user string, password string, passwordfile string, repo string) error {
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
	return repo_cleaner.CleanRepo(url, user, password, repo)
}
