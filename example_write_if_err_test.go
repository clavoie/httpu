package httpu_test

import (
	"net/http"

	"github.com/clavoie/httpu"
)

type JsonRequest struct {
	Value int
}

type JsonResult struct {
	Value int
}

func DoWork(r *JsonRequest) (*JsonResult, error) {
	return new(JsonResult), nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	request := new(JsonRequest)
	if httpu.DecodeJsonOr400(w, r, request, "Could not decode request") {
		return
	}

	result, err := DoWork(request)

	if httpu.Write500IfErr(err, w, "Could not process request for: %v", request.Value) {
		return
	}

	httpu.EncodeJsonOr500(w, result, "Could not encode the response for: %v", result.Value)
}

func ExampleWriteIfErr() {
	http.HandleFunc("/foo", Handler)
}
