package httpu

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clavoie/logu"
)

// DecodeJsonOr400 attempts to json decode the request body into the destination object. See
// encoding/json for details.
//
// The request body is closed when this function returns.
//
// If an error is encountered decoding the object a HTTP 400 is written to the response stream,
// and true is returned. If the decoding succeeds then false is returned
func DecodeJsonOr400(w http.ResponseWriter, r *http.Request, dst interface{}, format string, args ...interface{}) bool {
	return NewImpl(w, r, logu.NewGoLogger()).DecodeJsonOr400(dst, format, args...)
}

// EncodeJsonOr500 sets the Content-Type of the response to application/json, and encodes the
// src object into a json response stream. If there is any error encoding the object a
// HTTP 500 is returned instead.
//
// Returns true if there was an error encountered, and false otherwise.
func EncodeJsonOr500(w http.ResponseWriter, src interface{}, format string, args ...interface{}) bool {
	return NewImpl(w, nil, logu.NewGoLogger()).EncodeJsonOr500(src, format, args...)
}

// SetAsDownloadFileWithName sets the Content-Disposition of the response writer to that of
// an attachment with the specified file name.
func SetAsDownloadFileWithName(w http.ResponseWriter, filenameFmt string, args ...interface{}) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", fmt.Sprintf(filenameFmt, args...)))
}

// Write400IfErr works like WriteIfErr(err, http.StatusBadRequest, w, format, args...)
func Write400IfErr(err error, w http.ResponseWriter, format string, args ...interface{}) bool {
	return WriteIfErr(err, http.StatusBadRequest, w, format, args...)
}

// Write500IfErr works like WriteIfErr(err, http.StatusInternalServerError, w, format, args...)
func Write500IfErr(err error, w http.ResponseWriter, format string, args ...interface{}) bool {
	return WriteIfErr(err, http.StatusInternalServerError, w, format, args...)
}

// WriteIfErr writes an http status code and message if there is an error. If the error
// is non-nil then the given http status code is written to the http.ResponseWriter,
// and the message and error are written to the go log package.
//
// true is returned if an error was detected and false is returned if there is no error
func WriteIfErr(err error, statusCode int, w http.ResponseWriter, format string, args ...interface{}) bool {
	return WriteIfErrL(err, statusCode, w, log.Printf, format, args...)
}

// WriteIferrL works like WriteIfErr, writing its output using a specific logging function
func WriteIfErrL(err error, statusCode int, w http.ResponseWriter, logFn func(format string, args ...interface{}), format string, args ...interface{}) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(statusCode)
	args = append(args, err)
	logFn(format+": %v", args...)
	return true
}
