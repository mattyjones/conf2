package browse

import (
	"io"
	"schema"
	"bufio"
	"strconv"
	"errors"
)

const QUOTE = '"';
const COLON = ':';
const OPEN_OBJ = '{';
const CLOSE_OBJ = '}';
const OPEN_ARRAY = '[';
const CLOSE_ARRAY = ']';
const COMMA = ',';

type JsonWriter struct {
	out *bufio.Writer
	meta *schema.Module
	firstInContainer bool
	startingInsideList bool
	firstWrite bool
	closeArrayOnExit bool
}

func NewJsonWriter(out io.Writer, module *schema.Module) *JsonWriter {
	return &JsonWriter{
		out:bufio.NewWriter(out),
		meta:module,
		firstInContainer:true,
	}
}

func NewJsonFragmentWriter(out io.Writer) *JsonWriter {
	return &JsonWriter{
		out:bufio.NewWriter(out),
		firstInContainer:true,
	}
}

func (json *JsonWriter) Selector(path *Path, strategy Strategy) (s Selection, state *WalkState, err error) {
	if strategy != INSERT && strategy != UPSERT {
		return nil, nil, errors.New("Only [UP,IN]SERT strategy is supported. Consider using bucket first")
	}
	s, _ = json.selectJson()
	if json.meta != nil {
		s, state, err = WalkPath(NewWalkState(json.meta), s, path)
	}
	return
}

func (self *JsonWriter) Module() *schema.Module {
	return nil
}

func (json *JsonWriter) selectJson() (Selection, error) {
	s := &MySelection{}
	var created Selection
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (child Selection, err error) {
		nest := created
		created = nil
		return nest, nil
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case BEGIN_EDIT:
			_, err = json.out.WriteRune(OPEN_OBJ)
			json.startingInsideList = schema.IsList(meta)
			json.firstWrite = true
			return err
		case END_EDIT:
			if err = json.conditionallyCloseArrayOnLastWrite(); err == nil {
				if _, err = json.out.WriteRune(CLOSE_OBJ); err == nil {
					err = json.out.Flush()
				}
			}
		case CREATE_CHILD:
			err = json.beginContainer(meta.GetIdent())
			created, _ = json.selectJson()
		case POST_CREATE_CHILD:
			err = json.endContainer()
		case CREATE_LIST_ITEM:
			if err = json.conditionallyOpenArrayOnFirstWrite(meta.GetIdent()); err == nil {
				err = json.beginArrayItem()
			}
			created, _ = json.selectJson()
		case POST_CREATE_LIST_ITEM:
			err = json.endArrayItem()
		case CREATE_LIST:
			err = json.beginList(meta.GetIdent())
			created, _ = json.selectJson()
		case POST_CREATE_LIST:
			return json.endList()
		case UPDATE_VALUE:
			err = json.writeValue(meta, v)
		default:
			err = &browseError{Msg:"Operation not supported"}
		}
		json.firstWrite = false
		return
	}
	s.OnNext = func(state *WalkState, meta *schema.List, keys []*Value, first bool) (next Selection, err error) {
		next = created
		created = nil
		return next, nil
	}
	return s, nil
}

func (json *JsonWriter) conditionallyOpenArrayOnFirstWrite(ident string) error {
	var err error
	if json.firstWrite && json.startingInsideList {
		json.closeArrayOnExit = true
		err = json.beginList(ident)
	}
	return err
}

func (json *JsonWriter) conditionallyCloseArrayOnLastWrite() error {
	var err error
	if json.closeArrayOnExit {
		err = json.endList()
	}
	return err
}

func (json *JsonWriter) beginList(ident string) (err error) {
	if err = json.writeIdent(ident); err == nil {
		_, err = json.out.WriteRune(OPEN_ARRAY)
		json.firstInContainer = true;
	}
	return
}

func (json *JsonWriter) endList() (err error) {
	_, err = json.out.WriteRune(CLOSE_ARRAY);
	json.firstInContainer = false
	return
}

func (json *JsonWriter) beginContainer(ident string) (err error) {
	if err = json.writeIdent(ident); err != nil {
		return
	}
	if err = json.beginObject(); err != nil {
		return
	}
	return
}

func (json *JsonWriter) endContainer() (err error) {
	json.firstInContainer = false;
	_, err = json.out.WriteRune(CLOSE_OBJ)
	return
}

func (json *JsonWriter) writeValue(meta schema.Meta, v *Value) (err error) {
	json.writeIdent(meta.GetIdent());
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		err = json.writeBool(v.Bool)
	case schema.FMT_INT32:
		err = json.writeInt(v.Int)
	case schema.FMT_STRING, schema.FMT_ENUMERATION:
		err = json.writeString(v.Str)
	case schema.FMT_BOOLEAN_LIST:
		if _, err = json.out.WriteRune(OPEN_ARRAY); err != nil {
			return
		}
		for i, b := range v.Boollist {
			if i > 0 {
				if _,err = json.out.WriteRune(COMMA); err != nil {
					return
				}
			}
			if err = json.writeBool(b); err != nil {
				return
			}
		}
		_, err = json.out.WriteRune(CLOSE_ARRAY)
	case schema.FMT_INT32_LIST:
		if _, err = json.out.WriteRune(OPEN_ARRAY); err != nil {
			return
		}
		for i, n := range v.Intlist {
			if i > 0 {
				if _, err = json.out.WriteRune(COMMA); err != nil {
					return
				}
			}
			if err = json.writeInt(n); err != nil {
				return
			}
		}
		_, err = json.out.WriteRune(CLOSE_ARRAY)
	case schema.FMT_STRING_LIST, schema.FMT_ENUMERATION_LIST:
		if _, err = json.out.WriteRune(OPEN_ARRAY); err != nil {
			return
		}
		for i, s := range v.Strlist {
			if i > 0 {
				if _, err = json.out.WriteRune(COMMA); err != nil {
					return
				}
			}
			if err = json.writeString(s); err != nil {
				return
			}
		}
		_, err = json.out.WriteRune(CLOSE_ARRAY)
	}
	return
}

func (json *JsonWriter) writeBool(b bool) error {
	if b {
		return json.writeString("true")
	} else {
		return json.writeString("false")
	}
}

func (json *JsonWriter) writeInt(i int) (err error) {
	_, err = json.out.WriteString(strconv.Itoa(i))
	return
}

func (json *JsonWriter) writeString(s string) (err error) {
	if _, err = json.out.WriteRune(QUOTE); err == nil {
		if _, err = json.out.WriteString(s); err == nil {
			_, err = json.out.WriteRune(QUOTE);
		}
	}
	return
}

func (json *JsonWriter) beginArrayItem() (err error) {
	if err = json.writeDelim(); err != nil {
		return
	}
	if err = json.beginObject(); err != nil {
		return
	}
	return
}

func (json *JsonWriter) endArrayItem() (err error) {
	json.firstInContainer = false;
	_, err = json.out.WriteRune(CLOSE_OBJ);
	return
}

func (json *JsonWriter) beginObject() (err error) {
	if err == nil {
		_, err = json.out.WriteRune(OPEN_OBJ);
		json.firstInContainer = true;
	}
	return
}

func (json *JsonWriter) writeIdent(ident string) (err error) {
	if err = json.writeDelim(); err != nil {
		return
	}
	if _, err = json.out.WriteRune(QUOTE); err != nil {
		return
	}
	if _, err = json.out.WriteString(ident); err != nil {
		return
	}
	if _, err = json.out.WriteRune(QUOTE); err != nil {
		return
	}
	_, err = json.out.WriteRune(COLON)
	return
}

func (json *JsonWriter) writeDelim() (err error) {
	if json.firstInContainer {
		json.firstInContainer = false;
	} else {
		_, err = json.out.WriteRune(COMMA);
	}
	return
}
