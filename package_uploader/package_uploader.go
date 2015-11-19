package package_uploader

import (
	"net/http"
)

type PackageUploader interface {
	UploadPackage(api_url string, api_username string, api_password string, file string) error
}

type packageUploader struct {
	httpClientProvider HttpClientProvider
}

type HttpClientProvider func() *http.Client

func New(httpClientProvider HttpClientProvider) *packageUploader {
	p := new(packageUploader)
	p.httpClientProvider = httpClientProvider
	return p
}

func (p *packageUploader ) UploadPackage(api_url string, api_username string, api_password string, file string) error {
	req, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		return err
	}
	if len(api_username) > 0 || len(api_password) > 0 {
		req.SetBasicAuth(api_username, api_password)
	}
	client := p.httpClientProvider()
	resp, err := client.Do(req)
	if err != nil {
		return  err
	}
	defer resp.Body.Close()
	return nil
}
