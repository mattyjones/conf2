package browse

import (
	"schema"
	"reflect"
	"fmt"
)

type Browser interface {
	schema.Resource
	RootSelector() (Selection, error)
	Module() (*schema.Module)
}

type Value struct {
	Type *schema.DataType
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


type WalkController interface {
	ListIterator(s Selection, level int, first bool) (hasMore bool, err error)
	ContainerIterator(s Selection, level int) schema.MetaIterator
	CloseSelection(s Selection) error
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

func WalkPath(from Selection, path *Path) (s Selection, err error) {
	nest := path.FindTargetController()
	err = walk(from, nest, 0)
	return nest.target, err
}

func WalkExhaustive(selection Selection, controller WalkController) (err error) {
	return walk(selection, controller, 0)
}

func walk(selection Selection, controller WalkController, level int) (err error) {
	state := selection.WalkState()
	if schema.IsList(state.Meta) && !state.insideList {
		var hasMore bool
		if hasMore, err = controller.ListIterator(selection, level, true); err != nil {
			return
		}
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			state.insideList = true

			if err = walk(selection, controller, level); err != nil {
				return
			}
			if hasMore, err = controller.ListIterator(selection, level, false); err != nil {
				return
			}
		}
	} else {
		var child Selection
		i := controller.ContainerIterator(selection, level)
		for i.HasNextMeta() {
			state.Position = i.NextMeta()
			if choice, isChoice := state.Position.(*schema.Choice); isChoice {
				if state.Position, err = selection.Chooze(choice); err != nil {
					return
				}
			}
			if schema.IsLeaf(state.Position) {
				val := &Value{}
				if err = selection.Read(val); err != nil {
					return err
				}
			} else {
				child, err = selection.Select()
				if child != nil {
					child.WalkState().Meta = state.Position.(schema.MetaList)
					defer schema.CloseResource(child)
				}
				if err != nil {
					return
				}
				if !state.Found {
					continue
				}

				if err = walk(child, controller, level + 1); err != nil {
					return
				}

				err = selection.Unselect()
			}
		}
	}
	return
}

func ReadField(meta schema.HasDataType, obj interface{}, v *Value) error {
	return ReadFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func ReadFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	v.Type = meta.GetDataType()
	switch tmeta := meta.(type) {
	case *schema.Leaf:
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
	case *schema.LeafList:
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

func WriteField(meta schema.Meta, obj interface{}, v *Value) error {
	return WriteFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj, v)
}

func WriteFieldWithFieldName(fieldName string, meta schema.Meta, obj interface{}, v *Value) error {
	objType := reflect.ValueOf(obj).Elem()
	if ! objType.IsValid() {
		return &browseError{Msg:fmt.Sprintf("Cannot find property \"%s\" on invalid or nil %s", fieldName, reflect.TypeOf(obj))}
	}
	value := objType.FieldByName(fieldName)
	switch tmeta := meta.(type) {
		case *schema.Leaf:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			value.SetBool(v.Bool)
		case "int32":
			value.SetInt(int64(v.Int))
		default:
			value.SetString(v.Str)
		}
		case *schema.LeafList:
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