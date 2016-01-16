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
	s := &MyNode{
		Label:"Marshal " + reflect.TypeOf(Obj).Name(),
		Internal: Obj,
	}
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
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, Obj)
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *Value) error {
		return WriteField(meta, Obj, val)
	}
	return s
}

func MarshalList(list []interface{}) Node {
	panic("Not implemented")
	return nil
}

type MarshalIndexedList struct {
	Map interface{}
	OnNewItem func() interface{}
	OnSelectItem func(item interface{}) Node
}

func (self *MarshalIndexedList) Node() Node {
	mapReflect := reflect.ValueOf(self.Map)
	n := &MyNode{
		Label: "Marshal " + mapReflect.Type().Name(),
		Internal: self.Map,
	}
	//valueReflect := reflect.ValueOf(prototype)
	index := NewIndex(self.Map)
	n.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		var item interface{}
		if new {
			item = self.OnNewItem()
		} else if len(key) > 0 {
			mapKey := reflect.ValueOf(key[0].Value())
			item = mapReflect.MapIndex(mapKey).Interface()
		} else {
			nextKey := index.NextKey(first)
			sel.State.SetKey(SetValues(meta.KeyMeta(), nextKey))
			item = mapReflect.MapIndex(reflect.ValueOf(nextKey)).Interface()
		}
		if item != nil {
			return self.OnSelectItem(item), nil
		}
		return nil, nil
	}
	return n
}
