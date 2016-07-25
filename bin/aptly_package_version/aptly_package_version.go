package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_model_lister "github.com/bborbe/aptly_utils/package_detail_lister"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_package_versions "github.com/bborbe/aptly_utils/package_versions"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
	aptly_version "github.com/bborbe/version"
)

const (
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
	PARAMETER_NAME              = "name"
)

var (
	logger          = log.DefaultLogger
	logLevelPtr     = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	urlPtr          = flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr      = flag.String(PARAMETER_API_USER, "", "user")
	passwordPtr     = flag.String(PARAMETER_API_PASSWORD, "", "password")
	passwordFilePtr = flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr         = flag.String(PARAMETER_REPO, "", "repo")
	namePtr         = flag.String(PARAMETER_NAME, "", "name")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	packageLister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)
	packageDetailLister := aptly_model_lister.New(packageLister.ListPackages)
	packageVersion := aptly_package_versions.New(packageDetailLister.ListPackageDetails)

	writer := os.Stdout
	err := do(writer, packageVersion, *urlPtr, *apiUserPtr, *passwordPtr, *passwordFilePtr, *repoPtr, *namePtr)
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
		password = strings.TrimSpace(string(content))
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
	var versions []aptly_version.Version
	if versions, err = packageVersions.PackageVersions(aptly_model.NewApi(url, user, password), aptly_model.Repository(repo), aptly_model.Package(name)); err != nil {
		return err
	}
	if len(versions) == 0 {
		return fmt.Errorf("package %s not found", name)
	}
	sort.Sort(aptly_version.VersionByName(versions))
	fmt.Fprintf(writer, "%s\n", versions[len(versions)-1])
	return nil
}
