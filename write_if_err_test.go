package httpu_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clavoie/httpu"
)

func TestWriteIfErr(t *testing.T) {
	type Json struct {
		Field int
	}

	fieldValue := 100
	successfulJson := fmt.Sprintf(`{"Field":%v}`, fieldValue)
	format := "hello: %v"
	formatArgs := []interface{}{1}
	expectedFormat := fmt.Sprintf(format, formatArgs...)

	newImpl := func(data string, t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
		r := httptest.NewRequest("POST", "http://test.com", bytes.NewBufferString(data))
		r.Header.Set("Content-Type", `multipart/form-data; boundary="xxx"`)
		w := httptest.NewRecorder()

		return r, w
	}

	t.Run("DecodeJsonOr400", func(t *testing.T) {
		r, w := newImpl(successfulJson, t)

		j := new(Json)
		if httpu.DecodeJsonOr400(w, r, j, format, formatArgs...) {
			t.Fatal("Was not expecting an error")
		}

		if w.Code == http.StatusBadRequest {
			t.Fatalf("Unexpected code: %v", w.Code)
		}
	})

	t.Run("EncodeJsonOr500", func(t *testing.T) {
		_, w := newImpl(``, t)

		j := &Json{fieldValue}
		if httpu.EncodeJsonOr500(w, j, format, formatArgs...) {
			t.Fatal("Was not expecting an error")
		}

		actualJ := strings.TrimSpace(w.Body.String())
		if actualJ != successfulJson {
			t.Fatal("Was expecting: ", successfulJson, actualJ)
		}
	})

	t.Run("SetAsDownloadFileWithName", func(t *testing.T) {
		_, w := newImpl(``, t)

		httpu.SetAsDownloadFileWithName(w, format, formatArgs...)

		content := w.HeaderMap.Get("Content-Disposition")
		if content == "" {
			t.Fatal("Was expecting non empty content disposition")
		}

		expectedValue := fmt.Sprintf("attachment; filename=\"%v\"", expectedFormat)
		if content != expectedValue {
			t.Fatal("Content value did not match", content, expectedValue)
		}
	})

	t.Run("Write400IfErr", func(t *testing.T) {
		_, w := newImpl(``, t)

		if httpu.Write400IfErr(nil, w, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})

	t.Run("Write500IfErr", func(t *testing.T) {
		_, w := newImpl(``, t)

		if httpu.Write500IfErr(nil, w, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})

	t.Run("WriteIfErrNone", func(t *testing.T) {
		_, w := newImpl(``, t)

		if httpu.WriteIfErr(nil, http.StatusOK, w, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})

	//
	// TryDeocdeJsonFile
	//

	filename := "filename"

	t.Run("TryDecodeJsonFileParseFail", func(t *testing.T) {
		r, w := newImpl(``, t)

		j := new(Json)

		if httpu.TryDecodeJsonFile(w, r, filename, j) == false {
			t.Fatal("Was expecting parse failure")
		}

		if w.Code != 500 {
			t.Fatal("Was expecting 500 status code")
		}
	})
}
