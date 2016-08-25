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
	parameterLoglevel        = "loglevel"
	parameterRepoURL         = "repo-url"
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepo            = "repo"
)

var (
	logger             = log.DefaultLogger
	logLevelPtr        = flag.String(parameterLoglevel, log.INFO_STRING, log.FLAG_USAGE)
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
	apiURLPtr          = flag.String(parameterAPIURL, "", "api url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	repoPtr            = flag.String(parameterRepo, "", "repo")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	package_lister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		package_lister,
		*repoURLPtr,
		*apiURLPtr,
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
	repoURL string,
	apiURL string,
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
	if len(apiURL) == 0 {
		return fmt.Errorf("parameter %s missing", parameterAPIURL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", parameterRepo)
	}
	var err error
	var packages []map[string]string
	if packages, err = packageLister.ListPackages(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo)); err != nil {
		return err
	}
	for _, info := range packages {
		name := info["Package"]
		version := info["Version"]
		fmt.Fprintf(writer, "%s %s\n", name, version)
	}
	return nil
}
