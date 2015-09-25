package browse
import (
	"schema"
	"strconv"
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
	case schema.FMT_INT32:
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