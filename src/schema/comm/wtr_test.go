package comm

import (
	"bytes"
	"schema"
	"schema/browse"
	"testing"
)

func TestCommWriterValue(t *testing.T) {
	c := NewWriter()
	typ := schema.DataType{Format: schema.FMT_INT32}
	val := &browse.Value{Type: &typ, Int: 99}
	c.WriteValue(val)
	actual := c.Data()
	expected := []byte{10, 0, 0, 0, 99, 0, 0, 0}
	if bytes.Compare(expected, actual) != 0 {
		t.Error("Unexpected comparison", actual)
	}
}
