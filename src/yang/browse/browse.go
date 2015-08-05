package browse

import (
	"yang"
	"fmt"
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
}

type Iterate func(keys []string, first bool) (hasMore bool, err error)
type Select func(ident string) (*Selection, error)
type Read func(val *Value) (error)
type Edit func(op Operation, val *Value) (error)

type Writer interface {
	EnterContainer(yang.MetaList) error
	ExitContainer(yang.MetaList) error

	EnterList(*yang.List) error
	ExitList(*yang.List) error

	UpdateValue(meta yang.Meta, val *Value) error
}

type WriteableSelection struct {
	stack []*Selection
	selection *Selection
}

func NewWriteableSelection(root *Selection) (ws *WriteableSelection) {
	ws = &WriteableSelection{}
	ws.stack = make([]*Selection, 10)
	ws.selection = root
	ws.stack[0] = ws.selection
	return ws
}

func (ws *WriteableSelection) EnterContainer(m yang.MetaList) (err error) {
	ws.selection.Position = m
	if ws.selection.Edit == nil {
		return EditNotImplemented(m)
	}
	ws.selection.Edit(CREATE_CHILD, nil)
	err = ws.pushSelection(m)
	return
}

func (ws *WriteableSelection) ExitContainer(m yang.MetaList) (err error) {
	if err = ws.popSelection(); err == nil {
		err = ws.selection.Edit(POST_CREATE_CHILD, nil)
	}
	return nil
}

func (ws *WriteableSelection) pushSelection(m yang.MetaList) (err error) {
	if ws.selection, err = ws.selection.Select(m.GetIdent()); err == nil {
		if ws.selection == nil {
			msg := fmt.Sprint("expected selector for property ", m.GetIdent())
			err = &browseError{Msg:msg}
		} else {
			ws.selection.Meta = m
			ws.stack = append(ws.stack, ws.selection)
		}
	}
	return err
}

func (ws *WriteableSelection) popSelection() (err error) {
	if len(ws.stack) == 0 {
		return &browseError{Msg:"Empty selection stack"}
	}
	ws.selection = ws.stack[len(ws.stack) - 1]
	ws.stack = ws.stack[0:len(ws.stack) - 1]
	return
}

func (ws *WriteableSelection) EnterList(m *yang.List) (err error) {
	ws.selection.Position = m
	if ws.selection.Edit == nil {
		return EditNotImplemented(m)
	}
	if err = ws.selection.Edit(CREATE_LIST, nil); err != nil {
		ws.pushSelection(m)
	}
	return err
}

func (ws *WriteableSelection) ExitList(m *yang.List) (err error) {
	if err = ws.popSelection(); err != nil {
		ws.selection.Edit(POST_CREATE_LIST, nil)
	}
	return err
}

func (ws *WriteableSelection) UpdateValue(m yang.Meta, val *Value) (err error) {
	ws.selection.Position = m
	var unexpected *Selection
	if unexpected, err = ws.selection.Select(m.GetIdent()); err == nil {
		if ws.selection.Edit == nil {
			err = EditNotImplemented(m)
		} else {
			err = ws.selection.Edit(UPDATE_VALUE, val)
		}
	} else if (unexpected != nil) {
		msg := fmt.Sprint("unexpected leaf selector for property ", m.GetIdent())
		err = &browseError{Msg:msg}
	}

	return
}

type Operation int
const (
	CREATE_CHILD Operation = 1 + iota
	POST_CREATE_CHILD
	CREATE_LIST
	POST_CREATE_LIST
	UPDATE_VALUE
	DELETE_CHILD
)

func Walk(from *Selection, path *Path, to Writer) (err error) {
	nest := &readController{path:path, maxLevel:100}
	return read(from, to, nest)
}

type readController struct {
	level int
	maxLevel int
	path *Path
}

func (n *readController) isMaxLevel() bool {
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

func (n *readController) keys() []string {
	if n.path == nil || n.level >= len(n.path.Segments) {
		return []string{}
	}
	return n.path.Segments[n.level].Keys
}

func (n *readController) matches(ident string) bool {
	if n.path == nil || n.level >= len(n.path.Segments) {
		return true
	}
	return n.path.Segments[n.level].Ident == ident
}

func (n *readController) recurse() (*readController) {
	return &readController{path:n.path, level: n.level + 1, maxLevel:n.maxLevel}
}

func read(selection *Selection, wtr Writer, controller *readController) (err error) {
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
			val := Value{}
			if err = selection.Read(&val); err != nil {
				return
			}
			if err = wtr.UpdateValue(selection.Position, &val); err != nil {
				return
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
				list := child.Meta.(*yang.List)
				if err = wtr.EnterList(list); err != nil {
					return
				}

				for more {
					if err = read(child, wtr, controller.recurse()); err != nil {
						return
					}

					if more, err = child.Iterate(controller.keys(), false); err != nil {
						return
					}
				}

				if err = wtr.ExitList(list); err != nil {
					return
				}

			} else {
				if err = wtr.EnterContainer(child.Meta); err != nil {
					return
				}
				if err = read(child, wtr, controller.recurse()); err != nil {
					return
				}
				if err = wtr.ExitContainer(child.Meta); err != nil {
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