package browse

import (
	"yang"
	"strings"
	"fmt"
	"reflect"
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
	GetSelector() Selection
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
//	Selection Selection
//	CaseMeta yang.Meta
}

func (v *Value) SetString(meta yang.Meta, str string) {
	// TODO: Coerse value according to meta
	v.Str = str
}

func (v *Value) GetString(meta yang.Meta) string {
	// TODO: Coerse value according to meta
	return v.Str
}

func (v *Value) SetInt(meta yang.Meta, i int) {
	v.Int = i
}

type Selection func(op Operation, meta yang.Meta, v *Visitor) error

type Operation int
const (
	CREATE_CHILD Operation = 1 + iota
	CREATE_LIST
	CREATE_LIST_ITEM
	POST_CREATE_CHILD
	POST_CREATE_LIST
	POST_CREATE_LIST_ITEM
	READ_VALUE
	UPDATE_VALUE
	DELETE_CHILD
	SELECT_CHILD
)

func (v *Visitor) MakeSelection(path *Path, initialMeta yang.Meta) (meta yang.Meta, err error) {
	meta = initialMeta
	var val Value
	for i, segment := range path.Segments {
		val.Keys = segment.Keys
		v.Selection(SELECT_CHILD, meta, v)
		lastSegment := (i == len(path.Segments) - 1)
		if !lastSegment {
			// TODO: Need to handle case meta
			if metaList, isList := meta.(yang.MetaList); !isList {
				return nil, &browseError{Msg:"Invalid path"}
			} else {
				i := yang.NewMetaListIterator(metaList, true)
				meta = yang.FindByIdent(i, segment.Ident)
			}
		}
	}
	return
}

type Visitor struct {
	Out			*Visitor
	Val 		Value
	Selection	Selection
	Position    yang.Meta
}

func NewVisitor(selection Selection) (v *Visitor) {
	v = &Visitor{Selection:selection}
	return
}

func (v *Visitor) EnterContainer(meta yang.Meta) (err error) {
	if err = v.Out.Selection(CREATE_CHILD, meta, v.Out); err == nil {
		err = v.Out.Selection(SELECT_CHILD, meta, v.Out)
	}
	return
}

func (v *Visitor) EnterList(meta yang.Meta) (err error) {
	if err = v.Out.Selection(CREATE_LIST, meta, v.Out); err == nil {
		err = v.Out.Selection(SELECT_CHILD, meta, v.Out)
	}
	return
}

func (v *Visitor) EnterListItem(meta yang.Meta) (err error) {
	if err = v.Out.Selection(CREATE_LIST_ITEM, meta, v.Out); err == nil {
		err = v.Out.Selection(SELECT_CHILD, meta, v.Out)
	}
	return
}

func (v *Visitor) ExitContainer(meta yang.Meta) error {
	return v.Out.Selection(POST_CREATE_CHILD, meta, v.Out)
}

func (v *Visitor) ExitList(meta yang.Meta) error {
	return v.Out.Selection(POST_CREATE_LIST, meta, v.Out)
}

func (v *Visitor) ExitListItem(meta yang.Meta) error {
	return v.Out.Selection(POST_CREATE_LIST_ITEM, meta, v.Out)
}

func (v *Visitor) Send(meta yang.Meta) error {
	v.Out.Val = v.Val
	return v.Out.Selection(UPDATE_VALUE, meta, v.Out)
}

type ResponseCode int
const (
	UNSPECIFIED ResponseCode = iota
	NOT_IMPLEMENTED
	NOT_FOUND
	MISSING_KEY
)

func NotImplemented(meta yang.Meta) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("browsing of \"%s\" not implemented", meta.GetIdent())}
}

func (v *Visitor) NotFound(key string) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("item identified with key \"%s\" not found", key)}
}

func (v *Visitor) ListKeyRequired() error {
	return &browseError{Code:MISSING_KEY, Msg:fmt.Sprintf("List key required")}
}

func UseReflection(op Operation, meta yang.Meta, obj interface{}, v *Visitor) error {
	fieldName := yang.MetaNameToFieldName(meta.GetIdent())
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	switch tmeta := meta.(type) {
	case *yang.Leaf:
		switch op {
		case READ_VALUE:
			switch tmeta.GetDataType().Resolve().Ident {
			case "boolean":
				if value.Bool() {
					v.Val.Bool = true
				}
				v.Val.Bool = false
			case "int32":
				v.Val.Int = int(value.Int())
			default:
				v.Val.Str = value.String()
			}
			v.Send(meta)

		default:
			return NotImplemented(meta)
		}
	case *yang.LeafList:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			v.Val.Boollist = value.Interface().([]bool)
		case "int32":
			v.Val.Intlist = value.Interface().([]int)
		default:
			v.Val.Strlist = value.Interface().([]string)
		}
		v.Send(meta)
	default:
		return NotImplemented(meta)
	}
	return nil
}
