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
	s.OnSelect = func(sel *Selection, r ContainerRequest) (Node, error) {
		objType := reflect.ValueOf(Obj).Elem()
		fieldName := schema.MetaNameToFieldName(r.Meta.GetIdent())
		value := objType.FieldByName(fieldName)
		if schema.IsList(r.Meta) {
			if value.Kind() == reflect.Map {
				marshal := &MarshalMap{
					Map:value.Interface(),
				}
				return marshal.Node(), nil
			} else {
				marshal := &MarshalArray{
					ArrayValue: &value,
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
	ArrayValue   *reflect.Value
	OnNewItem    func() interface{}
	OnSelectItem func(item interface{}) Node
}

func (self *MarshalArray) Node() Node {
	//aryReflect := reflect.ValueOf(self.Array)
	var i int
	n := &MyNode{
		Label: "Marshal " + self.ArrayValue.Type().Name(),
		Peekables: map[string]interface{}{"internal": self.ArrayValue.Interface()},
	}
	n.OnNext = func(sel *Selection, r ListRequest) (next Node, key []*Value, err error) {
		var item interface{}
		if r.New {
			var itemValue reflect.Value
			if self.OnNewItem != nil {
				item = self.OnNewItem()
				itemValue = reflect.ValueOf(item)
			} else {
				itemValue = reflect.New(self.ArrayValue.Type().Elem().Elem())
				item = itemValue.Interface()
			}
			self.ArrayValue.Set(reflect.Append(*self.ArrayValue, itemValue))
		} else if len(r.Key) > 0 {
			// Not implemented, but could be...
			panic("Keys only implemented on MarshalMap, not MarshalArray")
		} else {
			if r.First {
				i = 0
			} else {
				i++
			}
			if i < self.ArrayValue.Len() {
				item = self.ArrayValue.Index(i).Interface()
			}
		}
		if item != nil {
			if self.OnSelectItem != nil {
				return self.OnSelectItem(item), nil, nil
			}
			return MarshalContainer(item), nil, nil
		}
		return nil, nil, nil
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
	n.OnNext = func(sel *Selection, r ListRequest) (next Node, key []*Value, err error) {
		var item interface{}
		key = r.Key
		if r.New {
			item = self.OnNewItem()
			mapKey := reflect.ValueOf(r.Key[0].Value())
			mapReflect.SetMapIndex(mapKey, reflect.ValueOf(item))
		} else if len(r.Key) > 0 {
			mapKey := reflect.ValueOf(r.Key[0].Value())
			itemVal := mapReflect.MapIndex(mapKey)
			if itemVal.IsValid() {
				item = itemVal.Interface()
			}
		} else {
			nextKey := index.NextKey(r.First)
			if nextKey != NO_VALUE {
				key = SetValues(r.Meta.KeyMeta(), nextKey.Interface())
				itemVal := mapReflect.MapIndex(nextKey)
				item = itemVal.Interface()
			}
		}
		if item != nil {
			if self.OnSelectItem != nil {
				return self.OnSelectItem(item), key, nil
			}
			return MarshalContainer(item), key, nil
		}
		return nil, nil, nil
	}
	return n
}
