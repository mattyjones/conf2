package data

import (
	"reflect"
	"schema"
)

// Uses reflection to marshal data into go structs
//func MarshalTo(from *Selection, Obj interface{}) error {
//	n := MarshalContainer(Obj)
//	to := NewSelectionFromState(n, from.State)
//	return Upsert(from, to)
//}
//
//func MarshalFrom(Obj interface{}, to *Selection) error {
//	n := MarshalContainer(Obj)
//	from := NewSelectionFromState(n, to.State)
//	return Upsert(from, to)
//}

func MarshalContainer(Obj interface{}) Node {
	s := &MyNode{Label:"Marshal " + reflect.TypeOf(Obj).Name()}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
		objType := reflect.ValueOf(Obj).Elem()
		fieldName := schema.MetaNameToFieldName(meta.GetIdent())
		value := objType.FieldByName(fieldName)
		if schema.IsList(meta) {
			return MarshalList(value.Interface().([]interface{})), nil
		} else {
			if value.Kind() == reflect.Struct {
				return MarshalContainer(value.Addr().Interface()), nil
			} else if value.CanAddr() {
				return MarshalContainer(value.Interface()), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, Obj)
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, Obj, val)
	}
	return s
}

func MarshalList(list []interface{}) Node {
	panic("Not implemented")
	return nil
}
