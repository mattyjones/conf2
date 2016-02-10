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
}

type closerFunc func() error

func NewJsonWriter(out io.Writer) *JsonWriter {
	return &JsonWriter{
		out:              bufio.NewWriter(out),
	}
}

func (json *JsonWriter) Node() Node {
	var closer closerFunc
	// JSON can begin at a container, inside a list or inside a container, each of these has
	// different results to make json legal
	return &Extend{
		Label: "JSON",
		Node: json.Container(json.endContainer),
		OnSelect: func(p Node, sel *Selection, meta schema.MetaList, new bool) (child Node, err error) {
			if closer == nil {
				json.beginObject()
				closer = json.endContainer
			}
			return p.Select(sel, meta, new)
		},
		OnNext: func(p Node, sel *Selection, meta *schema.List, new bool, keys []*Value, first bool) (next Node, err error) {
			if closer == nil {
				json.beginObject()
				json.beginList(meta.GetIdent())
				closer = func() (closeErr error) {
					if closeErr = json.endList(); closeErr == nil {
						closeErr = json.endContainer()
					}
					return closeErr
				}
			}
			return p.Next(sel, meta, new, keys, first)
		},
		OnWrite: func(p Node, sel *Selection, meta schema.HasDataType, v *Value) (err error) {
			if closer == nil {
				json.beginObject()
				closer = json.endContainer
			}
			return p.Write(sel, meta, v)
		},
		OnEvent:func(p Node, s *Selection, e Event) error {
			var err error
			switch e {
			case END_TREE_EDIT:
				if closer != nil {
					if err = closer(); err != nil {
						return err
					}
				}
				err = json.out.Flush()
			default:
				err = p.Event(s, e)
			}
			return err
		},
	}
}

func (json *JsonWriter) Container(closer closerFunc) Node {
	first := true
	delim := func() (err error) {
		if ! first {
			if _, err = json.out.WriteRune(','); err != nil {
				return
			}
		} else {
			first = false
		}
		return
	}
	s := &MyNode{Label: "JSON Write"}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (child Node, err error) {
		if ! new {
			return nil, nil
		}
		if err = delim(); err != nil {
			return nil, err
		}
		if schema.IsList(meta) {
			if err = json.beginList(meta.GetIdent()); err != nil {
				return nil, err
			}
			return json.Container(json.endList), nil

		}
		if err = json.beginContainer(meta.GetIdent()); err != nil {
			return nil, err
		}
		return json.Container(json.endContainer), nil
	}
	s.OnEvent = func(sel *Selection, e Event) (err error) {
		switch e {
		case LEAVE:
			err = closer()
		}
		return
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, v *Value) (err error) {
		if err = delim(); err != nil {
			return err
		}
		err = json.writeValue(meta, v)
		return
	}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, keys []*Value, first bool) (next Node, err error) {
		if ! new {
			return nil, nil
		}
		if err = delim(); err != nil {
			return nil, err
		}
		if err = json.beginObject(); err != nil {
			return nil, err
		}
		return json.Container(json.endContainer), nil
	}
	return s
}

func (json *JsonWriter) beginList(ident string) (err error) {
	if err = json.writeIdent(ident); err == nil {
		_, err = json.out.WriteRune('[')
	}
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

func (json *JsonWriter) beginObject() (err error) {
	if err == nil {
		_, err = json.out.WriteRune('{')
	}
	return
}

func (json *JsonWriter) writeIdent(ident string) (err error) {
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

func (json *JsonWriter) endList() (err error) {
	_, err = json.out.WriteRune(']')
	return
}

func (json *JsonWriter) endContainer() (err error) {
	_, err = json.out.WriteRune('}')
	return
}

func (json *JsonWriter) writeValue(meta schema.Meta, v *Value) (err error) {
	json.writeIdent(meta.GetIdent())
	if schema.IsListFormat(v.Type.Format()) {
		if _, err = json.out.WriteRune('['); err != nil {
			return
		}
	}
	switch v.Type.Format() {
	case schema.FMT_BOOLEAN:
		err = json.writeBool(v.Bool)
	case schema.FMT_ANYDATA:
		var s string
		s, err = v.Data.String()
		if err == nil {
			// TODO: don't assume output is json
			json.out.WriteString(s)
		}
	case schema.FMT_INT64:
		err = json.writeInt64(v.Int64)
	case schema.FMT_INT32:
		err = json.writeInt(v.Int)
	case schema.FMT_DECIMAL64:
		err = json.writeFloat(v.Float)
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
	if schema.IsListFormat(v.Type.Format()) {
		if _, err = json.out.WriteRune(']'); err != nil {
			return
		}
	}
	return
}

func (json *JsonWriter) writeBool(b bool) (err error) {
	if b {
		_, err = json.out.WriteString("true")
	} else {
		_, err = json.out.WriteString("false")
	}
	return
}

func (json *JsonWriter) writeFloat(f float64) (err error) {
	_, err = json.out.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
	return
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

