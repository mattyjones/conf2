package data

import (
	"fmt"
	"reflect"
	"schema"
	"conf2"
)

func ReadField(meta schema.HasDataType, obj interface{}) (*Value, error) {
	return ReadFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj)
}

func ReadFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}) (v *Value, err error) {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	v = &Value{Type: meta.GetDataType()}
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		v.Bool = value.Bool()
	case schema.FMT_BOOLEAN_LIST:
		v.Boollist = value.Interface().([]bool)
	case schema.FMT_INT32_LIST:
		v.Intlist = value.Interface().([]int)
	case schema.FMT_INT64_LIST:
		v.Int64list = value.Interface().([]int64)
	case schema.FMT_INT32:
		v.Int = int(value.Int())
	case schema.FMT_INT64:
		v.Int64 = value.Int()
	case schema.FMT_STRING:
		v.Str = value.String()
	case schema.FMT_STRING_LIST:
		v.Strlist = value.Interface().([]string)
	case schema.FMT_ENUMERATION:
		switch value.Type().Kind() {
		case reflect.String:
			v.SetEnumByLabel(value.String())
		default:
			v.SetEnum(int(value.Int()))
		}
	case schema.FMT_ANYDATA:
		if anyData, isAnyData := value.Interface().(AnyData); isAnyData {
			v.Data = anyData
		} else {
			return nil, conf2.NewErr("Cannot read anydata from value that doesn't implement AnyData")
		}
	default:
		panic(fmt.Sprintf("Format code %d not implemented", meta.GetDataType().Format))
	}
	return
}

func WriteField(meta schema.HasDataType, obj interface{}, v *Value) error {
	return WriteFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func WriteFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	if !objType.IsValid() {
		panic(fmt.Sprintf("Cannot find property \"%s\" on invalid or nil %s", fieldName, reflect.TypeOf(obj)))
	}
	value := objType.FieldByName(fieldName)
	if !value.IsValid() {
		panic(fmt.Sprintf("Invalid property \"%s\" on %s", fieldName, reflect.TypeOf(obj)))
	}
	switch v.Type.Format {
	case schema.FMT_BOOLEAN_LIST:
		value.Set(reflect.ValueOf(v.Boollist))
	case schema.FMT_BOOLEAN:
		value.SetBool(v.Bool)
	case schema.FMT_INT32_LIST:
		value.Set(reflect.ValueOf(v.Intlist))
	case schema.FMT_INT32:
		value.SetInt(int64(v.Int))
	case schema.FMT_INT64_LIST:
		value.Set(reflect.ValueOf(v.Int64list))
	case schema.FMT_INT64:
		value.SetInt(v.Int64)
	case schema.FMT_STRING_LIST:
		value.Set(reflect.ValueOf(v.Strlist))
	case schema.FMT_STRING:
		value.SetString(v.Str)
	case schema.FMT_ENUMERATION:
		switch value.Type().Kind() {
		case reflect.String:
			value.SetString(v.Str)
		default:
			value.SetInt(int64(v.Int))
		}
	case schema.FMT_ANYDATA:
		// could support writing to string as well
		value.Set(reflect.ValueOf(v.Data))

	// TODO: Enum list
	default:
		panic(meta.GetIdent() + " not implemented")
	}
	return nil
}
