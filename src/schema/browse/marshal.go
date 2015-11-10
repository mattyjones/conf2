package browse

import (
	"reflect"
	"schema"
)

// Uses reflection to marshal data into go structs

func Marshal(obj interface{}, sel *Selection) (err error) {
	n := &MarshalContainer{Obj: obj}
	err = UpsertByNode(sel, sel.Node(), n)
	return
}

type MarshalContainer struct {
	Obj         interface{}
	OnSelectNil SelectFunc
}

func (s *MarshalContainer) Select(selection *Selection, meta schema.MetaList) (Node, error) {
	objType := reflect.ValueOf(s.Obj).Elem()
	fieldName := schema.MetaNameToFieldName(meta.GetIdent())
	value := objType.FieldByName(fieldName)
	if schema.IsList(meta) {
		return MarshalList(value.Interface().([]interface{})), nil
	} else {
		if value.Kind() == reflect.Struct {
			return &MarshalContainer{Obj: value.Addr().Interface()}, nil
		} else if value.CanAddr() {
			return &MarshalContainer{Obj: value.Interface()}, nil
		}
	}
	return nil, nil
}

func (s *MarshalContainer) String() string {
	return "Marshal " + reflect.TypeOf(s.Obj).Name()
}

func (s *MarshalContainer) Next(state *Selection, meta *schema.List, keys []*Value, isFirst bool) (next Node, err error) {
	panic("Not implemented")
}

func (s *MarshalContainer) Read(state *Selection, meta schema.HasDataType) (*Value, error) {
	panic("Not implemented")
}

func (s *MarshalContainer) Choose(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
	panic("Not implemented")
}

func (s *MarshalContainer) Unselect(state *Selection, meta schema.MetaList) error {
	return nil
}

func (s *MarshalContainer) Action(state *Selection, meta *schema.Rpc, input Node) (output *Selection, err error) {
	panic("Not implemented")
}

func (s *MarshalContainer) Write(selection *Selection, meta schema.Meta, op Operation, val *Value) (err error) {
	switch op {
	case UPDATE_VALUE:
		return WriteField(meta.(schema.HasDataType), s.Obj, val)
	}
	return
}

func MarshalList(list []interface{}) Node {
	panic("Not implemented")
	return nil
}
