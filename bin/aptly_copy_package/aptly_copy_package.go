package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	aptly_package_uploader "github.com/bborbe/aptly_utils/package_uploader"
	aptly_repo_publisher "github.com/bborbe/aptly_utils/repo_publisher"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_client "github.com/bborbe/http/client"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
"strings"
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
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	requestbuilder_executor := aptly_requestbuilder_executor.New(client)
	requestbuilder := http_requestbuilder.NewHttpRequestBuilderProvider()
	repo_publisher := aptly_repo_publisher.New(requestbuilder_executor, requestbuilder)
	package_uploader := aptly_package_uploader.New(requestbuilder_executor, requestbuilder, repo_publisher.PublishRepo)
	package_copier := aptly_package_copier.New(package_uploader, requestbuilder, client)

	writer := os.Stdout
	err := do(writer, package_copier, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *sourcePtr, *targetPtr, *namePtr, *versionPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, package_copier aptly_package_copier.PackageCopier, url string, user string, password string, passwordfile string, source string, target string, name string, version string) error {
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
	if len(source) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_SOURCE)
	}
	if len(target) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_TARGET)
	}
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_NAME)
	}
	if len(version) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_VERSION)
	}
	return package_copier.CopyPackage(url, user, password, source, target, name, version)
}
