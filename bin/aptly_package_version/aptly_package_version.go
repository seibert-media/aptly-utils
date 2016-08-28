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
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

const (
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepoURL         = "repo-url"
	parameterRepo            = "repo"
	parameterName            = "name"
)

var (
	apiURLPtr          = flag.String(parameterAPIURL, "", "url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	repoPtr            = flag.String(parameterRepo, "", "repo")
	namePtr            = flag.String(parameterName, "", "name")
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	packageLister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)
	packageDetailLister := aptly_model_lister.New(packageLister.ListPackages)
	packageVersion := aptly_package_versions.New(packageDetailLister.ListPackageDetails)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		packageVersion,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
		*namePtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	writer io.Writer,
	packageVersions aptly_package_versions.PackageVersions,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	repo string,
	name string,
) error {
	if len(apiPasswordfile) > 0 {
		content, err := ioutil.ReadFile(apiPasswordfile)
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
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", parameterName)
	}

	var err error
	var versions []aptly_version.Version
	if versions, err = packageVersions.PackageVersions(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(repo), aptly_model.Package(name)); err != nil {
		return err
	}
	if len(versions) == 0 {
		return fmt.Errorf("package %s not found", name)
	}
	sort.Sort(aptly_version.VersionByName(versions))
	fmt.Fprintf(writer, "%s\n", versions[len(versions)-1])
	return nil
}
