package browse

import (
	"yang"
	"strings"
	"fmt"
	"reflect"
	"log"
)

type browseError struct {
	Code ResponseCode
	Msg string
}

func (err *browseError) Error() string {
	return err.Msg
}

type Path struct {
	Segments []PathSegment
	URL string
}

type PathSegment struct {
	Path *Path
	Index int
	Ident string
	Keys []string
}

func NewPath(path string) (p *Path) {
	p = &Path{}
	qmark := strings.Index(path, "?")
	if qmark >= 0 {
		p.URL = path[:qmark]
	} else {
		p.URL = path
	}
	segments := strings.Split(p.URL, "/")
	p.Segments = make([]PathSegment, len(segments))
	for i, segment := range segments {
		p.Segments[i] = PathSegment{Path:p, Index:i}
		p.Segments[i].parseSegment(segment)
	}
	return
}

func (ps *PathSegment) parseSegment(segment string) {
	equalsMark := strings.Index(segment, "=")
	if equalsMark >= 0 {
		ps.Ident = segment[:equalsMark]
		ps.Keys = strings.Split(segment[equalsMark + 1:], ",")
	} else {
		ps.Ident = segment
	}
}

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
	ListIterator ListIterator
	Selector Selector
	Reader Reader
	Edit Editor
}

type ListIterator func(keys []string, first bool) (hasMore bool, err error)
type Selector func(ident string) (*Selection, error)
type Reader func(val *Value) (error)
type Editor func(op Operation, val *Value) (error)

type Writer interface {
	EnterContainer(yang.MetaList) error
	ExitContainer(yang.MetaList) error

	EnterList(*yang.List) error
	ExitList(*yang.List) error

	EnterListItem(*yang.List) error
	ExitListItem(*yang.List) error

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
	if ws.selection, err = ws.selection.Selector(m.GetIdent()); err == nil {
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

func (ws *WriteableSelection) EnterListItem(m *yang.List) (err error) {
	ws.selection.Position = m
	if ws.selection.Edit == nil {
		return EditNotImplemented(m)
	}
	if err = ws.selection.Edit(CREATE_LIST_ITEM, nil); err != nil {
		err = ws.pushSelection(m)
	}
	return err
}

func (ws *WriteableSelection) ExitListItem(m *yang.List) (err error) {
	if err = ws.popSelection(); err == nil {
		err = ws.selection.Edit(POST_CREATE_LIST_ITEM, nil)
	}
	return err
}

func (ws *WriteableSelection) UpdateValue(m yang.Meta, val *Value) error {
	ws.selection.Position = m
	if ws.selection.Edit == nil {
		return EditNotImplemented(m)
	}
	return ws.selection.Edit(UPDATE_VALUE, val)
}


type Operation int
const (
	CREATE_CHILD Operation = 1 + iota
	CREATE_LIST
	CREATE_LIST_ITEM
	POST_CREATE_CHILD
	POST_CREATE_LIST
	POST_CREATE_LIST_ITEM
	UPDATE_VALUE
	DELETE_CHILD
)

type NullWriter struct {
}

func (NullWriter) EnterContainer(yang.MetaList) error {
	return nil
}

func (NullWriter) ExitContainer(yang.MetaList) error {
	return nil
}

func (NullWriter) EnterList(*yang.List) error {
	return nil
}

func (NullWriter) ExitList(*yang.List) error {
	return nil
}

func (NullWriter) EnterListItem(*yang.List) error {
	return nil
}

func (NullWriter) ExitListItem(*yang.List) error {
	return nil
}

func (NullWriter) UpdateValue(meta yang.Meta, val *Value) error {
	return nil
}

type DebuggingWriter struct {
	Delegate Writer
}

func (w *DebuggingWriter) EnterContainer(m yang.MetaList) error {
	log.Println("Entering Container", m.GetIdent())
	return w.Delegate.EnterContainer(m)
}

func (w *DebuggingWriter) ExitContainer(m yang.MetaList) error {
	log.Println("Exiting Container", m.GetIdent())
	return w.Delegate.ExitContainer(m)
}

func (w *DebuggingWriter) EnterList(m *yang.List) error {
	log.Println("Entering List", m.GetIdent())
	return w.Delegate.EnterList(m)
}

func (w *DebuggingWriter) ExitList(m *yang.List) error {
	log.Println("Exiting List", m.GetIdent())
	return w.Delegate.ExitList(m)
}

func (w *DebuggingWriter) EnterListItem(m *yang.List) error {
	log.Println("Entering List Item", m.GetIdent())
	return w.Delegate.EnterListItem(m)
}

func (w *DebuggingWriter) ExitListItem(m *yang.List) error {
	log.Println("Existing List Item", m.GetIdent())
	return w.Delegate.ExitListItem(m)
}

func (w *DebuggingWriter) UpdateValue(m yang.Meta, v *Value) error {
	log.Println("Updating Value", m.GetIdent())
	return w.Delegate.UpdateValue(m, v)
}

func Advance(selection *Selection, path *Path) (child *Selection, err error) {
	return readContainer(selection, NullWriter{}, path, 0)
}

func Transfer(from *Selection, to Writer) (err error) {
	_, err = readContainer(from, to, nil, 0)
	return
}

func readContainer(from *Selection, to Writer, path *Path, level int) (selection *Selection, err error) {
	var segment PathSegment
	if path != nil && len(path.Segments) > level {
		segment = path.Segments[level]
	}
	selection = from
	var fromChild *Selection
	i := yang.NewMetaListIterator(from.Meta, true)
	var isContainer bool
	for i.HasNextMeta() {
		meta := i.NextMeta()
		if segment.Ident != "" {
			if meta.GetIdent() != segment.Ident {
				continue
			}
		}
		val := Value{}
		if fromChild, err = from.Selector(meta.GetIdent()); err != nil {
			return
		}
		if from.Position == nil {
			continue
		}
		if fromChild != nil {
			fromChild.Meta, isContainer = from.Position.(yang.MetaList)
			if !isContainer {
				msg := fmt.Sprint("leaf node returned a selector:", from.Position.GetIdent())
				return nil, &browseError{Msg:msg}
			}
			if yang.IsList(from.Position) {


				if selection, err = readList(fromChild, to, path, level + 1); err != nil {
					return
				}


			} else {

				if err = to.EnterContainer(fromChild.Meta); err != nil {
					return
				}

				if selection, err = readContainer(fromChild, to, path, level + 1); err != nil {
					return
				}

				if err = to.ExitContainer(fromChild.Meta); err != nil {
					return
				}

			}
		} else {
			if err = from.Reader(&val); err != nil {
				return
			}
			if err = to.UpdateValue(from.Position, &val); err != nil {
				return
			}
		}
	}
	return
}

func readList(from *Selection, to Writer, path *Path, level int) (selection *Selection, err error) {
	var segment PathSegment
	if path != nil && len(path.Segments) > level {
		segment = path.Segments[level]
	}
	selection = from
	metaList := from.Meta.(*yang.List)
	var hasMore bool
	if hasMore, err = from.ListIterator(segment.Keys, true); err != nil {
		return
	} else if (!hasMore) {
		return
	}
	if err = to.EnterList(metaList); err != nil {
		return
	}
	for hasMore {
		// list in list is illegal AFAIK so assume container
		if err = to.EnterListItem(metaList); err != nil {
			return
		}

		if selection, err = readContainer(from, to, path, level + 1); err != nil {
			return
		}

		if err = to.ExitListItem(metaList); err != nil {
			return
		}

		if hasMore, err = from.ListIterator(segment.Keys, false); err != nil {
			return
		}
	}
	if err = to.ExitList(metaList); err != nil {
		return
	}
	return
}

type ResponseCode int
const (
	UNSPECIFIED ResponseCode = iota
	NOT_IMPLEMENTED
	NOT_FOUND
	MISSING_KEY
)

func EditNotImplemented(meta yang.Meta) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("editing of \"%s\" not implemented", meta.GetIdent())}
}

func NotImplementedByName(ident string) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("browsing of \"%s\" not implemented", ident)}
}

func NotImplemented(meta yang.Meta) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("browsing of \"%s\" not implemented", meta.GetIdent())}
}

func NotFound(key string) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("item identified with key \"%s\" not found", key)}
}

func ListKeyRequired() error {
	return &browseError{Code:MISSING_KEY, Msg:fmt.Sprintf("List key required")}
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