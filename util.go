package handlerutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON encode / decode here is OK, but anything dal related needs
// to be removed

// DecodeJsonOr400 returns true if there is an error
func DecodeJsonOr400(w http.ResponseWriter, r *http.Request, dst interface{}, format string) bool {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)

	return Write400IfErr(w, r, err, format)
}

func EncodeJsonOr500(w http.ResponseWriter, r *http.Request, src interface{}, format string) bool {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(src)

	return Write500IfErr(w, r, err, format)
}

func SetAsDownloadFileWithName(w http.ResponseWriter, filenameFmt string, args ...interface{}) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", fmt.Sprintf(filenameFmt, args...)))
}

// Write400IfErr works like WriteIfErr(err, http.StatusBadRequest, w, message)
func Write400IfErr(err error, w http.ResponseWriter, message string) bool {
	return WriteIfErr(err, http.StatusBadRequest, w, message)
}

// Write500IfErr works like WriteIfErr(err, http.StatusInternalServerError, w, message)
func Write500IfErr(err error, w http.ResponseWriter, message string) bool {
	return WriteIfErr(err, http.StatusInternalServerError, w, message)
}

// WriteIfErr writes an http status code and message if there is an error. If the error
// is non-nil then the given http status code is written to the http.ResponseWriter,
// and the message and error are written to the go log package.
//
// true is returned if an error was detected and false is returned if there is no error
func WriteIfErr(err error, statusCode int, w http.ResponseWriter, message string) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(statusCode)
	log.Println("%v: %v", message, err)
	return true
}
