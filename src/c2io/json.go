package c2io
import (
	"io"
	"yang"
	"encoding/json"
	"bufio"
	"fmt"
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
	first bool
}

// Receiver
func NewJsonReceiver(out io.Writer ) *JsonReceiver {
	return &JsonReceiver{out:bufio.NewWriter(out)}
}

func (self *JsonReceiver) StartTransaction() {
	self.out.WriteRune(OPEN_OBJ)
	self.first = true
}
func (self *JsonReceiver) NewListItem(y yang.Meta) {
	self.writeDelim()
	self.out.WriteRune(OPEN_OBJ)
}
func (self *JsonReceiver) NewObject(y yang.MetaList) {
	self.beginObject(y.GetIdent())
}
func (self *JsonReceiver) NewList(y *yang.List) {
	self.writeIdent(y.GetIdent());
	self.out.WriteRune(OPEN_ARRAY)
	self.first = true;
}
func (self *JsonReceiver) PutIntLeaf(y *yang.Leaf, value int) {
	self.writeDelim()
	self.writeIdent(y.GetIdent());
	self.out.WriteString(string(value))
}
func (self *JsonReceiver) PutStringLeaf(y *yang.Leaf, value string) {
	self.writeIdent(y.GetIdent());
	self.out.WriteRune(QUOTE);
	self.out.WriteString(value)
	self.out.WriteRune(QUOTE);
}
func (self *JsonReceiver) PutIntLeafList(y *yang.LeafList, ints []int) {
	self.writeIdent(y.GetIdent());
	self.out.WriteRune(OPEN_ARRAY)
	for i, n := range ints {
		if i > 0 {
			self.out.WriteRune(COMMA)
		}
		self.out.WriteString(string(n))
	}
	self.out.WriteRune(CLOSE_ARRAY)
}
func (self *JsonReceiver) PutStringLeafList(y *yang.LeafList, strings []string) {
	self.writeIdent(y.GetIdent());
	self.out.WriteRune(OPEN_ARRAY)
	for i, n := range strings {
		if i > 0 {
			self.out.WriteRune(COMMA)
		}
		self.out.WriteRune(QUOTE);
		self.out.WriteString(n)
		self.out.WriteRune(QUOTE);
	}
	self.out.WriteRune(CLOSE_ARRAY)
}
func (self *JsonReceiver) ExitObject(y yang.MetaList) {
	self.out.WriteRune(CLOSE_OBJ);
}
func (self *JsonReceiver) ExitListItem(y yang.Meta) {
	self.out.WriteRune(CLOSE_OBJ);
}
func (self *JsonReceiver) ExitList(y *yang.List) {
	self.out.WriteRune(CLOSE_ARRAY);
	self.first = false
}
func (self *JsonReceiver) EndTransaction() {
	self.out.WriteRune(CLOSE_OBJ);
	self.out.Flush()
}
// helper functions
func (self *JsonReceiver) beginObject(ident string) {
	self.writeIdent(ident);
	self.out.WriteRune(OPEN_OBJ);
	self.first = true;
}
func (self *JsonReceiver) writeIdent(ident string) {
	self.writeDelim()
	self.out.WriteRune(QUOTE);
	self.out.WriteString(ident);
	self.out.WriteRune(QUOTE);
	self.out.WriteRune(COLON);
}
func (self *JsonReceiver) writeDelim() {
	if self.first {
		self.first = false;
	} else {
		self.out.WriteRune(COMMA);
	}
}
type JsonTransmitter struct {
	in  io.Reader
	metaRoot yang.MetaList
	out Receiver
}
func (self *JsonTransmitter) Transmit() (err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err := d.Decode(&values); err == nil {
		self.out.StartTransaction()
		if err = self.ReadValues(self.metaRoot, values); err == nil {
			self.out.EndTransaction()
		}
	}
	return
}
func (self *JsonTransmitter) ReadValues(meta yang.MetaList, values map[string]interface{}) error {
	if values == nil || len(values) == 0 {
		return nil
	}
	i := yang.NewMetaListIterator(meta, true)
	for i.HasNextMeta() {
		field := i.NextMeta()
		self.ReadValue(field, values[field.GetIdent()])
	}
	return nil
}

func (self *JsonTransmitter) ReadListValues(meta *yang.List, value interface{}) (err error) {
	if value == nil {
		return nil
	}
	self.out.NewList(meta)
	if values, isArray := value.([]interface{}); isArray {
		for _, arrayItem := range values {
			// TODO: Doesn't resolve switch cases
			self.out.NewListItem(meta)
			if arrayContainer, isContainer := arrayItem.(map[string]interface{}); isContainer {
				self.ReadValues(meta, arrayContainer)
			} else {
				msg := fmt.Sprint("expected a container for array item", meta.GetIdent())
				return &yang.MetaError{msg}
			}
			self.out.ExitListItem(meta)
		}
	} else {
		msg := fmt.Sprint("expected an array for array item", meta.GetIdent())
		return &yang.MetaError{msg}
	}
	self.out.ExitList(meta)
	return nil
}

func (self *JsonTransmitter) ReadValue(meta yang.Meta, value interface{}) (err error) {
	if value == nil {
		return
	}
	switch m := meta.(type) {
		case *yang.List:
			self.ReadListValues(m, value)
		case *yang.Container:
			self.out.NewObject(m)
			if err = self.ReadValues(m, value.(map[string]interface{})); err != nil {
				return err
			}
			self.out.ExitObject(m)
		case *yang.Leaf:
			// TODO: Support PutIntLeaf by looking at type
			self.out.PutStringLeaf(m, value.(string))
		case *yang.LeafList:
			// TODO: Support PutIntLeafList by looking at type
			self.out.PutStringLeafList(m, value.([]string))
	}
	return nil
}


