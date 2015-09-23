package browse

import (
	"schema"
	"reflect"
	"fmt"
)

type Browser interface {
	RootSelector() (Selection, error)
	Module() (*schema.Module)
}

type WalkController interface {
	ListIterator(s Selection, level int, first bool) (hasMore bool, err error)
	ContainerIterator(s Selection, level int) schema.MetaIterator
	CloseSelection(s Selection) error
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
	if schema.IsList(state.Meta) && !state.InsideList {
		var hasMore bool
		if hasMore, err = controller.ListIterator(selection, level, true); err != nil {
			return
		}
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			state.InsideList = true

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
				if state.Position, err = selection.Choose(choice); err != nil {
					return
				}
			}
			if schema.IsLeaf(state.Position) {
				// only walking here, not interested in value
				if _, err = selection.Read(state.Position.(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := state.Position.(schema.MetaList)
				child, err = selection.Select(metaList)
				if err != nil {
					return
				} else if child == nil {
					continue
				}
				child.WalkState().Meta = metaList
				defer schema.CloseResource(child)

				if err = walk(child, controller, level + 1); err != nil {
					return
				}

				err = selection.Unselect(metaList)
			}
		}
	}
	return
}

func ReadField(meta schema.HasDataType, obj interface{}) (*Value, error) {
	return ReadFieldWithFieldName(schema.MetaNameToFieldName(meta.GetIdent()), meta, obj)
}

func ReadFieldWithFieldName(fieldName string, meta schema.HasDataType, obj interface{}) (*Value, error) {
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	_, isList := meta.(*schema.LeafList)
	switch meta.GetDataType().Format {
	case schema.FMT_BOOLEAN:
		if isList {
			return &Value{Boollist:value.Interface().([]bool)}, nil
		}
		return &Value{Bool:value.Bool()}, nil
	case schema.FMT_INT32:
		if isList {
			return &Value{Intlist:value.Interface().([]int), IsList:true}, nil
		}
		return &Value{Int:int(value.Int())}, nil
	default:
		if isList {
			return &Value{Strlist:value.Interface().([]string), IsList:true}, nil
		}
		return &Value{Str:value.String()}, nil
	}
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
	case schema.FMT_BOOLEAN:
		if v.IsList {
			value.Set(reflect.ValueOf(v.Boollist))
		} else {
			value.SetBool(v.Bool)
		}
	case schema.FMT_INT32:
		if v.IsList {
			value.Set(reflect.ValueOf(v.Intlist))
		} else {
			value.SetInt(int64(v.Int))
		}
	case schema.FMT_STRING:
		if v.IsList {
			value.Set(reflect.ValueOf(v.Strlist))
		} else {
			value.SetString(v.Str)
		}
	default:
		return NotImplemented(meta)
	}
	return nil
}