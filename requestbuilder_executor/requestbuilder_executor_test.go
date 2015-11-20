package requestbuilder_executor


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsRequestbuilderExecutor(t *testing.T) {
	b := New(nil)
	var i *RequestbuilderExecutor
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}