package data

import (
	"fmt"
	"reflect"
	"strconv"
	"schema"
	"conf2"
)

type Value struct {
	Type      *schema.DataType
	Bool      bool
	Int       int
	Int64     int64
	Int64list []int64
	Str       string
	Float     float64
	Intlist   []int
	Strlist   []string
	Boollist  []bool
	Keys      []string
	Data      AnyData
}

func (v *Value) Value() interface{} {
	switch v.Type.Format() {
	case schema.FMT_BOOLEAN:
		return v.Bool
	case schema.FMT_BOOLEAN_LIST:
		return v.Boollist
	case schema.FMT_INT64:
		return v.Int64
	case schema.FMT_INT64_LIST:
		return v.Int64list
	case schema.FMT_INT32, schema.FMT_ENUMERATION:
		return v.Int
	case schema.FMT_INT32_LIST, schema.FMT_ENUMERATION_LIST:
		return v.Intlist
	case schema.FMT_STRING:
		return v.Str
	case schema.FMT_STRING_LIST:
		return v.Strlist
	case schema.FMT_ANYDATA:
		return v.Data
	default:
		panic("Not implemented")
	}
}

func (a *Value) Equal(b *Value) bool {
	if a == nil {
		if b == nil {
			return true
		}
		return false
	}
	if b == nil {
		return false
	}
	if a.Type.Format() != b.Type.Format() {
		return false
	}
	if schema.IsListFormat(a.Type.Format()) {
		return reflect.DeepEqual(a.Value(), b.Value())
	}
	return a.Value() == b.Value()
}

func (v *Value) SetEnumList(intlist []int) bool {
	strlist := make([]string, len(intlist))
	en := v.Type.Enumeration()
	for i, n := range intlist {
		if n >= len(en) {
			return false
		}
		strlist[i] = en[n]
	}
	v.Intlist = intlist
	v.Strlist = strlist
	return true
}

func (v *Value) SetEnumListByLabels(labels []string) bool {
	intlist := make([]int, len(labels))
	en := v.Type.Enumeration()
	for i, s := range labels {
		var found bool
		for j, e := range en {
			if s == e {
				found = true
				intlist[i] = j
				break
			}
		}
		if !found {
			return false
		}
	}
	v.Intlist = intlist
	v.Strlist = labels
	return true
}

func (v *Value) SetEnum(n int) bool {
	en := v.Type.Enumeration()
	if n < len(en) {
		v.Int = n
		v.Str = en[n]
		return true
	}
	return false
}

func (v *Value) SetEnumByLabel(label string) bool {
	for i, n := range v.Type.Enumeration() {
		if n == label {
			v.Int = i
			v.Str = label
			return true
		}
	}
	return false
}

func (v *Value) String() string {
	switch v.Type.Format() {
	case schema.FMT_BOOLEAN:
		if v.Bool {
			return "true"
		}
		return "false"
	case schema.FMT_INT32:
		return strconv.Itoa(v.Int)
	case schema.FMT_INT64:
		return strconv.FormatInt(v.Int64, 10)
	case schema.FMT_STRING:
		return v.Str
	default:
		panic("Not implemented")
	}
}

func SetValues(meta []schema.HasDataType, objs ...interface{}) []*Value {
	var err error
	vals := make([]*Value, len(meta))
	for i, obj := range objs {
		vals[i], err = SetValue(meta[i].GetDataType(), obj)
		if err != nil {
			panic(err)
		}
	}
	return vals
}

// Incoming value should be of appropriate type according to given data type format
func SetValue(typ *schema.DataType, val interface{}) (*Value, error) {
	if val == nil {
		return nil, nil
	}
	reflectVal := reflect.ValueOf(val)
	v := &Value{Type: typ}
	switch typ.Format() {
	case schema.FMT_BOOLEAN:
		v.Bool = reflectVal.Bool()
	case schema.FMT_BOOLEAN_LIST:
		v.Boollist = InterfaceToBoollist(reflectVal.Interface())
	case schema.FMT_INT32_LIST:
		v.Intlist = InterfaceToIntlist(val)
	case schema.FMT_INT32:
		switch reflectVal.Kind() {
		case reflect.Float64:
			v.Int = int(reflectVal.Float())
		default:
			v.Int = int(reflectVal.Int())
		}
	case schema.FMT_DECIMAL64:
		v.Float = reflectVal.Float()
	case schema.FMT_INT64:
		v.Int64 = reflectVal.Int()
	case schema.FMT_INT64_LIST:
		v.Int64list = InterfaceToInt64list(val)
	case schema.FMT_STRING:
		v.Str = reflectVal.String()
	case schema.FMT_ENUMERATION:
		switch reflectVal.Kind() {
		case reflect.String:
			v.SetEnumByLabel(reflectVal.String())
		default:
			v.SetEnum(int(reflectVal.Int()))
		}
	case schema.FMT_ENUMERATION_LIST:
		val := reflectVal.Interface()
		strlist := InterfaceToStrlist(val)
		if len(strlist) > 0 {
			v.SetEnumListByLabels(strlist)
		} else {
			intlist := InterfaceToIntlist(val)
			v.SetEnumList(intlist)
		}
	case schema.FMT_STRING_LIST:
		v.Strlist = InterfaceToStrlist(reflectVal.Interface())
	case schema.FMT_ANYDATA:
		v.Data = reflectVal.Interface().(AnyData)
	default:
		panic(fmt.Sprintf("Format code %d not implemented", typ.Format))
	}
	return v, nil
}

func InterfaceToStrlist(o interface{}) (strlist []string) {
	switch arrayValues := o.(type) {
	case []string:
		return arrayValues
	case []interface{}:
		found := make([]string, len(arrayValues))
		for i, arrayValue := range arrayValues {
			if reflect.TypeOf(arrayValue).Kind() != reflect.String {
				return
			}
			found[i] = arrayValue.(string)
		}
		strlist = found
	}
	return
}

func InterfaceToBoollist(o interface{}) (boollist []bool) {
	switch arrayValues := o.(type) {
	case []bool:
		return arrayValues
	case []interface{}:
		boollist = make([]bool, len(arrayValues))
		for i, arrayValue := range arrayValues {
			boollist[i] = arrayValue.(bool)
		}
	}
	return
}

func InterfaceToInt64list(o interface{}) (intlist []int64) {
	switch arrayValues := o.(type) {
	case []int64:
		return arrayValues
	case []interface{}:
		intlist = make([]int64, len(arrayValues))
		for i, arrayValue := range arrayValues {
			intlist[i] = arrayValue.(int64)
		}
	}
	return
}

func InterfaceToIntlist(o interface{}) (intlist []int) {
	switch arrayValues := o.(type) {
	case []int:
		return arrayValues
	case []interface{}:
		intlist = make([]int, len(arrayValues))
		for i, arrayValue := range arrayValues {
			intlist[i] = arrayValue.(int)
		}
	}
	return
}

func (v *Value) CoerseStrValue(s string) error {
	switch v.Type.Format() {
	case schema.FMT_BOOLEAN:
		v.Bool = s == "true"
	case schema.FMT_INT64:
		var err error
		v.Int64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
	case schema.FMT_INT32:
		var err error
		v.Int, err = strconv.Atoi(s)
		if err != nil {
			return err
		}
	case schema.FMT_STRING:
		v.Str = s
	case schema.FMT_ENUMERATION:
		eid, err := strconv.Atoi(s)
		if err == nil {
			v.SetEnum(eid)
			return nil
		}
		if ! v.SetEnumByLabel(s) {
			return conf2.NewErr("Not an allowed enumation: " + s)
		}
	default:
		panic(fmt.Sprintf("Coersion not supported from data format " + v.Type.Format().String()))
	}
	return nil
}
