package schema

import (
	"testing"
)

func TestCoerseValue(t *testing.T) {
	v, err := SetValue(&DataType{Format: FMT_INT32}, 35)
	if err != nil {
		t.Error(err)
	} else if v.Int != 35 {
		t.Error("Coersion error")
	}
}
