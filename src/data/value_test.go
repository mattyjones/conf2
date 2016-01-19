package data

import (
	"testing"
	"schema"
)

func TestCoerseValue(t *testing.T) {
	v, err := SetValue(schema.NewDataType(nil, "int32"), 35)
	if err != nil {
		t.Error(err)
	} else if v.Int != 35 {
		t.Error("Coersion error")
	}
}
