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
	Enter Enter
	ReadValue ReadValue
	Edit Edit
	Exit Exit
	Choose ResolveChoice
	Found bool
	insideList bool
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
	return nest.target, err
}

type WalkController interface {
	ListIterator(s *Selection, level int, first bool) (hasMore bool, err error)
	ContainerIterator(s *Selection, level int) yang.MetaIterator
}

type exhaustiveController struct {}

func (e exhaustiveController) ListIterator(s *Selection, level int, first bool) (hasMore bool, err error) {
	return s.Iterate([]string{}, first)
}
func (e exhaustiveController) ContainerIterator(s *Selection, level int) yang.MetaIterator {
	return yang.NewMetaListIterator(s.Meta, true)
}

func WalkExhaustive(selection *Selection) (err error) {
	return walk(selection, exhaustiveController{}, 0)
}

func walk(selection *Selection, controller WalkController, level int) (err error) {
	if yang.IsList(selection.Meta) && !selection.insideList {
		var hasMore bool
		hasMore, err = controller.ListIterator(selection, level, true)
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			selection.insideList = true

			if err = walk(selection, controller, level); err != nil {
				return
			}
			hasMore, err = controller.ListIterator(selection, level, false)
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
				if err != nil {
					return
				}
				if !selection.Found {
					continue
				}
				child.Meta = selection.Position.(yang.MetaList)

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