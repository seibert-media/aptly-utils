package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"

	"strings"

	"io"
	"os"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_package_lister "github.com/seibert-media/aptly-utils/package_lister"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/io/util"
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
		glog.Exit(err)
	}
}

func do(
	writer io.Writer,
	packageLister aptly_package_lister.PackageLister,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	repo string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v passwordfile: %v repo: %v", repoURL, apiURL, apiUsername, apiPassword, apiPasswordfile, repo)
	if len(apiPasswordfile) > 0 {
		apiPasswordfile, err := util.NormalizePath(apiPasswordfile)
		if err != nil {
			return fmt.Errorf("normalize path %s failed: %v", apiPasswordfile, err)
		}
		content, err := ioutil.ReadFile(apiPasswordfile)
		if err != nil {
			return fmt.Errorf("read password from file failed: %v", err)
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
