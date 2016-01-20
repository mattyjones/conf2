package data

import (
	"reflect"
	"schema"
)

// Uses reflection to marshal data into go structs
func MarshalContainer(Obj interface{}) Node {
	s := &MyNode{
		Label:"Marshal " + reflect.TypeOf(Obj).Name(),
		Peekables: map[string]interface{}{"internal": Obj},
	}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
		objType := reflect.ValueOf(Obj).Elem()
		fieldName := schema.MetaNameToFieldName(meta.GetIdent())
		value := objType.FieldByName(fieldName)
		if schema.IsList(meta) {
			if value.Kind() == reflect.Map {
				marshal := &MarshalMap{
					Map:value.Interface(),
				}
				return marshal.Node(), nil
			} else {
				marshal := &MarshalArray{
					Array: value.Interface(),
				}
				return marshal.Node(), nil
			}
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

type MarshalArray struct {
	Array        interface{}
	OnNewItem    func() interface{}
	OnSelectItem func(item interface{}) Node
}

func (self *MarshalArray) Node() Node {
	aryReflect := reflect.ValueOf(self.Array)
	n := &MyNode{
		Label: "Marshal " + aryReflect.Type().Name(),
		Peekables: map[string]interface{}{"internal": self.Array},
	}
	n.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		panic("Not implemented")
		return nil, nil
	}
	return n
}

type MarshalMap struct {
	Map          interface{}
	OnNewItem    func() interface{}
	OnSelectItem func(item interface{}) Node
}

func (self *MarshalMap) Node() Node {
	mapReflect := reflect.ValueOf(self.Map)
	n := &MyNode{
		Label: "Marshal " + mapReflect.Type().Name(),
		Peekables: map[string]interface{}{"internal": self.Map},
	}
	index := NewIndex(self.Map)
	n.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		var item interface{}
		if new {
			item = self.OnNewItem()
			mapKey := reflect.ValueOf(key[0].Value())
			mapReflect.SetMapIndex(mapKey, reflect.ValueOf(item))
		} else if len(key) > 0 {
			mapKey := reflect.ValueOf(key[0].Value())
			itemVal := mapReflect.MapIndex(mapKey)
			if itemVal.IsValid() {
				item = itemVal.Interface()
			}
		} else {
			nextKey := index.NextKey(first)
			if nextKey != NO_VALUE {
				sel.path.key = SetValues(meta.KeyMeta(), nextKey.Interface())
				itemVal := mapReflect.MapIndex(nextKey)
				item = itemVal.Interface()
			}
		}
		if item != nil {
			if self.OnSelectItem != nil {
				return self.OnSelectItem(item), nil
			}
			return MarshalContainer(item), nil
		}
		return nil, nil
	}
	return n
}
