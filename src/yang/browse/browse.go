package browse

import (
	"yang"
	"reflect"
)

type Browser interface {
	RootSelector() (*Selection, error)
}

type Value struct {
	Bool bool
	Int int
	Str string
	Float float32
	Intlist []int
	Strlist []string
	Boollist []bool
	Keys []string
}

type Selection struct {
	Meta yang.MetaList
	Position yang.Meta
	Iterate Iterate
	Select Select
	Read Read
	Edit Edit
	Exit Exit
	Found bool
}

func (s *Selection) CreateChild() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(CREATE_CHILD, nil)
}

func (s *Selection) FinishCreateChild() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	if yang.IsList(s.Position) {
		return s.Edit(POST_CREATE_LIST, nil)
	}
	return s.Edit(POST_CREATE_CHILD, nil)
}

func (s *Selection) DeleteChild() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(DELETE_CHILD, nil)
}

func (s *Selection) DeleteList() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(DELETE_LIST, nil)
}

func (s *Selection) SetValue(val *Value) error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(UPDATE_VALUE, val)
}

func (s *Selection) CreateList() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(CREATE_LIST, nil)
}

func (s *Selection) FinishCreateList() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(POST_CREATE_LIST, nil)
}

type Iterate func(keys []string, first bool) (hasMore bool, err error)
type Select func(ident string) (*Selection, error)
type Read func(val *Value) (error)
type Edit func(op Operation, val *Value) (error)
type Exit func() (error)

func Walk(from *Selection, path *Path) (err error) {
	nest := &walkController{path:path, maxLevel:100}
	return walk(from, nest)
}

type walkController struct {
	level int
	maxLevel int
	path *Path
}

func (n *walkController) isMaxLevel() bool {
	if n.path != nil {
		if (n.path.Depth > 0) {
			calcDepth := len(n.path.Segments) + n.path.Depth
			if calcDepth < n.maxLevel {
				return n.level + 1 >= calcDepth
			}
		}
	}
	return n.level + 1 >= n.maxLevel
}

func (n *walkController) keys() []string {
	if n.path == nil || n.level >= len(n.path.Segments) {
		return []string{}
	}
	return n.path.Segments[n.level].Keys
}

func (n *walkController) matches(ident string) bool {
	if n.path == nil || n.level >= len(n.path.Segments) {
		return true
	}
	return n.path.Segments[n.level].Ident == ident
}

func (n *walkController) recurse() (*walkController) {
	return &walkController{path:n.path, level: n.level + 1, maxLevel:n.maxLevel}
}

func walk(selection *Selection, controller *walkController) (err error) {
	var child *Selection
	i := yang.NewMetaListIterator(selection.Meta, true)
	for i.HasNextMeta() {
		meta := i.NextMeta()
		if !controller.matches(meta.GetIdent()) {
			continue
		}
		child, err = selection.Select(meta.GetIdent())
		if selection.Position == nil {
			continue
		}
		if child == nil {
			val := &Value{}
			if err = selection.Read(val); err != nil {
				return err
			}
		} else if ! controller.isMaxLevel() {
			child.Meta = selection.Position.(yang.MetaList)
			if (child.Iterate != nil) {
				var more bool
				if more, err = child.Iterate(controller.keys(), true); err != nil {
					return
				} else if (!more) {
					continue
				}
				for more {
					if err = walk(child, controller.recurse()); err != nil {
						return
					}

					if more, err = child.Iterate(controller.keys(), false); err != nil {
						return
					}
				}
			} else {
				if err = walk(child, controller.recurse()); err != nil {
					return
				}
			}
			if selection.Exit != nil {
				if err = selection.Exit(); err != nil {
					return
				}
			}
		}
	}
	return
}

func ReadField(meta yang.Meta, obj interface{}, v *Value) error {
	return ReadFieldWithFieldName(yang.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func ReadFieldWithFieldName(fieldName string, meta yang.Meta, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	switch tmeta := meta.(type) {
	case *yang.Leaf:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			if value.Bool() {
				v.Bool = true
			}
			v.Bool = false
		case "int32":
			v.Int = int(value.Int())
		default:
			v.Str = value.String()
		}
	case *yang.LeafList:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			v.Boollist = value.Interface().([]bool)
		case "int32":
			v.Intlist = value.Interface().([]int)
		default:
			v.Strlist = value.Interface().([]string)
		}
	default:
		return NotImplemented(meta)
	}
	return nil
}

func WriteField(meta yang.Meta, obj interface{}, v *Value) error {
	return WriteFieldWithFieldName(yang.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func WriteFieldWithFieldName(fieldName string, meta yang.Meta, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	switch tmeta := meta.(type) {
		case *yang.Leaf:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			value.SetBool(v.Bool)
		case "int32":
			value.SetInt(int64(v.Int))
		default:
			value.SetString(v.Str)
		}
		case *yang.LeafList:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			value.Set(reflect.ValueOf(v.Boollist))
		case "int32":
			value.Set(reflect.ValueOf(v.Intlist))
		default:
			value.Set(reflect.ValueOf(v.Strlist))
		}
		default:
		return NotImplemented(meta)
	}
	return nil
}