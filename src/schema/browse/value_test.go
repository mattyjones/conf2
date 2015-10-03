package browse
import (
	"testing"
	"schema"
)

func TestCoerseValue(t *testing.T) {
	v, err := SetValue(&schema.DataType{Format:schema.FMT_INT32}, 35)
	if err != nil {
		t.Error(err)
	} else if v.Int != 35 {
		t.Error("Coersion error")
	}
}