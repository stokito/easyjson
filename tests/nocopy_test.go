package tests

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/mailru/easyjson"
)

// verifies if string pointer belongs to the given buffer or outside of it
func strBelongsTo(s string, buf []byte) bool {
	sPtr := (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	bufPtr := (*reflect.SliceHeader)(unsafe.Pointer(&buf)).Data

	if bufPtr <= sPtr && sPtr < bufPtr+uintptr(len(buf)) {
		return true
	}
	return false
}

func TestNocopy(t *testing.T) {
	data := []byte(`{"a": "valueA", "b": "valueB"}`)
	exp := NocopyStruct{
		A: "valueA",
		B: "valueB",
	}
	res := NocopyStruct{}

	easyjson.Unmarshal(data, &res)
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("TestNocopy(): got=%+v, exp=%+v", res, exp)
	}

	if strBelongsTo(res.A, data) {
		t.Error("TestNocopy(): field A was not copied and refers to buffer")
	}
	if !strBelongsTo(res.B, data) {
		t.Error("TestNocopy(): field B was copied rather than refer to bufferr")
	}

	data = []byte(`{"b": "valueNoCopy"}`)
	res = NocopyStruct{}
	allocsPerRun := testing.AllocsPerRun(1000, func() {
		easyjson.Unmarshal(data, &res)
		if res.B != "valueNoCopy" {
			t.Fatalf("wrong value: %q", res.B)
		}
	})
	if allocsPerRun != 1 {
		t.Fatalf("expected 1 allocs, got %f", allocsPerRun)
	}
}
