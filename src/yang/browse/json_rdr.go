package browse

import (
	"io"
	"yang"
	"encoding/json"
	"fmt"
)

type JsonReader struct {
	in  io.Reader
}

func NewJsonReader(in io.Reader) *JsonReader {
	return &JsonReader{in:in}
}

func (self *JsonReader) GetSelector(meta yang.MetaList) (s *Selection, err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err = d.Decode(&values); err != nil {
		return
	}
	if yang.IsList(meta) {
		list := values[meta.GetIdent()]
		s, err = enterJson(nil, list.([]interface{}))
	} else {
		s, err = enterJson(values, nil)
	}
	s.Meta = meta
	return
}

func readLeafOrLeafList(meta yang.Meta, data interface{}, val *Value) (err error) {
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
		switch tmeta.DataType.Resolve().Ident {
		case "int32":
			a := data.([]float64)
			val.Intlist = make([]int, len(a))
			for i, f := range a {
				val.Intlist[i] = int(f)
			}
		case "string":
			a := data.([]string)
			val.Strlist = a
		case "boolean":
			a := data.([]string)
			val.Boollist = make([]bool, len(a))
			for i, s := range a {
				val.Boollist[i] = ("true" == s)
			}
		}
		val.Strlist = data.([]string)
	}
	return
}

func enterJson(values map[string]interface{}, list []interface{}) (s *Selection, err error) {
	s = &Selection{}
	var value interface{}
	var container = values
	var i int
	s.Enter = func() (child *Selection, e error) {
		value, s.Found = container[s.Position.GetIdent()]
fmt.Println("json_rdr: found container", s.Position.GetIdent(), s.Found)
		if s.Found {
			if yang.IsList(s.Position) {
				return enterJson(nil, value.([]interface{}))
			} else {
				return enterJson(value.(map[string]interface{}), nil)
			}
		}
		return
	}
	s.ReadValue = func (val *Value) (err error) {
		value, s.Found = container[s.Position.GetIdent()]
		if s.Found {
			return readLeafOrLeafList(s.Position, value, val)
		}
		return
	}
	s.Iterate = func(keys []string, first bool) (hasMore bool, err error) {
		container = nil
		if len(keys) > 0 {
			if first {
				keyFields := s.Meta.(*yang.List).Keys
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if jsonKeyMatches(keyFields, candidate, keys) {
						container = candidate
						break
					}
				}
			}
		} else {
			if first {
				i = 0
			} else {
				i++
			}
			if i < len(list) {
				container = list[i].(map[string]interface{})
			}
		}
fmt.Println("json_rdr: found list", s.Meta.GetIdent(), container != nil, list)
		return container != nil, nil
	}
	return
}

func jsonKeyMatches(keyFields []string, candidate map[string]interface{}, target []string) bool {
	for i, field := range keyFields {
		if candidate[field] != target[i] {
			return false
		}
	}
	return true
}
