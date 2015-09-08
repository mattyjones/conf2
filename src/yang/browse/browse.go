package browse

import (
	"yang"
	"reflect"
)

type Browser interface {
	yang.Resource
	RootSelector() (*Selection, error)
	Module() (*yang.Module)
}

type Value struct {
	Type *yang.DataType
	IsList bool
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
	Enter Enter
	ReadValue ReadValue
	Edit Edit
	Exit Exit
	Choose ResolveChoice
	Found bool
	insideList bool
	Resource yang.Resource
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

func (s *Selection) Close() (err error){
	if s.Resource != nil {
		err = s.Resource.Close()
		s.Resource = nil
	}
	return
}

func (s *Selection) CreateChild() error {
	if s.Edit == nil {
		return &browseError{Msg:"Not editable"}
	}
	return s.Edit(CREATE_CHILD, nil)
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

type Iterate func(keys []string, first bool) (hasMore bool, err error)
type Enter func() (child *Selection, err error)
type ReadValue func(val *Value) (error)
type Edit func(op Operation, val *Value) (error)
type Exit func() (error)
type ResolveChoice func(choice *yang.Choice) (m yang.Meta, err error)

func WalkPath(from *Selection, path *Path) (s *Selection, err error) {
	nest := newPathController(path)
	err = walk(from, nest, 0)
	if nest.target != nil {
		nest.target.Resource = nest.resource
	}
	return nest.target, err
}

func WalkExhaustive(selection *Selection) (err error) {
	return walk(selection, &exhaustiveController{MaxDepth:32}, 0)
}

func walk(selection *Selection, controller WalkController, level int) (err error) {
	if yang.IsList(selection.Meta) && !selection.insideList {
		var hasMore bool
		if hasMore, err = controller.ListIterator(selection, level, true); err != nil {
			return
		}
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			selection.insideList = true

			if err = walk(selection, controller, level); err != nil {
				return
			}
			if hasMore, err = controller.ListIterator(selection, level, false); err != nil {
				return
			}
		}
	} else {
		var child *Selection
		i := controller.ContainerIterator(selection, level)
		for i.HasNextMeta() {
			selection.Position = i.NextMeta()
			if choice, isChoice := selection.Position.(*yang.Choice); isChoice {
				if selection.Position, err = selection.Choose(choice); err != nil {
					return
				}
			}
			if yang.IsLeaf(selection.Position) {
				val := &Value{}
				if err = selection.ReadValue(val); err != nil {
					return err
				}
			} else {
				child, err = selection.Enter()
				if child != nil {
					child.Meta = selection.Position.(yang.MetaList)
					defer child.Close()
				}
				if err != nil {
					return
				}
				if !selection.Found {
					continue
				}

				if err = walk(child, controller, level + 1); err != nil {
					return
				}

				if selection.Exit != nil {
					if err = selection.Exit(); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func ReadField(meta yang.HasDataType, obj interface{}, v *Value) error {
	return ReadFieldWithFieldName(yang.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func ReadFieldWithFieldName(fieldName string, meta yang.HasDataType, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	v.Type = meta.GetDataType()
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
		v.IsList = true
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