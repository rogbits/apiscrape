package tools

import (
	"testing"
)

func TestIdGen(t *testing.T) {
	id := GenerateId(8)
	if len(id) != 8 {
		t.FailNow()
	}
}
