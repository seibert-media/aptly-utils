package package_deleter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aptly_model "github.com/bborbe/aptly_utils/model"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	aptly_version "github.com/bborbe/version"
	"github.com/golang/glog"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type NewHTTPRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PublishRepo func(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution) error

type PackageDeleter interface {
	DeletePackageByNameAndVersion(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, name aptly_model.Package, version aptly_version.Version) error
	DeletePackagesByKey(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, keys []aptly_model.Key) error
}

type packageDeleter struct {
	executeRequest        ExecuteRequest
	newHTTPRequestBuilder NewHTTPRequestBuilder
	publishRepo           PublishRepo
}

func New(executeRequest ExecuteRequest, newHTTPRequestBuilder NewHTTPRequestBuilder, publishRepo PublishRepo) *packageDeleter {
	p := new(packageDeleter)
	p.executeRequest = executeRequest
	p.newHTTPRequestBuilder = newHTTPRequestBuilder
	p.publishRepo = publishRepo
	return p
}

type JSONStruct []map[string]string

func (p *packageDeleter) DeletePackageByNameAndVersion(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, name aptly_model.Package, version aptly_version.Version) error {
	keys, err := p.findKeys(api, repository, name, version)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("no package found")
	}
	return p.DeletePackagesByKey(api, repository, distribution, keys)
}

func (p *packageDeleter) DeletePackagesByKey(api aptly_model.API, repository aptly_model.Repository, distribution aptly_model.Distribution, keys []aptly_model.Key) error {
	if err := p.deletePackagesByKey(api, repository, keys); err != nil {
		return err
	}
	return p.publishRepo(api, repository, distribution)
}

func (p *packageDeleter) findKeys(api aptly_model.API, repository aptly_model.Repository, packageName aptly_model.Package, version aptly_version.Version) ([]aptly_model.Key, error) {
	glog.V(2).Infof("PackageVersions - repo: %s package: %s", repository, packageName)
	url := fmt.Sprintf("%s/api/repos/%s/packages?format=details", api.APIUrl, repository)
	requestbuilder := p.newHTTPRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request to %s failed with status %d", url, resp.StatusCode)
	}

	var jsonStruct JSONStruct
	err = json.Unmarshal(content, &jsonStruct)
	if err != nil {
		return nil, err
	}
	var keys []aptly_model.Key
	for _, info := range jsonStruct {
		if info["Package"] == string(packageName) && info["Version"] == string(version) {
			key := info["Key"]
			glog.V(2).Infof("found key: %s", key)
			keys = append(keys, aptly_model.Key(key))
		}
	}
	return keys, nil
}

func (p *packageDeleter) deletePackagesByKey(api aptly_model.API, repository aptly_model.Repository, keys []aptly_model.Key) error {
	glog.V(2).Infof("delete package with keys: %v", keys)
	url := fmt.Sprintf("%s/api/repos/%s/packages?format=details", api.APIUrl, repository)
	requestbuilder := p.newHTTPRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.APIUsername), string(api.APIPassword))
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestContent, err := json.Marshal(map[string][]aptly_model.Key{"PackageRefs": keys})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(requestContent))
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request to %s failed with status %d", url, resp.StatusCode)
	}
	return nil
}
