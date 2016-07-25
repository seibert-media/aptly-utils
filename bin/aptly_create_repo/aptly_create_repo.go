package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

const (
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
	PARAMETER_DISTRIBUTION      = "distribution"
	PARAMETER_ARCHITECTURE      = "architecture"
)

var (
	logger          = log.DefaultLogger
	logLevelPtr     = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	urlPtr          = flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr      = flag.String(PARAMETER_API_USER, "", "user")
	passwordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	passwordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr         = flag.String(PARAMETER_REPO, "", "repo")
	distributionPtr = flag.String(PARAMETER_DISTRIBUTION, string(aptly_model.DISTRIBUTION_DEFAULT), "distribution")
	architecturePtr = flag.String(PARAMETER_ARCHITECTURE, string(aptly_model.ARCHITECTURE_DEFAULT), "architecture")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	requestbuilder := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	repo_creater := aptly_repo_creater.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishNewRepo)

	writer := os.Stdout
	err := do(writer, repo_creater, *urlPtr, *apiUserPtr, *passwordPtr, *passwordFilePtr, *repoPtr, *distributionPtr, strings.Split(*architecturePtr, ","))
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repo_creater aptly_repo_creater.RepoCreater, url string, user string, password string, passwordfile string, repo string, distribution string, architectures []string) error {
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
	return repo_creater.CreateRepo(aptly_model.NewApi(url, user, password), aptly_model.Repository(repo), aptly_model.Distribution(distribution), aptly_model.ParseArchitectures(architectures))
}
