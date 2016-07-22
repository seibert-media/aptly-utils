package package_uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"strings"

	"io/ioutil"

	aptly_api "github.com/bborbe/aptly_utils/api"
	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PublishRepo func(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution) error

type PackageUploader interface {
	UploadPackageByFile(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, file string) error
	UploadPackageByReader(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, file string, src io.Reader) error
}

type packageUploader struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	publishRepo                PublishRepo
}

var logger = log.DefaultLogger

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, publishRepo PublishRepo) *packageUploader {
	p := new(packageUploader)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.publishRepo = publishRepo
	return p
}

func FromFileName(path string) string {
	slashPos := strings.LastIndex(path, "/")
	if slashPos != -1 {
		return path[slashPos+1:]
	}
	return path
}

func (p *packageUploader) UploadPackageByFile(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, file string) error {
	logger.Debugf("UploadPackageByFile - repo: %s file: %s", repository, file)
	name := FromFileName(file)
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	return p.UploadPackageByReader(api, repository, distribution, name, fh)
}

func (p *packageUploader) UploadPackageByReader(api aptly_api.Api, repository aptly_repository.Repository, distribution aptly_distribution.Distribution, file string, src io.Reader) error {
	logger.Debugf("UploadPackageByReader - repo: %s dist: %s file: %s", repository, distribution, file)
	if err := p.uploadFile(api, file, src); err != nil {
		return err
	}
	if err := p.addPackageToRepo(api, repository, file); err != nil {
		return err
	}
	if err := p.publishRepo(api, repository, distribution); err != nil {
		return err
	}
	return nil
}

func (p *packageUploader) uploadFile(api aptly_api.Api, file string, src io.Reader) error {
	logger.Debugf("uploadFile - package: %s", file)

	logger.Debugf("write upload to temp file ...")
	f, err := ioutil.TempFile("", "upload")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	bodyWriter := multipart.NewWriter(f)
	fileWriter, err := bodyWriter.CreateFormFile("file", file)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return err
	}
	bodyWriter.Close()
	f.Seek(0, 0)

	logger.Debugf("write upload to temp file finshed")
	logger.Debugf("build request ...")
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/files/%s", api.Url, file))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType(bodyWriter.FormDataContentType())
	requestbuilder.SetBody(f)
	logger.Debugf("build request finished")
	logger.Debugf("uploading ...")
	defer logger.Debugf("uploading finished")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}

func (p *packageUploader) addPackageToRepo(api aptly_api.Api, repository aptly_repository.Repository, file string) error {
	logger.Debugf("addPackageToRepo - repo: %s file: %s", repository, file)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", api.Url, repository, file))
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	logger.Debugf("addPackageToRepo ...")
	defer logger.Debugf("addPackageToRepo finished")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
