package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	aptly_repo_deleter "github.com/bborbe/aptly_utils/repo_deleter"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_FILE              = "file"
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
	filePtr := flag.String(PARAMETER_FILE, "", "file")
	apiUrlPtr := flag.String(PARAMETER_API_URL, "", "url")
	apiUserPtr := flag.String(PARAMETER_API_USER, "", "user")
	apiPasswordPtr := flag.String(PARAMETER_API_PASSWORD, "", "password")
	apiPasswordFilePtr := flag.String(PARAMETER_API_PASSWORD_FILE, "", "passwordfile")
	repoPtr := flag.String(PARAMETER_REPO, "", "repo")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	repo_deleter := aptly_repo_deleter.New()

	writer := os.Stdout
	err := do(writer, repo_deleter, *apiUrlPtr, *apiUserPtr, *apiPasswordPtr, *apiPasswordFilePtr, *filePtr, *repoPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, repo_deleter aptly_repo_deleter.RepoDeleter, url string, user string, password string, passwordfile string, file string, repo string) error {
	if len(passwordfile) > 0 {
		content, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			return err
		}
		password = string(content)
	}
	if len(file) == 0 {
		return fmt.Errorf("parameter file missing")
	}
	return repo_deleter.DeleteRepo()
}
