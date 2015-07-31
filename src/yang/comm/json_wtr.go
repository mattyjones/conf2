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
	return &JsonReceiver{out:bufio.NewWriter(out), firstInDoc:true}
}

func (json *JsonReceiver) Flush() error {
	return json.out.Flush()
}

func (self *JsonReceiver) GetSelector() browse.Selection {
	return func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		switch op {
			case browse.CREATE_LIST:
				self.writeIdent(meta.GetIdent());
				self.out.WriteRune(OPEN_ARRAY)
				self.firstInContainer = true;
			case browse.CREATE_LIST_ITEM:
				self.beginArrayItem()
			case browse.CREATE_CHILD:
				self.beginObject(meta.GetIdent())
			case browse.POST_CREATE_LIST_ITEM:
				self.endArrayItem()
			case browse.POST_CREATE_LIST:
				self.out.WriteRune(CLOSE_ARRAY);
				self.firstInContainer = false
			case browse.POST_CREATE_CHILD:
				self.out.WriteRune(CLOSE_OBJ);
			case browse.UPDATE_VALUE:
				self.writeIdent(meta.GetIdent());
				switch tmeta := meta.(type) {
				case *yang.Leaf:
					switch tmeta.GetDataType().Resolve().Ident {
					case "boolean":
						self.writeBool(v.Val.Bool)
					case "int32":
						self.writeInt(v.Val.Int)
					case "string":
						self.writeString(v.Val.Str)
					}
				case *yang.LeafList:
					switch tmeta.GetDataType().Resolve().Ident {
					case "int32":
						self.out.WriteRune(OPEN_ARRAY)
						for i, n := range v.Val.Intlist {
							if i > 0 {
								self.out.WriteRune(COMMA)
							}
							self.writeInt(n)
						}
						self.out.WriteRune(CLOSE_ARRAY)
					case "string":
						self.out.WriteRune(OPEN_ARRAY)
						for i, s := range v.Val.Strlist {
							if i > 0 {
								self.out.WriteRune(COMMA)
							}
							self.writeString(s)
						}
						self.out.WriteRune(CLOSE_ARRAY)
					case "boolean":
						self.out.WriteRune(OPEN_ARRAY)
						for i, b := range v.Val.Boollist {
							if i > 0 {
								self.out.WriteRune(COMMA)
							}
							self.writeBool(b)
						}
						self.out.WriteRune(CLOSE_ARRAY)
					}
				}
		default:
			return browse.NotImplemented(meta)
		}
		return
	}
}

func (self *JsonReceiver) writeBool(b bool) {
	if b {
		self.writeString("true")
	} else {
		self.writeString("false")
	}
}

func (self *JsonReceiver) writeInt(i int) {
	self.out.WriteString(strconv.Itoa(i))
}

func (self *JsonReceiver) writeString(s string) {
	self.out.WriteRune(QUOTE);
	self.out.WriteString(s)
	self.out.WriteRune(QUOTE);
}

func (self *JsonReceiver) beginArrayItem() {
	self.out.WriteRune(OPEN_OBJ)
	self.firstInContainer = true
}
func (self *JsonReceiver) endArrayItem() {
	self.out.WriteRune(CLOSE_OBJ);
}
// helper functions
func (self *JsonReceiver) beginObject(ident string) {
	if self.firstInDoc {
		self.firstInDoc = false
	} else {
		self.writeIdent(ident);
	}
	self.out.WriteRune(OPEN_OBJ);
	self.firstInContainer = true;
}
func (self *JsonReceiver) writeIdent(ident string) {
	self.writeDelim()
	self.out.WriteRune(QUOTE);
	self.out.WriteString(ident);
	self.out.WriteRune(QUOTE);
	self.out.WriteRune(COLON);
}
func (self *JsonReceiver) writeDelim() {
	if self.firstInContainer {
		self.firstInContainer = false;
	} else {
		self.out.WriteRune(COMMA);
	}
}

