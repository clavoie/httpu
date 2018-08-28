package httpu_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clavoie/httpu"
	"github.com/clavoie/logu/mock_logu"
	"github.com/golang/mock/gomock"
)

func TestImpl(t *testing.T) {
	type Json struct {
		Field int
	}

	fieldValue := 100
	successfulJson := fmt.Sprintf(`{"Field":%v}`, fieldValue)
	format := "hello: %v"
	errFormat := format + ": %v"
	formatArgs := []interface{}{1}
	expectedFormat := fmt.Sprintf(format, formatArgs...)
	err := errors.New("my error")

	newImpl := func(data string, t *testing.T) (*http.Request, *httptest.ResponseRecorder, *mock_logu.MockLogger, httpu.Impl, func()) {
		r := httptest.NewRequest("POST", "http://test.com", bytes.NewBufferString(data))
		r.Header.Set("Content-Type", `multipart/form-data; boundary="xxx"`)
		w := httptest.NewRecorder()
		ctrl := gomock.NewController(t)
		l := mock_logu.NewMockLogger(ctrl)
		i := httpu.NewImpl(w, r, l)

		return r, w, l, i, ctrl.Finish
	}
	newErrArgs := func(args ...interface{}) []interface{} {
		errArgs := append(make([]interface{}, 0, 2), formatArgs...)
		return append(errArgs, args...)
	}

	errArgs := newErrArgs(err)

	t.Run("DecodeJsonOr400", func(t *testing.T) {
		_, w, _, i, finish := newImpl(successfulJson, t)
		defer finish()

		j := new(Json)
		if i.DecodeJsonOr400(j, format, formatArgs...) {
			t.Fatal("Was not expecting an error")
		}

		if w.Code == http.StatusBadRequest {
			t.Fatalf("Unexpected code: %v", w.Code)
		}
	})
	t.Run("DecodeJsonOr400Fail", func(t *testing.T) {
		_, w, l, i, finish := newImpl(`{"Invalid":x}`, t)
		defer finish()

		errArgs := newErrArgs(NonEmptyStr())
		l.EXPECT().Warningf(errFormat, errArgs...)

		j := new(Json)
		if i.DecodeJsonOr400(j, format, formatArgs...) == false {
			t.Fatal("Expecting to fail on decode")
		}

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Unexpected code: %v", w.Code)
		}
	})

	t.Run("EncodeJsonOr500", func(t *testing.T) {
		_, w, _, i, finish := newImpl(``, t)
		defer finish()

		j := &Json{fieldValue}
		if i.EncodeJsonOr500(j, format, formatArgs...) {
			t.Fatal("Was not expecting an error")
		}

		actualJ := strings.TrimSpace(w.Body.String())
		if actualJ != successfulJson {
			t.Fatal("Was expecting: ", successfulJson, actualJ)
		}
	})
	t.Run("EncodeJsonOr500Fail", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		type Err struct {
			E ErrType
		}
		e := new(Err)

		errArgs := newErrArgs(NonEmptyStr())
		l.EXPECT().Errorf(errFormat, errArgs...)

		if i.EncodeJsonOr500(e, format, formatArgs...) == false {
			t.Fatal("Was expecting an error")
		}

		if w.Code != http.StatusInternalServerError {
			t.Fatal("Was expecting internal server error", w.Code)
		}
	})

	t.Run("SetAsDownloadFileWithName", func(t *testing.T) {
		_, w, _, i, finish := newImpl(``, t)
		defer finish()

		i.SetAsDownloadFileWithName(format, formatArgs...)

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
		_, _, _, i, finish := newImpl(``, t)
		defer finish()

		if i.Write400IfErr(nil, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})
	t.Run("Write400IfErrFail", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		l.EXPECT().Warningf(errFormat, errArgs...)
		code := http.StatusBadRequest

		if i.Write400IfErr(err, format, formatArgs...) == false {
			t.Fatal("Was expecting err")
		}

		if w.Code != code {
			t.Fatalf("expecting code: %v, found %v", code, w.Code)
		}
	})

	t.Run("Write500IfErr", func(t *testing.T) {
		_, _, _, i, finish := newImpl(``, t)
		defer finish()

		if i.Write500IfErr(nil, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})
	t.Run("Write500IfErrFail", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		l.EXPECT().Errorf(errFormat, errArgs...)
		code := http.StatusInternalServerError

		if i.Write500IfErr(err, format, formatArgs...) == false {
			t.Fatal("Was expecting err")
		}

		if w.Code != code {
			t.Fatalf("expecting code: %v, found %v", code, w.Code)
		}
	})

	t.Run("WriteIfErrNone", func(t *testing.T) {
		_, _, _, i, finish := newImpl(``, t)
		defer finish()

		if i.WriteIfErr(nil, http.StatusOK, format, formatArgs...) {
			t.Fatal("Was not expecting err")
		}
	})
	t.Run("WriteIfErrWarn", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		l.EXPECT().Warningf(errFormat, errArgs...)
		code := 400

		if i.WriteIfErr(err, code, format, formatArgs...) == false {
			t.Fatal("Was expecting err")
		}

		if w.Code != code {
			t.Fatalf("expecting code: %v, found %v", code, w.Code)
		}
	})
	t.Run("WriteIfErrError", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		l.EXPECT().Errorf(errFormat, errArgs...)
		code := 500

		if i.WriteIfErr(err, code, format, formatArgs...) == false {
			t.Fatal("Was expecting err")
		}

		if w.Code != code {
			t.Fatalf("expecting code: %v, found %v", code, w.Code)
		}
	})

	//
	// TryDeocdeJsonFile
	//

	filename := "filename"
	postData := fmt.Sprintf(`--xxx
Content-Disposition: form-data; name="%v"; filename="%v"
Content-Type: application/octet-stream
Content-Transfer-Encoding: binary

%v
--xxx--
`, filename, filename, successfulJson)

	t.Run("TryDecodeJsonFileParseFail", func(t *testing.T) {
		_, w, l, i, finish := newImpl(``, t)
		defer finish()

		l.EXPECT().Errorf(NonEmptyStr(), filename)
		j := new(Json)

		if i.TryDecodeJsonFile(filename, j) == false {
			t.Fatal("Was expecting parse failure")
		}

		if w.Code != 500 {
			t.Fatal("Was expecting 500 status code")
		}
	})
	t.Run("TryDecodeJsonFileFileNotFound", func(t *testing.T) {
		_, w, l, i, finish := newImpl(postData, t)
		defer finish()

		badFilename := filename + "x"
		l.EXPECT().Warningf(NonEmptyStr(), badFilename)
		j := new(Json)

		if i.TryDecodeJsonFile(badFilename, j) == false {
			t.Fatal("Was expecting parse failure")
		}

		if w.Code != 400 {
			t.Fatal("Was expecting 400 status code")
		}
	})
	t.Run("TryDecodeJsonFileSuccess", func(t *testing.T) {
		_, _, _, i, finish := newImpl(postData, t)
		defer finish()

		j := new(Json)
		if i.TryDecodeJsonFile(filename, j) {
			t.Fatal("Was expecting success")
		}

		if j.Field != fieldValue {
			t.Fatal(j.Field, fieldValue)
		}
	})
}
