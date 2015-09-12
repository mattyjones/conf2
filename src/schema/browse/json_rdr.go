package browse

import (
	"io"
	"schema"
	"encoding/json"
)

type JsonReader struct {
	in  io.Reader
}

func NewJsonReader(in io.Reader) *JsonReader {
	return &JsonReader{in:in}
}

func (self *JsonReader) GetSelector(meta schema.MetaList) (s *Selection, err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err = d.Decode(&values); err != nil {
		return
	}
	if schema.IsList(meta) {
		list := values[meta.GetIdent()]
		s, err = enterJson(nil, list.([]interface{}))
	} else {
		s, err = enterJson(values, nil)
	}
	s.Meta = meta
	return
}

func readLeafOrLeafList(meta schema.Meta, data interface{}, val *Value) (err error) {
	switch tmeta := meta.(type) {
	case *schema.Leaf:
		switch tmeta.DataType.Format {
		case schema.FMT_INT32:
			val.Int = int(data.(float64))
		case schema.FMT_STRING:
			s := data.(string)
			val.Str = s
		case schema.FMT_BOOLEAN:
			s := data.(string)
			val.Bool = ("true" == s)
		}
	case *schema.LeafList:
		switch tmeta.DataType.Format {
		case schema.FMT_INT32:
			a := data.([]float64)
			val.Intlist = make([]int, len(a))
			for i, f := range a {
				val.Intlist[i] = int(f)
			}
		case schema.FMT_STRING:
			a := data.([]string)
			val.Strlist = a
		case schema.FMT_BOOLEAN:
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
		if s.Found {
			if schema.IsList(s.Position) {
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
				keyFields := s.Meta.(*schema.List).Keys
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
