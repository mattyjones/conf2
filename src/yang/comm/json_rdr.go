package comm

import (
	"io"
	"yang"
	"yang/browse"
	"encoding/json"
	"fmt"
)

type JsonTransmitter struct {
	in  io.Reader
}

func NewJsonTransmitter(in io.Reader) *JsonTransmitter {
	return &JsonTransmitter{in:in}
}

func (self *JsonTransmitter) GetSelector() browse.Selection {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	parseErr := d.Decode(&values)
	return func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		if parseErr != nil {
			return parseErr
		}
		switch op {
		case browse.READ_VALUE:
			return self.readContainer(meta.(yang.MetaList), v, values)
		default:
			return browse.NotImplemented(meta);
		}
		return
	}
}

func (self *JsonTransmitter) readContainer(meta yang.MetaList, v *browse.Visitor, values map[string]interface{}) error {
	if values == nil || len(values) == 0 {
		return nil
	}
	v.EnterContainer(meta)
	i := yang.NewMetaListIterator(meta, true)
	for i.HasNextMeta() {
		field := i.NextMeta()
		self.readLeaf(field, v, values[field.GetIdent()])
	}
	v.ExitContainer(meta)
	return nil
}

func (self *JsonTransmitter) readList(meta *yang.List, v *browse.Visitor, value interface{}) (err error) {
	if value == nil {
		return nil
	}
	v.EnterList(meta)
	if values, isArray := value.([]interface{}); isArray {
		for _, arrayItem := range values {
			v.EnterListItem(meta)
			// TODO: Doesn't resolve switch cases
			if arrayContainer, isContainer := arrayItem.(map[string]interface{}); isContainer {
fmt.Println(arrayContainer)
				i := yang.NewMetaListIterator(meta, true)
				for ident, container := range arrayContainer {
					listItemMeta := yang.FindByIdent(i, ident).(yang.MetaList)
					self.readContainer(listItemMeta, v, container.(map[string]interface{}))
				}
			} else {
				msg := fmt.Sprint("expected a container for array item", meta.GetIdent())
				return &commError{msg}
			}
			v.ExitListItem(meta)
		}
	} else {
		msg := fmt.Sprint("expected an array for array item", meta.GetIdent())
		return &commError{msg}
	}
	v.ExitList(meta)
	return nil
}

func (self *JsonTransmitter) readLeaf(meta yang.Meta, v *browse.Visitor, value interface{}) (err error) {
	if value == nil {
		return
	}
	switch m := meta.(type) {
	case *yang.List:
		self.readList(m, v, value)
	case *yang.Container:
		if err = self.readContainer(m, v, value.(map[string]interface{})); err != nil {
			return err
		}
	case *yang.Leaf:
		// TODO: Support PutIntLeaf by looking at type
		v.Val.SetString(m, value.(string))
		v.Send(m)
	case *yang.LeafList:
		// TODO: Support PutIntLeafList by looking at type
		v.Val.Strlist = value.([]string)
		v.Send(m)
	}
	return nil
}


