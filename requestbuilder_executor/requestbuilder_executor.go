package requestbuilder_executor
import (
	"fmt"
	"io/ioutil"
	"net/http"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
)

type RequestbuilderExecutor interface {
	BuildRequestAndExecute(requestbuilder http_requestbuilder.HttpRequestBuilder) error
}

type requestbuilderExecutor struct {
	client *http.Client
}

func New(client  *http.Client) *requestbuilderExecutor {
	r := new(requestbuilderExecutor)
	r.client = client
	return r
}

func (r *requestbuilderExecutor) BuildRequestAndExecute(requestbuilder http_requestbuilder.HttpRequestBuilder) error {
	req, err := requestbuilder.GetRequest()
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode / 100 != 2 {
		return fmt.Errorf("upload file failed: %s", string(content))
	}
	return nil
}
