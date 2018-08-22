package handlerutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func TryDecodeJsonFile(w http.ResponseWriter, r *http.Request, filename string, dst interface{}) bool {
	err := r.ParseMultipartForm(10000000)
	// https://golang.org/pkg/mime/multipart/#Form

	if Write500IfErr(w, r, err, fmt.Sprintf("Could not parse multipart form datafor file %v: %%v", filename)) {
		return true
	}

	file, _, err := r.FormFile(filename)

	if Write400IfErr(w, r, err, fmt.Sprintf("Could not find file with name %v in upload: %%v", filename)) {
		return true
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(dst)

	return Write400IfErr(w, r, err, fmt.Sprintf("Could not decode json from file %v: %%v", filename))
}
