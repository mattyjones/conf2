package browse

import (
	"io"
	"schema"
	"encoding/json"
	"fmt"
	"errors"
)

type JsonReader struct {
	in  io.Reader
}

func NewJsonReader(in io.Reader) *JsonReader {
	return &JsonReader{in:in}
}

// need to pass in walk state so it knows what to look for in initial data
// particularly when in the middle of a list
func (self *JsonReader) GetSelector(meta schema.MetaList, insideList bool) (s Selection, err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err = d.Decode(&values); err != nil {
		return
	}
	if schema.IsList(meta) {
		if insideList {
			singletonList := []interface{} {
				values,
			}
			s, err = self.enterJson(nil, singletonList, insideList)
		} else {
			list, found := values[meta.GetIdent()]
			if !found {
				msg := fmt.Sprintf("Could not find json data %s", meta.GetIdent())
				return nil, &browseError{Msg:msg}
			}
			s, err = self.enterJson(nil, list.([]interface{}), false)
		}
	} else {
		s, err = self.enterJson(values, nil, false)
	}
	s.WalkState().Meta = meta
	return
}

func (self *JsonReader) readLeafOrLeafList(meta schema.HasDataType, data interface{}) (v *Value, err error) {
	v = &Value{}
	switch meta.GetDataType().Format {
	case schema.FMT_INT32:
		v.Int = int(data.(float64))
	case schema.FMT_INT32_LIST:
		a := data.([]float64)
		v.Intlist = make([]int, len(a))
		for i, f := range a {
			v.Intlist[i] = int(f)
		}
	case schema.FMT_STRING:
		v.Str = data.(string)
	case schema.FMT_STRING_LIST:
		a := data.([]string)
		v.Strlist = a
	case schema.FMT_BOOLEAN:
		s := data.(string)
		v.Bool = ("true" == s)
	case schema.FMT_BOOLEAN_LIST:
		a := data.([]string)
		v.Boollist = make([]bool, len(a))
		for i, s := range a {
			v.Boollist[i] = ("true" == s)
		}
	default:
		return nil, errors.New("Not implemented")
	}
	return
}

func (self *JsonReader) enterJson(values map[string]interface{}, list []interface{}, insideList bool) (Selection, error) {
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
	s.OnSelect = func(meta schema.MetaList) (child Selection, e error) {
		var found bool
		if value, found = container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return self.enterJson(nil, value.([]interface{}), false)
			} else {
				return self.enterJson(value.(map[string]interface{}), nil, false)
			}
		}
		return
	}
	s.OnRead = func (meta schema.HasDataType) (val *Value, err error) {
		var found bool
		if value, found = container[s.State.Position.GetIdent()]; found {
			return self.readLeafOrLeafList(meta, value)
		}
		return
	}
	s.OnNext = func(key []*Value, first bool) (hasMore bool, err error) {
		container = nil
		if len(key) > 0 {
			if insideList {
				if first && len(list) > 0 {
					container = list[0].(map[string]interface{})
					return true, nil
				} else {
					return false, nil
				}
			} else if first {
				keyFields := s.State.Meta.(*schema.List).Keys
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if self.jsonKeyMatches(keyFields, candidate, key) {
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

func (self *JsonReader) jsonKeyMatches(keyFields []string, candidate map[string]interface{}, key []*Value) bool {
	for i, field := range keyFields {
		if candidate[field] != key[i].String() {
			return false
		}
	}
	return true
}
