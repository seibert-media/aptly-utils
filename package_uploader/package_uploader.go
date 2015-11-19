package package_uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PackageUploader interface {
	UploadPackage(api_url string, api_username string, api_password string, file string, repo string) error
}

type packageUploader struct {
	httpClientProvider         HttpClientProvider
	httpRequestBuilderProvider HttpRequestBuilderProvider
}

type HttpClientProvider func() *http.Client

type HttpRequestBuilderProvider func(url string) http_requestbuilder.HttpRequestBuilder

var logger = log.DefaultLogger

func New(httpClientProvider HttpClientProvider, httpRequestBuilderProvider HttpRequestBuilderProvider) *packageUploader {
	p := new(packageUploader)
	p.httpClientProvider = httpClientProvider
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	return p
}

func (p *packageUploader) UploadPackage(api_url string, api_username string, api_password string, file string, repo string) error {
	logger.Debugf("UploadPackage")
	if err := p.upload_file(api_url, api_username, api_password, file); err != nil {
		return err
	}
	if err := p.add_package_to_repo(api_url, api_username, api_password, file, repo); err != nil {
		return err
	}
	if err := p.publish_repo(api_url, api_username, api_password, file, repo, "default"); err != nil {
		return err
	}
	return nil
}

func extractNameOfFile(path string) string {
	slashPos := strings.LastIndex(path, "/")
	if slashPos != -1 {
		return path[slashPos+1:]
	}
	return path
}

func (p *packageUploader) upload_file(api_url string, api_username string, api_password string, file string) error {
	logger.Debugf("upload_file")
	name := extractNameOfFile(file)
	requestbuilder := p.httpRequestBuilderProvider(fmt.Sprintf("%s/api/files/%s", api_url, name))
	requestbuilder.AddBasicAuth(api_username, api_password)
	requestbuilder.SetMethod("POST")

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", fmt.Sprintf("%s.deb", name))
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(file)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	bodyWriter.Close()

	requestbuilder.AddContentType(bodyWriter.FormDataContentType())
	requestbuilder.SetBody(bodyBuf)
	return p.buildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) add_package_to_repo(api_url string, api_username string, api_password string, file string, repo string) error {
	logger.Debugf("add_package_to_repo")
	name := extractNameOfFile(file)
	requestbuilder := p.httpRequestBuilderProvider(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", api_url, repo, name))
	requestbuilder.AddBasicAuth(api_username, api_password)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	return p.buildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) publish_repo(api_url string, api_username string, api_password string, file string, repo string, distribution string) error {
	logger.Debugf("publish_repo")
	requestbuilder := p.httpRequestBuilderProvider(fmt.Sprintf("%s/api/publish/%s/%s", api_url, repo, distribution))
	requestbuilder.AddBasicAuth(api_username, api_password)
	requestbuilder.SetMethod("PUT")
	requestbuilder.AddContentType("application/json")

	content, err := json.Marshal(map[string]bool{"ForceOverwrite": true})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(content))
	return p.buildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) buildRequestAndExecute(requestbuilder http_requestbuilder.HttpRequestBuilder) error {
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	client := p.httpClientProvider()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("upload file failed: %s", string(content))
	}
	return nil
}
