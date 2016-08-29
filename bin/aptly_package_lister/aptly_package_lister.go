package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"

	"strings"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

const (
	parameterRepoURL         = "repo-url"
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepo            = "repo"
)

var (
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
	apiURLPtr          = flag.String(parameterAPIURL, "", "api url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	repoPtr            = flag.String(parameterRepo, "", "repo")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	package_lister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	err := do(
		package_lister,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*repoPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	packageLister aptly_package_lister.PackageLister,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	passwordfile string,
	repo string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v passwordfile: %v repo: %v", repoURL, apiURL, apiUsername, apiPassword, passwordfile, repo)
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
