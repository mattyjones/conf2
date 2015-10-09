package browse
import (
	"schema"
	"strconv"
	"reflect"
	"errors"
	"fmt"
)


type Value struct {
	Type *schema.DataType
	Bool bool
	Int int
	Str string
	Float float32
	Intlist []int
	Strlist []string
	Boollist []bool
	Keys []string
}

func (v *Value) Value() interface{} {
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		return v.Bool
	case schema.FMT_BOOLEAN_LIST:
		return v.Boollist
	case schema.FMT_INT32, schema.FMT_ENUMERATION:
		return v.Int
	case schema.FMT_INT32_LIST:
		return v.Intlist
	case schema.FMT_STRING:
		return v.Str
	case schema.FMT_STRING_LIST:
		return v.Strlist
	default:
		panic("Not implemented")
	}
}

func (a *Value) Equal(b *Value) bool {
	return a.Value() == b.Value()
}

func (v *Value) SetEnumList(intlist []int) {
	v.Strlist = make([]string, len(intlist))
	for i, n := range intlist {
		v.Strlist[i] = v.Type.Enumeration[n]
	}
}

func (v *Value) SetEnum(n int) {
	v.Int = n
	v.Str = v.Type.Enumeration[n]
}

func (v *Value) String() string {
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		if v.Bool {
			return "true"
		}
		return "false"
	case schema.FMT_INT32:
		return strconv.Itoa(v.Int)
	case schema.FMT_STRING:
		return v.Str
	default:
		panic("Not implemented")
	}
}

// Incoming value should be of appropriate type according to given data type format
func SetValue(typ *schema.DataType, val interface{}) (*Value, error) {
	reflectVal := reflect.ValueOf(val)
	v := &Value{}
	switch typ.Format {
	case schema.FMT_BOOLEAN:
		v.Bool = reflectVal.Bool()
	case schema.FMT_BOOLEAN_LIST:
		v.Boollist = reflectVal.Interface().([]bool)
	case schema.FMT_INT32_LIST:
		v.Intlist = reflectVal.Interface().([]int)
	case schema.FMT_INT32:
		v.Int = int(reflectVal.Int())
	case schema.FMT_STRING, schema.FMT_ENUMERATION:
		v.Str = reflectVal.String()
	case schema.FMT_STRING_LIST:
		v.Strlist = reflectVal.Interface().([]string)
	default:
		return nil, errors.New(fmt.Sprintf("Format code %d not implemented", typ.Format))
	}
	return v, nil
}

func (v *Value) CoerseStrValue(s string) error {
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		v.Bool = s == "true"
	case schema.FMT_INT32:
		var err error
		v.Int, err = strconv.Atoi(s)
		if err != nil {
			return err
		}
	case schema.FMT_STRING:
		v.Str = s
	default:
		return &browseError{Msg:"Coersion not supported from this data format"}
	}
	return nil
}