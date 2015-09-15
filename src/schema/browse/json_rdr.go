package browse

import (
	"io"
	"schema"
	"encoding/json"
	"fmt"
)

type JsonReader struct {
	in  io.Reader
}

func NewJsonReader(in io.Reader) *JsonReader {
	return &JsonReader{in:in}
}

func (self *JsonReader) GetSelector(meta schema.MetaList) (s Selection, err error) {
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
	s.WalkState().Meta = meta
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

func enterJson(values map[string]interface{}, list []interface{}) (Selection, error) {
	s := &MySelection{}
	var value interface{}
	var container = values
	var i int
	s.OnChoose = func(choice *schema.Choice) (m schema.Meta, err error) {
		// go thru each case and if there are any properties in the data that are not
		// part of the schema, that disqualifies that case and we move onto next case
		// until one case aligns with data.  If no cases align then input in inconclusive
		// i.e. non-discriminating and we should error out.
		cases := schema.NewMetaListIterator(choice, false)
		for cases.HasNextMeta() {
			kase := cases.NextMeta().(*schema.ChoiceCase)
			aligned := true
			props := schema.NewMetaListIterator(kase, true)
			for props.HasNextMeta() {
				prop := props.NextMeta()
				_, found := container[prop.GetIdent()]
				if !found {
					aligned = false
					break;
				} else {
					m = prop
				}
			}
			if aligned {
				return m, nil
			}
		}
		msg := fmt.Sprintf("No discriminating data for choice schema %s ", s.ToString())
		return nil, &browseError{Msg:msg}
	}
	s.OnSelect = func() (child Selection, e error) {
		state := s.WalkState()
		value, state.Found = container[s.State.Position.GetIdent()]
		if state.Found {
			if schema.IsList(s.State.Position) {
				return enterJson(nil, value.([]interface{}))
			} else {
				return enterJson(value.(map[string]interface{}), nil)
			}
		}
		return
	}
	s.OnRead = func (val *Value) (err error) {
		state := s.WalkState()
		value, state.Found = container[s.State.Position.GetIdent()]
		if state.Found {
			return readLeafOrLeafList(s.State.Position, value, val)
		}
		return
	}
	s.OnNext = func(key []Value, first bool) (hasMore bool, err error) {
		container = nil
		if len(key) > 0 {
			if first {
				keyFields := s.State.Meta.(*schema.List).Keys
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if jsonKeyMatches(keyFields, candidate, key) {
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
	return s, nil
}

func jsonKeyMatches(keyFields []string, candidate map[string]interface{}, key []Value) bool {
	for i, field := range keyFields {
		if candidate[field] != key[i].String() {
			return false
		}
	}
	return true
}
