package httpu_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

type nonEmptyStrMatcher struct {
}

func NonEmptyStr() gomock.Matcher {
	return &nonEmptyStrMatcher{}
}

func (m *nonEmptyStrMatcher) Matches(x interface{}) bool {
	var str string

	if stringer, isStringer := x.(fmt.Stringer); isStringer {
		str = stringer.String()
	} else if err, isErr := x.(error); isErr {
		str = err.Error()
	} else if s, isStr := x.(string); isStr {
		str = s
	}

	return str != ""
}

func (m *nonEmptyStrMatcher) String() string {
	return "Non empty string"
}
