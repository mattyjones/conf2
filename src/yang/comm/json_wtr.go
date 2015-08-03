package comm

import (
	"io"
	"yang"
	"yang/browse"
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

type JsonReceiver struct {
	out *bufio.Writer
	firstInContainer bool
	firstInDoc bool
}

func NewJsonReceiver(out io.Writer ) *JsonReceiver {
	return &JsonReceiver{out:bufio.NewWriter(out), firstInDoc:true, firstInContainer: true}
}

func (json *JsonReceiver) Flush() error {
	return json.out.Flush()
}

func (json *JsonReceiver) EnterList(meta *yang.List) (err error) {
	if err = json.writeIdent(meta.GetIdent()); err == nil {
		_, err = json.out.WriteRune(OPEN_ARRAY)
		json.firstInContainer = true;
	}
	return
}

func (json *JsonReceiver) ExitList(meta *yang.List) (err error) {
	_, err = json.out.WriteRune(CLOSE_ARRAY);
	json.firstInContainer = false
	return
}

func (json *JsonReceiver) EnterContainer(meta yang.MetaList) error {
	return json.beginObject(meta.GetIdent())
}
func (json *JsonReceiver) ExitContainer(meta yang.MetaList) (err error) {
	_, err = json.out.WriteRune(CLOSE_OBJ)
	return
}

func (json *JsonReceiver) UpdateValue(meta yang.Meta, v *browse.Value) (err error) {
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

func (json *JsonReceiver) EnterListItem(meta *yang.List) error {
	return json.beginArrayItem()
}
func (json *JsonReceiver) ExitListItem(meta *yang.List) error {
	return json.endArrayItem()
}

func (json *JsonReceiver) writeBool(b bool) error {
	if b {
		return json.writeString("true")
	} else {
		return json.writeString("false")
	}
}

func (json *JsonReceiver) writeInt(i int) (err error) {
	_, err = json.out.WriteString(strconv.Itoa(i))
	return
}

func (json *JsonReceiver) writeString(s string) (err error) {
	if _, err = json.out.WriteRune(QUOTE); err == nil {
		if _, err = json.out.WriteString(s); err == nil {
			_, err = json.out.WriteRune(QUOTE);
		}
	}
	return
}

func (json *JsonReceiver) beginArrayItem() (err error) {
	_, err = json.out.WriteRune(OPEN_OBJ)
	json.firstInContainer = true
	return
}
func (json *JsonReceiver) endArrayItem() (err error) {
	_, err = json.out.WriteRune(CLOSE_OBJ);
	return
}
// helper functions
func (json *JsonReceiver) beginObject(ident string) (err error) {
	if json.firstInDoc {
		json.firstInDoc = false
	} else {
		err = json.writeIdent(ident);
	}
	if err == nil {
		_, err = json.out.WriteRune(OPEN_OBJ);
		json.firstInContainer = true;
	}
	return
}
func (json *JsonReceiver) writeIdent(ident string) (err error) {
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
func (json *JsonReceiver) writeDelim() (err error) {
	if json.firstInContainer {
		json.firstInContainer = false;
	} else {
		_, err = json.out.WriteRune(COMMA);
	}
	return
}

