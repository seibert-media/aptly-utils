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
	PARAMETER_REPO              = "repo"
)

var (
	logger          = log.DefaultLogger
	logLevelPtr     = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	urlPtr          = flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr      = flag.String(PARAMETER_API_USER, "", "user")
	passwordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	passwordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
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

	writer := os.Stdout
	err := do(writer, repo_lister, *urlPtr, *apiUserPtr, *passwordPtr, *passwordFilePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repoLister aptly_repo_lister.RepoLister, url string, user string, password string, passwordfile string) error {
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
	var err error
	var repos []map[string]string
	if repos, err = repoLister.ListRepos(aptly_api.New(url, user, password)); err != nil {
		return err
	}
	for _, info := range repos {
		name := info["Name"]
		fmt.Fprintf(writer, "%s\n", name)
	}
	return nil
}
