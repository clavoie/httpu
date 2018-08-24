package httpu_test

import (
	"testing"
)

func TestNewDiDefs(t *testing.T) {
	defs := httpu.NewDiDefs()

	if defs == nil {
		t.Fatal("Expecting non-nil defs")
	}

	defs2 := httpu.NewDiDefs()
	if defs[0] == defs2[0] {
		t.Fatal("Not expecting defs to match")
	}
}
