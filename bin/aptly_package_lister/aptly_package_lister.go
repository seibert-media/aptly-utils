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
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
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
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	urlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	passwordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	passwordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	package_lister := aptly_package_lister.New(client.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)

	writer := os.Stdout
	err := do(writer, package_lister, *urlPtr, *apiUserPtr, *passwordPtr, *passwordFilePtr, *repoPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, packageLister aptly_package_lister.PackageLister, url string, user string, password string, passwordfile string, repo string) error {
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
	var err error
	var packages []map[string]string
	if packages, err = packageLister.ListPackages(aptly_api.New(url, user, password), aptly_repository.Repository(repo)); err != nil {
		return err
	}
	for _, info := range packages {
		name := info["Package"]
		version := info["Version"]
		fmt.Fprintf(writer, "%s %s\n", name, version)
	}
	return nil
}
