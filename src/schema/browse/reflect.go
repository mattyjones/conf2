package browse
import (
	"schema"
	"reflect"
	"fmt"
	"errors"
)

func ReadField(meta schema.HasDataType, obj interface{}) (*Value, error) {
	return ReadFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj)
}

func ReadFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}) (v *Value, err error) {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	v = &Value{}
	switch meta.GetDataType().Format {
	case schema.FMT_BOOLEAN:
		v.Bool = value.Bool()
	case schema.FMT_BOOLEAN_LIST:
		v.Boollist = value.Interface().([]bool)
	case schema.FMT_INT32_LIST:
		v.Intlist = value.Interface().([]int)
	case schema.FMT_INT32:
		v.Int = int(value.Int())
	case schema.FMT_STRING:
		v.Str = value.String()
	case schema.FMT_STRING_LIST:
		v.Strlist = value.Interface().([]string)
	default:
		err = errors.New(fmt.Sprintf("Format code %d not implemented", meta.GetDataType().Format))
	}
	return
}

func WriteField(meta schema.HasDataType, obj interface{}, v *Value) error {
	return WriteFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func WriteFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	if ! objType.IsValid() {
		return &browseError{Msg:fmt.Sprintf("Cannot find property \"%s\" on invalid or nil %s", fieldName, reflect.TypeOf(obj))}
	}
	value := objType.FieldByName(fieldName)
	switch v.Type.Format {
	case schema.FMT_BOOLEAN_LIST:
		value.Set(reflect.ValueOf(v.Boollist))
	case schema.FMT_BOOLEAN:
		value.SetBool(v.Bool)
	case schema.FMT_INT32_LIST:
		value.Set(reflect.ValueOf(v.Intlist))
	case schema.FMT_INT32:
		value.SetInt(int64(v.Int))
	case schema.FMT_STRING_LIST:
		value.Set(reflect.ValueOf(v.Strlist))
	case schema.FMT_STRING:
		value.SetString(v.Str)
	default:
		return NotImplemented(meta)
	}
	return nil
}