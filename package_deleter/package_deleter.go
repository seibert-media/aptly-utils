package package_deleter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aptly_distribution "github.com/bborbe/aptly_utils/distribution"
	aptly_key "github.com/bborbe/aptly_utils/key"
	aptly_package_name "github.com/bborbe/aptly_utils/package_name"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type PublishRepo func(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution) error

type PackageDeleter interface {
	DeletePackageByNameAndVersion(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repository aptly_repository.Repository,
		distribution aptly_distribution.Distribution,
		name aptly_package_name.PackageName,
		version aptly_version.Version) error
	DeletePackagesByKey(
		apiUrl aptly_url.Url,
		apiUsername aptly_user.User,
		apiPassword aptly_password.Password,
		repository aptly_repository.Repository,
		keys []aptly_key.Key) error
}

type packageDeleter struct {
	executeRequest        ExecuteRequest
	newHttpRequestBuilder NewHttpRequestBuilder
	publishRepo           PublishRepo
}

var logger = log.DefaultLogger

func New(executeRequest ExecuteRequest, newHttpRequestBuilder NewHttpRequestBuilder, publishRepo PublishRepo) *packageDeleter {
	p := new(packageDeleter)
	p.executeRequest = executeRequest
	p.newHttpRequestBuilder = newHttpRequestBuilder
	p.publishRepo = publishRepo
	return p
}

type JsonStruct []map[string]string

func (p *packageDeleter) DeletePackageByNameAndVersion(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository,
	distribution aptly_distribution.Distribution,
	name aptly_package_name.PackageName,
	version aptly_version.Version) error {
	keys, err := p.findKeys(apiUrl, apiUsername, apiPassword, repository, name, version)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("no package found")
	}
	if err = p.DeletePackagesByKey(apiUrl, apiUsername, apiPassword, repository, keys); err != nil {
		return err
	}
	return p.publishRepo(apiUrl, apiUsername, apiPassword, repository, distribution)
}

func (p *packageDeleter) findKeys(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository,
	packageName aptly_package_name.PackageName,
	version aptly_version.Version) ([]aptly_key.Key, error) {
	logger.Debugf("PackageVersions - repo: %s package: %s", repository, packageName)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", apiUrl, repository))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.GetRequest()
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
		return nil, fmt.Errorf("request failed: %s", (content))
	}

	var jsonStruct JsonStruct
	err = json.Unmarshal(content, &jsonStruct)
	if err != nil {
		return nil, err
	}
	var keys []aptly_key.Key
	for _, info := range jsonStruct {
		if info["Package"] == string(packageName) && info["Version"] == string(version) {
			key := info["Key"]
			logger.Debugf("found key: %s", key)
			keys = append(keys, aptly_key.Key(key))
		}
	}
	return keys, nil
}

func (p *packageDeleter) DeletePackagesByKey(
	apiUrl aptly_url.Url,
	apiUsername aptly_user.User,
	apiPassword aptly_password.Password,
	repository aptly_repository.Repository,
	keys []aptly_key.Key) error {
	logger.Debugf("delete package with keys: %v", keys)
	requestbuilder := p.newHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s/packages?format=details", apiUrl, repository))
	requestbuilder.AddBasicAuth(string(apiUsername), string(apiPassword))
	requestbuilder.SetMethod("DELETE")
	requestbuilder.AddContentType("application/json")
	requestContent, err := json.Marshal(map[string][]aptly_key.Key{"PackageRefs": keys})
	if err != nil {
		return err
	}
	requestbuilder.SetBody(bytes.NewBuffer(requestContent))
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request failed: %s", (content))
	}
	return nil
}
