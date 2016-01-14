package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	aptly_architecture "github.com/bborbe/aptly_utils/architecture"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	http_client "github.com/bborbe/http/client"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

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

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	distributionPtr := flag.String(PARAMETER_DISTRIBUTION, string(aptly_distribution.DEFAULT), "distribution")
	architecturePtr := flag.String(PARAMETER_ARCHITECTURE, string(aptly_architecture.DEFAULT), "architecture")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	requestbuilder_executor := aptly_requestbuilder_executor.New(client)
	requestbuilder := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	repo_creater := aptly_repo_creater.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishNewRepo)

	writer := os.Stdout
	err := do(writer, repo_creater, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *repoPtr, *distributionPtr, strings.Split(*architecturePtr, ","))
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
	return repo_creater.CreateRepo(aptly_url.Url(url), aptly_user.User(user), aptly_password.Password(password), aptly_repository.Repository(repo), aptly_distribution.Distribution(distribution), aptly_architecture.Parse(architectures))
}
