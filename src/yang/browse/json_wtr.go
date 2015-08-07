package browse

import (
	"io"
	"yang"
	"bufio"
	"strconv"
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
	firstInContainer bool
}

func NewJsonWriter(out io.Writer) *JsonWriter {
	return &JsonWriter{
		out:bufio.NewWriter(out),
		firstInContainer: true,
	}
}

func (json *JsonWriter) GetSelector() (*Selection, error) {
	return json.selectJson(false)
}

func (json *JsonWriter) selectJson(isList bool) (*Selection, error) {
	s := &Selection{}
	s.Select = func(ident string) (child *Selection, err error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		if yang.IsList(s.Position) {
			return json.selectJson(true)
		} else if yang.IsContainer(s.Position) {
			return json.selectJson(false)
		}
		return nil, nil
	}
	s.Edit = func(op Operation, v *Value) (err error) {
		switch op {
		case BEGIN_EDIT:
			_, err = json.out.WriteRune(OPEN_OBJ)
		case END_EDIT:
			if _, err = json.out.WriteRune(CLOSE_OBJ); err != nil {
				return err
			}
			return json.out.Flush()
		case CREATE_CHILD:
			if isList {
				err = json.beginArrayItem()
			}
			err = json.beginContainer(s.Position)
		case POST_CREATE_CHILD:
			err = json.endContainer()
			if isList {
				err = json.endArrayItem()
			}
		case CREATE_LIST:
			return json.beginList(s.Position)
		case POST_CREATE_LIST:
			return json.endList()
		case UPDATE_VALUE:
			return json.writeValue(s.Position, v)
		default:
			return &browseError{Msg:"Operation not supported"}
		}
		return
	}
	if isList {
		s.Iterate = func(keys []string, hasMore bool) (bool, error) {
			// never finds values because always INSERT mode
			return false, nil
		}
	}
	return s, nil
}

func (json *JsonWriter) beginList(meta yang.Meta) (err error) {
	if err = json.writeIdent(meta.GetIdent()); err == nil {
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

func (json *JsonWriter) beginContainer(meta yang.Meta) (err error) {
	if err = json.writeIdent(meta.GetIdent()); err != nil {
		return
	}
	if err = json.beginObject(); err != nil {
		return
	}
	return
}
func (json *JsonWriter) endContainer() (err error) {
	_, err = json.out.WriteRune(CLOSE_OBJ)
	return
}

func (json *JsonWriter) writeValue(meta yang.Meta, v *Value) (err error) {
	json.writeIdent(meta.GetIdent());
	switch tmeta := meta.(type) {
	case *yang.Leaf:
		switch tmeta.GetDataType().Resolve().Ident {
		case "boolean":
			err = json.writeBool(v.Bool)
		case "int32":
			err = json.writeInt(v.Int)
		case "string":
			err = json.writeString(v.Str)
		}
	case *yang.LeafList:
		switch tmeta.GetDataType().Resolve().Ident {
		case "int32":
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
		case "string":
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
		case "boolean":
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
		}
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
