package package_uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/bborbe/aptly/defaults"
	"github.com/bborbe/aptly/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PackageUploader interface {
	UploadPackageByFile(apiUrl string, apiUsername string, apiPassword string, repo string, file string) error
	UploadPackageByReader(apiUrl string, apiUsername string, apiPassword string, repo string, name string, src io.Reader) error
}

type packageUploader struct {
	buildRequestAndExecute     requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *packageUploader {
	p := new(packageUploader)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (p *packageUploader) UploadPackageByFile(apiUrl string, apiUsername string, apiPassword string, repo string, file string) error {
	logger.Debugf("UploadPackageByFile - repo: %s file: %s", repo, file)
	name := extractPkgOfFile(file)
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	return p.UploadPackageByReader(apiUrl, apiUsername, apiPassword, repo, name, fh)
}

func (p *packageUploader) UploadPackageByReader(apiUrl string, apiUsername string, apiPassword string, repo string, pkg string, src io.Reader) error {
	logger.Debugf("UploadPackageByReader - repo: %s package: %s", repo, pkg)
	if err := p.uploadFile(apiUrl, apiUsername, apiPassword, pkg, src); err != nil {
		return err
	}
	if err := p.addPackageToRepo(apiUrl, apiUsername, apiPassword, repo, pkg); err != nil {
		return err
	}
	if err := p.publishRepo(apiUrl, apiUsername, apiPassword, repo, defaults.DEFAULT_DISTRIBUTION); err != nil {
		return err
	}
	return nil
}

func extractPkgOfFile(path string) string {
	slashPos := strings.LastIndex(path, "/")
	if slashPos != -1 {
		return path[slashPos+1:]
	}
	return path
}

func (p *packageUploader) uploadFile(apiUrl string, apiUsername string, apiPassword string, pkg string, src io.Reader) error {
	logger.Debugf("uploadFile - package: %s", pkg)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/files/%s", apiUrl, pkg))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", pkg)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return err
	}
	bodyWriter.Close()
	requestbuilder.AddContentType(bodyWriter.FormDataContentType())
	requestbuilder.SetBody(bodyBuf)
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) addPackageToRepo(apiUrl string, apiUsername string, apiPassword string, repo string, pkg string) error {
	logger.Debugf("addPackageToRepo - repo: %s package: %s", repo, pkg)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", apiUrl, repo, pkg))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) publishRepo(apiUrl string, apiUsername string, apiPassword string, repo string, distribution string) error {
	logger.Debugf("publishRepo - repo: %s distribution: %s", repo, distribution)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/publish/%s/%s", apiUrl, repo, distribution))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("PUT")
	requestbuilder.AddContentType("application/json")
	content, err := json.Marshal(map[string]bool{"ForceOverwrite": true})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
