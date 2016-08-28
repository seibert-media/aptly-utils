package package_uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"strings"

	"io/ioutil"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

type PublishRepo func(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error

type PackageUploader interface {
	UploadPackageByFile(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, file string) error
	UploadPackageByReader(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, file string, src io.Reader) error
}

type packageUploader struct {
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider
	publishRepo                PublishRepo
}

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HTTPRequestBuilderProvider, publishRepo PublishRepo) *packageUploader {
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

func (p *packageUploader) UploadPackageByFile(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, file string) error {
	glog.V(2).Infof("UploadPackageByFile - repo: %s file: %s", repository, file)
	name := FromFileName(file)
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	return p.UploadPackageByReader(api, repository, distribution, name, fh)
}

func (p *packageUploader) UploadPackageByReader(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, file string, src io.Reader) error {
	glog.V(2).Infof("UploadPackageByReader - repo: %s dist: %s file: %s", repository, distribution, file)
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

func (p *packageUploader) uploadFile(api aptly_model.API, file string, src io.Reader) error {
	glog.V(2).Infof("uploadFile - package: %s", file)

	glog.V(2).Infof("write upload to temp file ...")
	f, err := ioutil.TempFile("", "upload")
	if err != nil {
		return err
	}
	defer f.Close()
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			glog.Warningf("remove %s failed: %v", f.Name(), err)
		}
	}()
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
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	glog.V(2).Infof("write upload to temp file finshed")
	glog.V(2).Infof("build request ...")
	requestbuilder := p.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/files/%s", api.APIUrl, file))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType(bodyWriter.FormDataContentType())
	requestbuilder.SetContentLength(fileInfo.Size())
	requestbuilder.SetBody(f)
	glog.V(2).Infof("build request finished")
	glog.V(2).Infof("uploading ...")
	if err := p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder); err != nil {
		return err
	}
	glog.V(2).Infof("uploading finished")
	return nil
}

func (p *packageUploader) addPackageToRepo(api aptly_model.API, repository aptly_model.Repository, file string) error {
	glog.V(2).Infof("addPackageToRepo - repo: %s file: %s", repository, file)
	requestbuilder := p.httpRequestBuilderProvider.NewHTTPRequestBuilder(fmt.Sprintf("%s/api/repos/%s/file/%s?forceReplace=1", api.APIUrl, repository, file))
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	glog.V(2).Infof("addPackageToRepo ...")
	if err := p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder); err != nil {
		return err
	}
	glog.V(2).Infof("addPackageToRepo finished")
	return nil
}
