package comm

import (
	"io"
	"yang"
	"yang/browse"
	"encoding/json"
	"strconv"
)

type JsonTransmitter struct {
	in  io.Reader
}

func NewJsonTransmitter(in io.Reader) *JsonTransmitter {
	return &JsonTransmitter{in:in}
}

func (self *JsonTransmitter) GetSelector(meta yang.MetaList) (s *browse.Selection, err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err = d.Decode(&values); err != nil {
		return
	}
	s, err = selectJsonContainer(values)
	s.Meta = meta
	return
}

type JsonNode struct {
	s *browse.Selection
	data interface{}
	found bool
}

func readLeafOrLeafList(meta yang.Meta, data interface{}, val *browse.Value) (err error) {
	switch tmeta := meta.(type) {
	case *yang.Leaf:
		switch tmeta.DataType.Resolve().Ident {
		case "int32":
			val.Int = int(data.(float64))
		case "string":
			s := data.(string)
			val.Str = s
		case "boolean":
			s := data.(string)
			val.Bool = ("true" == s)
		}
	case *yang.LeafList:
		a := data.([]string)
		switch tmeta.DataType.Resolve().Ident {
		case "int32":
			val.Intlist = make([]int, len(a))
			for i, s := range a {
				if val.Intlist[i], err = strconv.Atoi(s); err != nil {
					return err
				}
			}
		case "string":
			val.Strlist = a
		case "boolean":
			val.Boollist = make([]bool, len(a))
			for i, s := range a {
				val.Boollist[i] = ("true" == s)
			}
		}
		val.Strlist = data.([]string)
	}
	return
}

func defaultSelector(meta yang.Meta, data interface{}) (s *browse.Selection, err error) {
	switch meta.(type) {
	case *yang.List:
		return selectJsonList(data.([]interface{}))
	case *yang.Container:
		return selectJsonContainer(data.(map[string]interface{}))
	}
	return
}

func selectJsonContainer(values map[string]interface{}) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	var found bool
	var data interface{}
	s.Selector = func(ident string) (child *browse.Selection, e error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		data, found = values[ident]
		if !found {
			s.Position = nil
		} else if !yang.IsLeaf(s.Position) {
			return defaultSelector(s.Position, data)
		}
		return
	}
	s.Reader = func (val *browse.Value) (err error) {
		if !found {
			return
		}
		return readLeafOrLeafList(s.Position, data, val)
	}
	return
}

func selectJsonList(list []interface{}) (s *browse.Selection, err error) {
	var i int
	s = &browse.Selection{}
	node := JsonNode{s:s}
	s.ListIterator = func(keys []string, first bool) (bool, error) {
		/* ignoring keys, cannot see the use case */
		if (first) {
			i = 0
		} else {
			i += 1
		}
		if i < len(list) {
			node.data = list[i]
			return true, nil
		}
		return false, nil
	}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return selectJsonContainer(node.data.(map[string]interface{}))
	}
	return
}
