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
	s, err = self.enterJson(self.values)
	return WalkPath(NewWalkState(self.Meta), s, path)
}

func (self *JsonReader) fragmentSelector(path *Path, strategy Strategy) (s Selection,  state *WalkState, err error) {
	// try to determine if we're reading a list
	if (len(self.values) == 1) {
		// only 1 item - iterates only once
		for _, value := range self.values {
			if list, isList := value.([]interface{}); isList {
				lastSegment := path.LastSegment()
				if lastSegment != nil {
					if len(lastSegment.Keys) == 0 {
						s, err = self.enterJsonList(list)
					}
				}
			}
		}
	}
	if s == nil {
		s, err = self.enterJson(self.values)
	}
	return
}

func (self *JsonReader) readLeafOrLeafList(meta schema.HasDataType, data interface{}) (v *Value, err error) {
	// TODO: Consider using CoerseValue
	v = &Value{}
	switch meta.GetDataType().Format {
	case schema.FMT_INT32:
		v.Int = int(data.(float64))
	case schema.FMT_INT32_LIST:
		a := data.([]interface{})
		v.Intlist = make([]int, len(a))
		for i, f := range a {
			v.Intlist[i] = int(f.(float64))
		}
	case schema.FMT_STRING:
		v.Str = data.(string)
	case schema.FMT_STRING_LIST:
		v.Strlist = asStringArray(data.([]interface{}))
	case schema.FMT_BOOLEAN:
		s := data.(string)
		v.Bool = ("true" == s)
	case schema.FMT_BOOLEAN_LIST:
		a := data.([]interface{})
		v.Boollist = make([]bool, len(a))
		for i, s := range a {
			v.Boollist[i] = ("true" == s.(string))
		}
	default:
		return nil, errors.New("Not implemented")
	}
	return
}

func asStringArray(data []interface{}) []string {
	s := make([]string, len(data))
	for i, d := range data {
		s[i] = d.(string)
	}
	return s
}

func (self *JsonReader) enterJsonList(list []interface{}) (Selection, error) {
	s := &MySelection{}
	var i int
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, first bool) (next Selection, err error) {
		if len(key) > 0 {
			if first {
				keyFields := meta.Keys
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if self.jsonKeyMatches(keyFields, candidate, key) {
						state.SetKey(key)
						return self.enterJson(candidate)
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
				container := list[i].(map[string]interface{})
				// TODO: compound keys
				keyStrs := []string{container[meta.Keys[0]].(string)}
				key, err = CoerseKeys(meta, keyStrs)
				state.SetKey(key)
				return self.enterJson(container)
			}
		}
		return nil, nil
	}
	return s, nil
}

func (self *JsonReader) enterJson(container map[string]interface{}) (Selection, error) {
	s := &MySelection{}
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
		if value, found := container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return self.enterJsonList(value.([]interface{}))
			} else {
				return self.enterJson(value.(map[string]interface{}))
			}
		}
		return
	}
	s.OnRead = func (state *WalkState, meta schema.HasDataType) (val *Value, err error) {
		if value, found := container[meta.GetIdent()]; found {
			return self.readLeafOrLeafList(meta, value)
		}
		return
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
