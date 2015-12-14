package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	aptly_repo_creater "github.com/bborbe/aptly_utils/repo_creater"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
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
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	requestbuilder_executor := aptly_requestbuilder_executor.New(client.GetClientWithoutProxy())
	repo_creater := aptly_repo_creater.New(requestbuilder_executor, requestbuilder.NewHttpRequestBuilderProvider())

	writer := os.Stdout
	err := do(writer, repo_creater, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *repoPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repo_creater aptly_repo_creater.RepoCreater, url string, user string, password string, passwordfile string, repo string) error {
	if len(passwordfile) > 0 {
		content, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			return err
		}
		password = string(content)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter repo missing")
	}
	return repo_creater.CreateRepo(url, user, password, repo)
}
