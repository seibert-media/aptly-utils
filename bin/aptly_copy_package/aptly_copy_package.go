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
	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_model_latest_lister "github.com/bborbe/aptly_utils/package_detail_latest_lister"
	aptly_model_lister "github.com/bborbe/aptly_utils/package_detail_lister"
	aptly_package_latest_version "github.com/bborbe/aptly_utils/package_latest_version"
	aptly_package_lister "github.com/bborbe/aptly_utils/package_lister"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_package_versions "github.com/bborbe/aptly_utils/package_versions"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

const (
	parameterSource          = "source"
	parameterTarget          = "target"
	parameterName            = "name"
	parameterVersion         = "version"
	parameterAPIURL          = "url"
	parameterAPIUser         = "username"
	parameterAPIPassword     = "password"
	parameterAPIPasswordFile = "passwordfile"
	parameterDistribution    = "distribution"
	parameterRepoURL         = "repo-url"
)

var (
	apiURLPtr             = flag.String(parameterAPIURL, "", "url")
	apiUserPtr            = flag.String(parameterAPIUser, "", "user")
	apiPasswordPtr        = flag.String(parameterAPIPassword, "", "password")
	apiPasswordFilePtr    = flag.String(parameterAPIPasswordFile, "", "passwordfile")
	sourcePtr             = flag.String(parameterSource, "", "source")
	targetPtr             = flag.String(parameterTarget, "", "target")
	namePtr               = flag.String(parameterName, "", "name")
	versionPtr            = flag.String(parameterVersion, "", "version")
	targetDistributionPtr = flag.String(parameterDistribution, string(aptly_model.DistribuionDefault), "distribution")
	repoURLPtr            = flag.String(parameterRepoURL, "", "repo url")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClient := http_client_builder.New().WithoutProxy().Build()
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	requestbuilder_executor := aptly_requestbuilder_executor.New(httpClient.Do)
	requestbuilder := http_requestbuilder.NewHTTPRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	packageUploader := aptly_package_uploader.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishRepo)
	packageCopier := aptly_package_copier.New(packageUploader, requestbuilder, httpClient.Do)
	packageLister := aptly_package_lister.New(httpClient.Do, httpRequestBuilderProvider.NewHTTPRequestBuilder)
	packageDetailLister := aptly_model_lister.New(packageLister.ListPackages)
	packageVersion := aptly_package_versions.New(packageDetailLister.ListPackageDetails)
	packageLastestVersion := aptly_package_latest_version.New(packageVersion.PackageVersions)
	packageDetailLatestLister := aptly_model_latest_lister.New(packageDetailLister.ListPackageDetails)

	if len(*repoURLPtr) == 0 {
		*repoURLPtr = *apiURLPtr
	}

	writer := os.Stdout
	err := do(
		writer,
		packageCopier,
		packageLastestVersion,
		packageDetailLatestLister,
		*repoURLPtr,
		*apiURLPtr,
		*apiUserPtr,
		*apiPasswordPtr,
		*apiPasswordFilePtr,
		*sourcePtr,
		*targetPtr,
		*targetDistributionPtr,
		*namePtr,
		*versionPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	writer io.Writer,
	packageCopier aptly_package_copier.PackageCopier,
	packageLatestVersion aptly_package_latest_version.PackageLatestVersion,
	packageDetailLatestLister aptly_model_latest_lister.PackageDetailLatestLister,
	repoURL string,
	apiURL string,
	apiUsername string,
	apiPassword string,
	apiPasswordfile string,
	sourceRepo string,
	targetRepo string,
	targetDistribution,
	name string,
	version string,
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
	if len(sourceRepo) == 0 {
		return fmt.Errorf("parameter %s missing", parameterSource)
	}
	if len(targetRepo) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTarget)
	}
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", parameterName)
	}
	if len(version) == 0 {
		return fmt.Errorf("parameter %s missing", parameterVersion)
	}
	return copy(packageCopier, packageLatestVersion, packageDetailLatestLister, aptly_model.NewAPI(repoURL, apiURL, apiUsername, apiPassword), aptly_model.Repository(sourceRepo), aptly_model.Repository(targetRepo), aptly_model.Distribution(targetDistribution), aptly_model.Package(name), aptly_version.Version(version))
}

func copy(packageCopier aptly_package_copier.PackageCopier, packageLatestVersion aptly_package_latest_version.PackageLatestVersion, packageDetailLatestLister aptly_model_latest_lister.PackageDetailLatestLister, api aptly_model.API, sourceRepo aptly_model.Repository, targetRepo aptly_model.Repository, targetDistribution aptly_model.Distribution, packageName aptly_model.Package, version aptly_version.Version) error {
	if packageName == aptly_model.PackageAll && version != aptly_version.LATEST {
		return fmt.Errorf("can't copy with package all and version != latest")
	}
	var list []aptly_model.PackageDetail
	var err error
	if packageName == aptly_model.PackageAll {
		list, err = packageDetailLatestLister.ListLatestPackageDetails(api, sourceRepo)
		if err != nil {
			return err
		}
	} else {
		list = []aptly_model.PackageDetail{aptly_model.NewPackageDetail(packageName, version)}
	}
	return copyList(packageCopier, packageLatestVersion, api, sourceRepo, targetRepo, targetDistribution, list)
}

func copyList(packageCopier aptly_package_copier.PackageCopier, packageLatestVersion aptly_package_latest_version.PackageLatestVersion, api aptly_model.API, sourceRepo aptly_model.Repository, targetRepo aptly_model.Repository, targetDistribution aptly_model.Distribution, list []aptly_model.PackageDetail) error {
	for _, e := range list {
		version := e.Version
		packageName := e.Package
		if version == aptly_version.LATEST {
			latestVersion, err := packageLatestVersion.PackageLatestVersion(api, sourceRepo, packageName)
			if err != nil {
				return err
			}
			version = *latestVersion
		}
		err := packageCopier.CopyPackage(api, sourceRepo, targetRepo, targetDistribution, packageName, version)
		if err != nil {
			return err
		}
	}
	return nil
}
