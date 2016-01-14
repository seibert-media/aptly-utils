package package_uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	"github.com/bborbe/aptly_utils/package_name"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type PublishRepo func(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error

type PackageUploader interface {
	UploadPackageByFile(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		file string) error
	UploadPackageByReader(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		packageName package_name.PackageName,
		src io.Reader) error
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

func (p *packageUploader) UploadPackageByFile(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	file string) error {
	logger.Debugf("UploadPackageByFile - repo: %s file: %s", repository, file)
	name := package_name.FromFileName(file)
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	return p.UploadPackageByReader(url, user, password, repository, distribution, name, fh)
}

func (p *packageUploader) UploadPackageByReader(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	packageName package_name.PackageName,
	src io.Reader) error {
	logger.Debugf("UploadPackageByReader - repo: %s package: %s", repository, packageName)
	if err := p.uploadFile(url, user, password, packageName, src); err != nil {
		return err
	}
	if err := p.addPackageToRepo(url, user, password, repository, packageName); err != nil {
		return err
	}
	if err := p.publishRepo(url, user, password, repository, distribution); err != nil {
		return err
	}
	return nil
}

func (p *packageUploader) uploadFile(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	packageName package_name.PackageName,
	src io.Reader) error {
	logger.Debugf("uploadFile - package: %s", packageName)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/files/%s", url, packageName))
	requestbuilder.AddBasicAuth(string(user), string(password))
	requestbuilder.SetMethod("POST")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", string(packageName))
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

func (p *packageUploader) addPackageToRepo(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	packageName package_name.PackageName) error {
	logger.Debugf("addPackageToRepo - repo: %s package: %s", repository, packageName)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", url, repository, packageName))
	requestbuilder.AddBasicAuth(string(user), string(password))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
