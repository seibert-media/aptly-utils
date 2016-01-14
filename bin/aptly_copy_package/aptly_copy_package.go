package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_package_latest_version "github.com/bborbe/aptly_utils/package_latest_version"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_package_versions "github.com/bborbe/aptly_utils/package_versions"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	http_client "github.com/bborbe/http/client"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_SOURCE            = "source"
	PARAMETER_TARGET            = "target"
	PARAMETER_NAME              = "name"
	PARAMETER_VERSION           = "version"
	PARAMETER_LOGLEVEL          = "loglevel"
	PARAMETER_API_URL           = "url"
	PARAMETER_API_USER          = "username"
	PARAMETER_API_PASSWORD      = "password"
	PARAMETER_API_PASSWORD_FILE = "passwordfile"
	PARAMETER_REPO              = "repo"
	PARAMETER_DISTRIBUTION      = "distribution"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	sourcePtr := flag.String(PARAMETER_SOURCE, "", "source")
	targetPtr := flag.String(PARAMETER_TARGET, "", "target")
	namePtr := flag.String(PARAMETER_NAME, "", "name")
	versionPtr := flag.String(PARAMETER_VERSION, "", "version")
	targetDistributionPtr := flag.String(PARAMETER_DISTRIBUTION, string(aptly_distribution.DEFAULT), "distribution")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	requestbuilder_executor := aptly_requestbuilder_executor.New(client)
	requestbuilder := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	packageUploader := aptly_package_uploader.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishRepo)
	packageCopier := aptly_package_copier.New(packageUploader, requestbuilder, client)
	packageLister := aptly_package_lister.New(client.Do, httpRequestBuilderProvider.NewHttpRequestBuilder)
	packageVersion := aptly_package_versions.New(packageLister.ListPackages)
	packageLastestVersion := aptly_package_latest_version.New(packageVersion.PackageVersions)

	writer := os.Stdout
	err := do(writer, packageCopier, packageLastestVersion, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *sourcePtr, *targetPtr, *targetDistributionPtr, *namePtr, *versionPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, packageCopier aptly_package_copier.PackageCopier, packageLatestVersion aptly_package_latest_version.PackageLatestVersion, url string, user string, password string, passwordfile string, sourceRepo string, targetRepo string, targetDistribution, name string, version string) error {
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
	if len(sourceRepo) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_SOURCE)
	}
	if len(targetRepo) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_TARGET)
	}
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_NAME)
	}
	if len(version) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_VERSION)
	}
	if version == "latest" {
		latestVersion, err := packageLatestVersion.PackageLatestVersion(aptly_url.Url(url), aptly_user.User(user), aptly_password.Password(password), aptly_repository.Repository(sourceRepo), package_name.PackageName(name))
		if err != nil {
			return err
		}
		version = string(*latestVersion)
	}
	if name == "all" {

	}
	return packageCopier.CopyPackage(aptly_url.Url(url), aptly_user.User(user), aptly_password.Password(password), aptly_repository.Repository(sourceRepo), aptly_repository.Repository(targetRepo), aptly_distribution.Distribution(targetDistribution), package_name.PackageName(name), aptly_version.Version(version))
}
