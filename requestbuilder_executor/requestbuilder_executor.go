package requestbuilder_executor

import (
	"fmt"
	"net/http"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type RequestbuilderExecutor interface {
	BuildRequestAndExecute(requestbuilder http_requestbuilder.HttpRequestBuilder) error
}

type requestbuilderExecutor struct {
	executeRequest ExecuteRequest
}

func New(executeRequest ExecuteRequest) *requestbuilderExecutor {
	r := new(requestbuilderExecutor)
	r.executeRequest = executeRequest
	return r
}

func (r *requestbuilderExecutor) BuildRequestAndExecute(requestbuilder http_requestbuilder.HttpRequestBuilder) error {
	req, err := requestbuilder.Build()
	if err != nil {
		return err
	}
	url := req.URL.String()
	method := req.Method
	glog.V(2).Infof("build %s request to %s", method, url)
	resp, err := r.executeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("%s request to %s failed with status %d", method, url, resp.StatusCode)
	}
	return nil
}
