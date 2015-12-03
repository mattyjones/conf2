package data

import (
	"bufio"
	"io"
	"schema"
	"strconv"
)

const QUOTE = '"'

type JsonWriter struct {
	out                *bufio.Writer
	firstInContainer   bool
	startingInsideList bool
	firstWrite         bool
	closeArrayOnExit   bool
}

func NewJsonWriter(out io.Writer) *JsonWriter {
	return &JsonWriter{
		out:              bufio.NewWriter(out),
		firstInContainer: true,
	}
}

func (json *JsonWriter) Node() Node {
	s := &MyNode{Label: "JSON Write"}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (child Node, err error) {
		if ! new {
			return nil, nil
		}
		if schema.IsList(meta) {
			err = json.beginList(meta.GetIdent())
		} else {
			err = json.beginContainer(meta.GetIdent())
		}
		return s, nil
	}
	s.OnEvent = func(sel *Selection, e Event) (err error) {
		switch e {
		case BEGIN_EDIT:
			_, err = json.out.WriteRune('{')
			json.startingInsideList = schema.IsList(sel.State.SelectedMeta())
			json.firstWrite = true
			return err
		case END_EDIT:
			if err = json.conditionallyCloseArrayOnLastWrite(); err != nil {
				return err
			}
			if _, err = json.out.WriteRune('}'); err != nil {
				return err
			}
			err = json.out.Flush()
		case LEAVE:
			if schema.IsList(sel.State.SelectedMeta()) {
				err = json.endList()
			} else {
				err = json.endContainer()
			}
		case LEAVE_ITEM:
			err = json.endArrayItem()
		}
		return
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, v *schema.Value) (err error) {
		err = json.writeValue(meta, v)
		json.firstWrite = false
		return
	}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (next Node, err error) {
		if ! new {
			return nil, nil
		}
		if err = json.conditionallyOpenArrayOnFirstWrite(meta.GetIdent()); err != nil {
			return nil, err
		}
//		if ! first {
//			if err = json.writeDelim(); err != nil {
//				return
//			}
//		}
//		if err = json.beginObject(); err != nil {
//			return
//		}
		return s, json.beginArrayItem()
	}
	return s
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
		_, err = json.out.WriteRune('[')
		json.firstInContainer = true
	}
	return
}

func (json *JsonWriter) endList() (err error) {
	_, err = json.out.WriteRune(']')
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
	json.firstInContainer = false
	_, err = json.out.WriteRune('}')
	return
}

func (json *JsonWriter) writeValue(meta schema.Meta, v *schema.Value) (err error) {
	json.writeIdent(meta.GetIdent())
	if schema.IsListFormat(v.Type.Format) {
		if _, err = json.out.WriteRune('['); err != nil {
			return
		}
	}
	switch v.Type.Format {
	case schema.FMT_BOOLEAN:
		err = json.writeBool(v.Bool)
	case schema.FMT_INT64:
		err = json.writeInt64(v.Int64)
	case schema.FMT_INT32:
		err = json.writeInt(v.Int)
	case schema.FMT_STRING, schema.FMT_ENUMERATION:
		err = json.writeString(v.Str)
	case schema.FMT_BOOLEAN_LIST:
		for i, b := range v.Boollist {
			if i > 0 {
				if _, err = json.out.WriteRune(','); err != nil {
					return
				}
			}
			if err = json.writeBool(b); err != nil {
				return
			}
		}
		_, err = json.out.WriteRune(']')
	case schema.FMT_INT32_LIST:
		for i, n := range v.Intlist {
			if i > 0 {
				if _, err = json.out.WriteRune(','); err != nil {
					return
				}
			}
			if err = json.writeInt(n); err != nil {
				return
			}
		}
	case schema.FMT_INT64_LIST:
		for i, n := range v.Int64list {
			if i > 0 {
				if _, err = json.out.WriteRune(','); err != nil {
					return
				}
			}
			if err = json.writeInt64(n); err != nil {
				return
			}
		}
	case schema.FMT_STRING_LIST, schema.FMT_ENUMERATION_LIST:
		for i, s := range v.Strlist {
			if i > 0 {
				if _, err = json.out.WriteRune(','); err != nil {
					return
				}
			}
			if err = json.writeString(s); err != nil {
				return
			}
		}
	}
	if schema.IsListFormat(v.Type.Format) {
		if _, err = json.out.WriteRune(']'); err != nil {
			return
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

func (json *JsonWriter) writeInt64(i int64) (err error) {
	_, err = json.out.WriteString(strconv.FormatInt(i, 10))
	return
}

func (json *JsonWriter) writeString(s string) (err error) {
	if _, err = json.out.WriteRune(QUOTE); err == nil {
		if _, err = json.out.WriteString(s); err == nil {
			_, err = json.out.WriteRune(QUOTE)
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
	json.firstInContainer = false
	_, err = json.out.WriteRune('}')
	return
}

func (json *JsonWriter) beginObject() (err error) {
	if err == nil {
		_, err = json.out.WriteRune('{')
		json.firstInContainer = true
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
	_, err = json.out.WriteRune(':')
	return
}

func (json *JsonWriter) writeDelim() (err error) {
	if json.firstInContainer {
		json.firstInContainer = false
	} else {
		_, err = json.out.WriteRune(',')
	}
	return
}
