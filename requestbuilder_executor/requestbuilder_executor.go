package requestbuilder_executor

import (
	"fmt"
	"io/ioutil"
	"net/http"

	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)


type ExecuteRequest func (req *http.Request) (resp *http.Response, err error)

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
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := r.executeRequest(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request file failed: %s", string(content))
	}
	return nil
}
