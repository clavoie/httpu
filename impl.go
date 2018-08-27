package httpu

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clavoie/logu"
)

// Impl is a wrapper around all top level package functions
type Impl interface {
	// DecodeJsonOr400 attempts to json decode the request body into the destination object. See
	// encoding/json for details.
	//
	// The request body is closed when this function returns.
	//
	// If an error is encountered decoding the object a HTTP 400 is written to the response stream,
	// and true is returned. If the decoding succeeds then false is returned
	DecodeJsonOr400(dst interface{}, format string, args ...interface{}) bool

	// EncodeJsonOr500 sets the Content-Type of the response to application/json, and encodes the
	// src object into a json response stream. If there is any error encoding the object a
	// HTTP 500 is returned instead.
	//
	// Returns true if there was an error encountered, and false otherwise.
	EncodeJsonOr500(src interface{}, format string, args ...interface{}) bool

	// SetAsDownloadFileWithName sets the Content-Disposition of the response writer to that of
	// an attachment with the specified file name.
	SetAsDownloadFileWithName(filenameFmt string, args ...interface{})

	// TryDecodeJsonFile attempts to parse a file upload from the request, and json deserialize
	// its contents into a destination object. If the multipart form cannot be parsed from the
	// request a HTTP 500 is written to the response and true is returned. If the file with the
	// given filename cannot be found, a HTTP 400 is written to the response and true is
	// returned. If the body of the file cannot be successfully json decoded into the destination
	// object, a HTTP 400 is written to the response and true is returned.
	//
	// If the entire operation is a success false is returned.
	TryDecodeJsonFile(filename string, dst interface{}) bool

	// Write400IfErr works like WriteIfErr(err, http.StatusBadRequest, format, args...)
	Write400IfErr(err error, format string, args ...interface{}) bool

	// Write500IfErr works like WriteIfErr(err, http.StatusInternalServerError, format, args...)
	Write500IfErr(err error, format string, args ...interface{}) bool

	// WriteIfErr writes an http status code and message if there is an error. If the error
	// is non-nil then the given http status code is written to the http.ResponseWriter,
	// and the message and error are written to the go log package.
	//
	// true is returned if an error was detected and false is returned if there is no error
	WriteIfErr(err error, statusCode int, format string, args ...interface{}) bool
}

// impl is an implementation of Impl
type impl struct {
	l logu.Logger
	r *http.Request
	w http.ResponseWriter
}

// NewImpl returns a new instance of Impl
func NewImpl(w http.ResponseWriter, r *http.Request, l logu.Logger) Impl {
	return &impl{
		l: l,
		r: r,
		w: w,
	}
}

func (i *impl) DecodeJsonOr400(dst interface{}, format string, args ...interface{}) bool {
	defer i.r.Body.Close()

	decoder := json.NewDecoder(i.r.Body)
	err := decoder.Decode(dst)

	return i.Write400IfErr(err, format, args...)
}

func (i *impl) EncodeJsonOr500(src interface{}, format string, args ...interface{}) bool {
	i.w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(i.w)
	err := encoder.Encode(src)

	return i.Write500IfErr(err, format, args...)
}

func (i *impl) TryDecodeJsonFile(filename string, dst interface{}) bool {
	err := i.r.ParseMultipartForm(10000000)

	if i.Write500IfErr(err, "Could not parse multipart form for file %v", filename) {
		return true
	}

	defer func() {
		err := i.r.MultipartForm.RemoveAll()

		if err != nil {
			i.l.Errorf("%v", err)
		}
	}()
	file, _, err := i.r.FormFile(filename)

	if i.Write400IfErr(err, "Could not find file with name %v in upload", filename) {
		return true
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(dst)

	return i.Write400IfErr(err, "Could not decode json from file %v", filename)
}

func (i *impl) SetAsDownloadFileWithName(filenameFmt string, args ...interface{}) {
	i.w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", fmt.Sprintf(filenameFmt, args...)))
}

func (i *impl) Write400IfErr(err error, format string, args ...interface{}) bool {
	return i.WriteIfErr(err, http.StatusBadRequest, format, args...)
}

func (i *impl) Write500IfErr(err error, format string, args ...interface{}) bool {
	return i.WriteIfErr(err, http.StatusInternalServerError, format, args...)
}

func (i *impl) WriteIfErr(err error, statusCode int, format string, args ...interface{}) bool {
	logFn := i.l.Warningf

	if statusCode >= 500 {
		logFn = i.l.Errorf
	}

	if err == nil {
		return false
	}

	i.w.WriteHeader(statusCode)
	args = append(args, err)
	logFn(format+": %v", args...)
	return true
}
