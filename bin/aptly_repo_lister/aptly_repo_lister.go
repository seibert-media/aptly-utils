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
	aptly_repo_lister "github.com/bborbe/aptly_utils/repo_lister"
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
	PARAMETER_REPO_URL          = "repo-url"
)

var (
	logger             = log.DefaultLogger
	logLevelPtr        = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr          = flag.String(PARAMETER_API_URL, "", "api url")
	repoUrlPtr         = flag.String(PARAMETER_REPO_URL, "", "repo url")
	apiUserPtr         = flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_lister := aptly_repo_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)

	if len(*repoUrlPtr) == 0 {
		*repoUrlPtr = *apiUrlPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		repo_lister,
		*repoUrlPtr,
		*apiUrlPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(
	writer io.Writer,
	repoLister aptly_repo_lister.RepoLister,
	repoUrl string,
	apiUrl string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
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
	var err error
	var repos []map[string]string
	if repos, err = repoLister.ListRepos(aptly_model.NewApi(repoUrl, apiUrl, apiUsername, apiPassword)); err != nil {
		return err
	}
	for _, repo := range repos {
		logger.Debugf("repo: %v", repo)
		name := repo["Name"]
		fmt.Fprintf(writer, "%s\n", name)
	}
	return nil
}
