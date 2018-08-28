package httpu_test

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/clavoie/di"
	"github.com/clavoie/erru"
	"github.com/clavoie/httpu"
	"github.com/clavoie/logu"
)

func onResolveErr(err *di.ErrResolve, w http.ResponseWriter, r *http.Request) {
	logger := logu.NewAppEngineLogger(r)
	logger.Errorf("err encountered while resolving dependencies: %v", err.String())

	httpErr, isHttpErr := err.Err.(erru.HttpErr)
	if isHttpErr {
		w.WriteHeader(httpErr.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type Request struct {
	Value int
}

type Result struct {
	Value int
}

type Dependency interface {
	DoWork(*Request) (*Result, error)
}

type dependency struct{}

func (d *dependency) DoWork(r *Request) (*Result, error) { return new(Result), nil }
func NewDependency() Dependency                          { return new(dependency) }

var defs = []*di.Def{{NewDependency, di.PerHttpRequest}}

func MyHandler(dep Dependency, helper httpu.Impl) {
	request := new(Request)
	if helper.DecodeJsonOr400(request, "Could not decode request") {
		return
	}

	result, err := dep.DoWork(request)

	if helper.Write500IfErr(err, "Could not process request for: %v", request.Value) {
		return
	}

	helper.EncodeJsonOr500(result, "Could not encode the response for: %v", result.Value)
}

func ExampleImpl() {
	resolver, err := di.NewResolver(onResolveErr, defs, logu.NewAppEngineDiDefs(), httpu.NewDiDefs())

	if err != nil {
		log.Fatal(err)
	}

	httpHandler, err := resolver.HttpHandler(MyHandler)

	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	httpHandler(w, req)
}
