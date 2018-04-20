package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"strings"

	aptly_model "github.com/seibert-media/aptly-utils/model"
	aptly_repo_lister "github.com/seibert-media/aptly-utils/repo_lister"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/io/util"
	"github.com/golang/glog"
)

const (
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterRepoURL         = "repo-url"
)

var (
	apiURLPtr          = flag.String(parameterAPIURL, "", "api url")
	repoURLPtr         = flag.String(parameterRepoURL, "", "repo url")
	apiUserPtr         = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr     = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr = flag.String(parameterAPIPasswordFile, "", "passwordfile")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	repo_lister := aptly_repo_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		repo_lister,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	writer io.Writer,
	repoLister aptly_repo_lister.RepoLister,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
) error {
	glog.Infof("repoURL: %v apiURL: %v apiUsername: %v apiPassword: %v apiPasswordfile: %v", repoURL, apiURL, apiUsername, apiPassword, apiPasswordfile)
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
	var err error
	var repos []map[string]string
	if repos, err = repoLister.ListRepos(aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword)); err != nil {
		return err
	}
	for _, repo := range repos {
		glog.V(2).Infof("repo: %v", repo)
		name := repo["Name"]
		fmt.Fprintf(writer, "%s\n", name)
	}
	return nil
}
