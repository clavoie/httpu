package httpu_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/clavoie/logu/mock_logu"
	"github.com/golang/mock/gomock"
)

func TestImpl(t *testing.T) {
	type Json struct {
		Field int
	}

	fieldValue := 100
	successfulJson := fmt.Sprintf(`{"Field":%v}`, fieldValue)

	t.Run("DecodeJsonOr400", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://test.com", bytes.NewBufferString(successfulJson))
		w := httptest.NewResponseRecorder()
		ctrl := gomock.NewController(t)
		l := mock_logu.NewMockLogger(ctrl)
	})
}
