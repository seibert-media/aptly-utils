package package_uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"github.com/bborbe/log"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/aptly/requestbuilder_executor"
	"github.com/bborbe/aptly/defaults"
)

type PackageUploader interface {
	UploadPackage(apiUrl string, apiUsername string, apiPassword string, file string, repo string) error
}

type packageUploader struct {
	buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider) *packageUploader {
	p := new(packageUploader)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (p *packageUploader) UploadPackage(apiUrl string, apiUsername string, apiPassword string, file string, repo string) error {
	logger.Debugf("UploadPackage")
	if err := p.uploadFile(apiUrl, apiUsername, apiPassword, file); err != nil {
		return err
	}
	if err := p.addPackageToRepo(apiUrl, apiUsername, apiPassword, file, repo); err != nil {
		return err
	}
	if err := p.publishRepo(apiUrl, apiUsername, apiPassword, file, repo, defaults.DEFAULT_DISTRIBUTION); err != nil {
		return err
	}
	return nil
}

func extractNameOfFile(path string) string {
	slashPos := strings.LastIndex(path, "/")
	if slashPos != -1 {
		return path[slashPos + 1:]
	}
	return path
}

func (p *packageUploader) uploadFile(apiUrl string, apiUsername string, apiPassword string, file string) error {
	logger.Debugf("uploadFile")
	name := extractNameOfFile(file)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/files/%s", apiUrl, name))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", fmt.Sprintf("%s.deb", name))
	if err != nil {
		return err
	}
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	bodyWriter.Close()
	requestbuilder.AddContentType(bodyWriter.FormDataContentType())
	requestbuilder.SetBody(bodyBuf)
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) addPackageToRepo(apiUrl string, apiUsername string, apiPassword string, file string, repo string) error {
	logger.Debugf("addPackageToRepo")
	name := extractNameOfFile(file)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", apiUrl, repo, name))
	requestbuilder.AddBasicAuth(apiUsername, apiPassword)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) publishRepo(apiUrl string, apiUsername string, apiPassword string, file string, repo string, distribution string) error {
	logger.Debugf("publishRepo")
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
