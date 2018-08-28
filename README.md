# httpu [![GoDoc](https://godoc.org/github.com/clavoie/httpu?status.svg)](http://godoc.org/github.com/clavoie/httpu) [![Build Status](https://travis-ci.org/clavoie/httpu.svg?branch=master)](https://travis-ci.org/clavoie/httpu) [![codecov](https://codecov.io/gh/clavoie/httpu/branch/master/graph/badge.svg)](https://codecov.io/gh/clavoie/httpu) [![Go Report Card](https://goreportcard.com/badge/github.com/clavoie/httpu)](https://goreportcard.com/report/github.com/clavoie/httpu)

Http handler utilities for Go.

httpu provides several convenience functions to remove some of the boilerplate from top level http handlers:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
  request := new(MyRequest)
  
  if httpu.DecodeJsonOr400(w, r, request, "Could not decode request") {
    // http.StatusBadRequest has been written to the response, we can now exit the handler
    return
  }
  
  result, err := DoWork(request)
  
  if httpu.Write500IfErr(err, w, "Could not process request for: %v", request.Value) {
    // http.StatusInternalServerError has been written to the response
    return
  }
  
  httpu.EncodeJsonOr400(w, result, "Could not encode response json for: %v", result.Value)
}
```
A full example is available [here](https://godoc.org/github.com/clavoie/httpu#example-WriteIfErr)

## Dependency Injection

httpu provides an interface that wraps all top level functions if you would prefer to inject the package into your project. A full example of using httpu with dependency injection is [here](https://godoc.org/github.com/clavoie/httpu#example-Impl)
