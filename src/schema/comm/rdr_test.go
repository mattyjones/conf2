package comm

import (
	"testing"
	"schema"
)

func TestCommReaderValue(t *testing.T) {
	data := []byte{10, 0, 0, 0, 99, 0, 0, 0 }
	c := NewReader(data)
	typ := schema.DataType{Format:schema.FMT_INT32}
	if val, err := c.ReadValue(&typ); err != nil {
		t.Error(err)
	} else {
		if val.Int != 99 {
			t.Error("Unexpected comparison", val.Int)
		}
	}
}

