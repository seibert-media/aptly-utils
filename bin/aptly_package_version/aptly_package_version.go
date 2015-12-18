package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"sort"

	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_package_versions "github.com/bborbe/aptly_utils/package_versions"
	"github.com/bborbe/aptly_utils/version"
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
	PARAMETER_NAME              = "name"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	namePtr := flag.String(PARAMETER_NAME, "", "name")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	package_lister := aptly_package_lister.New(client.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)
	package_version := aptly_package_versions.New(package_lister.ListPackages)

	writer := os.Stdout
	err := do(writer, package_version, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *repoPtr, *namePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, packageVersions aptly_package_versions.PackageVersions, url string, user string, password string, passwordfile string, repo string, name string) error {
	if len(passwordfile) > 0 {
		content, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			return err
		}
		password = string(content)
	}

	if len(url) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_API_URL)
	}
	if len(repo) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_REPO)
	}
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_NAME)
	}

	var err error
	var versions []version.Version
	if versions, err = packageVersions.PackageVersions(url, user, password, repo, name); err != nil {
		return err
	}
	if len(versions) == 0 {
		return fmt.Errorf("package %s not found", name)
	}
	sort.Sort(version.VersionByName(versions))
	fmt.Fprintf(writer, "%s\n", versions[len(versions)-1])
	return nil
}
