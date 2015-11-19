package main

import (
	"flag"
	"os"

	"runtime"

	"io"

	"fmt"

	"net/http"

	aptly_package_uploader "github.com/bborbe/aptly/package_uploader"
	"github.com/bborbe/http/client"
	"github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_FILE         = "file"
	PARAMETER_LOGLEVEL     = "loglevel"
	PARAMETER_API_URL      = "url"
	PARAMETER_API_USER     = "username"
	PARAMETER_API_PASSWORD = "password"
	PARAMETER_REPO         = "repo"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	filePtr := flag.String(PARAMETER_FILE, "", "file")
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	package_uploader := aptly_package_uploader.New(func() *http.Client {
		return client.GetClientWithoutProxy()
	}, requestbuilder.NewHttpRequestBuilderProvider().NewHttpRequestBuilder)

	writer := os.Stdout
	err := do(writer, package_uploader, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *filePtr, *repoPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, package_uploader aptly_package_uploader.PackageUploader, url string, user string, pass string, file string, repo string) error {
	if len(file) == 0 {
		return fmt.Errorf("parameter file missing")
	}
	return package_uploader.UploadPackage(url, user, pass, file, repo)
}
