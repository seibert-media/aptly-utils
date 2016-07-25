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
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

const (
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_REPO_URL          = "repo-url"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
)

var (
	logger             = log.DefaultLogger
	logLevelPtr        = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	repoUrlPtr         = flag.String(PARAMETER_REPO_URL, "", "repo url")
	apiUrlPtr          = flag.String(PARAMETER_API_URL, "", "api url")
	apiUserPtr         = flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr            = flag.String(PARAMETER_REPO, "", "repo")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	package_lister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)

	if len(*repoUrlPtr) == 0 {
		*repoUrlPtr = *apiUrlPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		package_lister,
		*repoUrlPtr,
		*apiUrlPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer,
	packageLister aptly_package_lister.PackageLister,
	repoUrl string,
	apiUrl string,
	apiUsername string,
	apiPassword string,
	passwordfile string,
	repo string,
) error {
	if len(passwordfile) > 0 {
		content, err := ioutil.ReadFile(passwordfile)
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
	var err error
	var packages []map[string]string
	if packages, err = packageLister.ListPackages(aptly_model.NewApi(repoUrl, apiUrl, apiUsername, apiPassword), aptly_model.Repository(repo)); err != nil {
		return err
	}
	for _, info := range packages {
		name := info["Package"]
		version := info["Version"]
		fmt.Fprintf(writer, "%s %s\n", name, version)
	}
	return nil
}
