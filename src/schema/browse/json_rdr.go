package browse

import (
	"io"
	"schema"
	"encoding/json"
	"fmt"
	"errors"
)

type JsonReader struct {
	In io.Reader
	Meta *schema.Module
	values map[string]interface{}
}

func NewJsonReader(in io.Reader, module *schema.Module) *JsonReader {
	r := &JsonReader{In:in, Meta:module}
	return r
}

func NewJsonFragmentReader(in io.Reader) *JsonReader {
	r := &JsonReader{In:in}
	return r
}

func (self *JsonReader) Module() *schema.Module {
	return self.Meta
}

func (self *JsonReader) Selector(path *Path, strategy Strategy) (s Selection,  state *WalkState, err error) {
	if strategy != READ {
		return nil, nil, errors.New("Only read is supported")
	}
	if self.values == nil {
		d := json.NewDecoder(self.In)
		if err = d.Decode(&self.values); err != nil {
			return nil, nil, err
		}
	}
	if self.Meta == nil {
		return self.fragmentSelector(path, strategy)
	}
	s, err = self.enterJson(self.values, nil, false)
	return WalkPath(NewWalkState(self.Meta), s, path)
}

func (self *JsonReader) fragmentSelector(path *Path, strategy Strategy) (s Selection,  state *WalkState, err error) {
	// try to determine if we're reading a list
	if (len(self.values) == 1) {
		for _, value := range self.values {
			if list, isList := value.([]interface{}); isList {
				var insideList bool
				lastSegment := path.LastSegment()
				if lastSegment != nil {
					insideList = len(lastSegment.Keys) > 0
				}
				s, err = self.enterJson(self.values, list, insideList)
			}
		}
	}
	if s == nil {
		s, err = self.enterJson(self.values, nil, false)
	}
	return
}

// need to pass in walk state so it knows what to look for in initial data
// particularly when in the middle of a list
//func (self *JsonReader) GetSelector(state *WalkState) (s Selection, err error) {
//	var values map[string]interface{}
//	d := json.NewDecoder(self.in)
//	if err = d.Decode(&values); err != nil {
//		return
//	}
//	if schema.IsList(state.SelectedMeta()) {
//		if state.InsideList() {
////			singletonList := []interface{} {
////				values,
////			}
//			s, err = self.enterJson(values, nil, true)
//		} else {
//			list, found := values[state.SelectedMeta().GetIdent()]
//			if !found {
//				msg := fmt.Sprintf("Could not find json data %s", state.String())
//				return nil, &browseError{Msg:msg}
//			}
//			s, err = self.enterJson(nil, list.([]interface{}), false)
//		}
//	} else {
//		s, err = self.enterJson(values, nil, false)
//	}
//	return
//}

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
	s.OnChoose = func(state *WalkState, choice *schema.Choice) (m schema.Meta, err error) {
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
		msg := fmt.Sprintf("No discriminating data for choice schema %s ", state.String())
		return nil, &browseError{Msg:msg}
	}
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (child Selection, e error) {
fmt.Printf("json-rdr - OnSelect %s\n", state.String())
		var found bool
		if value, found = container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
fmt.Printf("json-rdr - OnSelect %s FOUND list, len=%d\n", state.String(), len(value.([]interface{})))
				return self.enterJson(nil, value.([]interface{}), false)
			} else {
fmt.Printf("json-rdr - OnSelect %s FOUND container\n", state.String())
				return self.enterJson(value.(map[string]interface{}), nil, false)
			}
		}
fmt.Printf("json-rdr - OnSelect %s NOT found\n", state.String())
		return
	}
	s.OnRead = func (state *WalkState, meta schema.HasDataType) (val *Value, err error) {
		var found bool
		if value, found = container[meta.GetIdent()]; found {
			return self.readLeafOrLeafList(meta, value)
		}
		return
	}
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, first bool) (hasMore bool, err error) {
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
				keyFields := meta.Keys
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
