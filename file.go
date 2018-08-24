package httpu

import (
	"net/http"

	"github.com/clavoie/logu"
)

// TryDecodeJsonFile attempts to parse a file upload from the request, and json deserialize
// its contents into a destination object. If the multipart form cannot be parsed from the
// request a HTTP 500 is written to the response and true is returned. If the file with the
// given filename cannot be found, a HTTP 400 is written to the response and true is
// returned. If the body of the file cannot be successfully json decoded into the destination
// object, a HTTP 400 is written to the response and true is returned.
//
// If the entire operation is a success false is returned.
func TryDecodeJsonFile(w http.ResponseWriter, r *http.Request, filename string, dst interface{}) bool {
	return NewImpl(w, r, logu.NewGoLogger()).TryDecodeJsonFile(filename, dst)
}
