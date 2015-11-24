package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"schema"
)

type JsonReader struct {
	In     io.Reader
	values map[string]interface{}
}

func NewJsonReader(in io.Reader) *JsonReader {
	r := &JsonReader{In: in}
	return r
}

func (self *JsonReader) Node() (Node, error) {
	values, err := self.decode()
	if err != nil {
		return nil, err
	}

	return self.Container(values), nil
}

func (self *JsonReader) Selector(state *WalkState) (*Selection, error) {
	var n Node
	values, err := self.decode()
	if err != nil {
		return nil, err
	}
	if schema.IsList(state.SelectedMeta()) {
		if state.InsideList() {
			n = self.Container(values)
		} else {
			ident := state.SelectedMeta().GetIdent()
			foundValues, found := values[ident]
			list, ok := foundValues.([]interface{})
			if len(values) != 1 || !found || !ok {
				msg := fmt.Sprintf("Expected { %s: [] }", ident)
				return nil, errors.New(msg)
			}
			n = self.List(list)
		}
	} else {
		n = self.Container(values)
	}
	return NewSelectionFromState(n, state), nil
}

func (self *JsonReader) decode() (map[string]interface{}, error) {
	if self.values == nil {
		d := json.NewDecoder(self.In)
		if err := d.Decode(&self.values); err != nil {
			return nil, err
		}
	}
	return self.values, nil
}

func (self *JsonReader) readLeafOrLeafList(meta schema.HasDataType, data interface{}) (v *Value, err error) {
	// TODO: Consider using CoerseValue
	v = &Value{Type: meta.GetDataType()}
	switch v.Type.Format {
	case schema.FMT_INT64:
		v.Int64 = int64(data.(float64))
	case schema.FMT_INT64_LIST:
		a := data.([]interface{})
		v.Int64list = make([]int64, len(a))
		for i, f := range a {
			v.Int64list[i] = int64(f.(float64))
		}
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
		switch vdata := data.(type) {
		case string:
			s := data.(string)
			v.Bool = ("true" == s)
		case bool:
			v.Bool = vdata
		}
	case schema.FMT_BOOLEAN_LIST:
		a := data.([]interface{})
		v.Boollist = make([]bool, len(a))
		for i, s := range a {
			v.Boollist[i] = ("true" == s.(string))
		}
	case schema.FMT_ENUMERATION:
		v.SetEnumByLabel(data.(string))
	case schema.FMT_ENUMERATION_LIST:
		strlist := InterfaceToStrlist(data)
		if len(strlist) > 0 {
			v.SetEnumListByLabels(strlist)
		} else {
			intlist := InterfaceToIntlist(data)
			v.SetEnumList(intlist)
		}
	default:
		msg := fmt.Sprint("JSON reading value type not implemented ", meta.GetDataType().Format)
		return nil, errors.New(msg)
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

func (self *JsonReader) List(list []interface{}) Node {
	s := &MyNode{Label: "JSON Read List"}
	var i int
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		if new {
			panic("Cannot write to JSON reader")
		}
		if len(key) > 0 {
			if first {
				keyFields := meta.Keys
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if self.jsonKeyMatches(keyFields, candidate, key) {
						sel.State.SetKey(key)
						return self.Container(candidate), nil
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
				if len(meta.Keys) > 0 {
					// TODO: compound keys
					keyData, hasKey := container[meta.Keys[0]]
					// Key may legitimately not exist when inserting new data
					if hasKey {
						keyStrs := []string{keyData.(string)}
						key, err = CoerseKeys(meta, keyStrs)
						sel.State.SetKey(key)
					}
				}
				return self.Container(container), nil
			}
		}
		return nil, nil
	}
	return s
}

func (self *JsonReader) Container(container map[string]interface{}) Node {
	s := &MyNode{Label: "JSON Read Container"}
	s.OnChoose = func(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
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
					break
				} else {
					m = prop
				}
			}
			if aligned {
				return m, nil
			}
		}
		msg := fmt.Sprintf("No discriminating data for choice schema %s ", state.String())
		return nil, &browseError{Msg: msg}
	}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (child Node, e error) {
		if new {
			panic("Cannot write to JSON reader")
		}
		if value, found := container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return self.List(value.([]interface{})), nil
			} else {
				return self.Container(value.(map[string]interface{})), nil
			}
		}
		return
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (val *Value, err error) {
		if value, found := container[meta.GetIdent()]; found {
			return self.readLeafOrLeafList(meta, value)
		}
		return
	}
	return s
}

func (self *JsonReader) jsonKeyMatches(keyFields []string, candidate map[string]interface{}, key []*Value) bool {
	for i, field := range keyFields {
		if candidate[field] != key[i].String() {
			return false
		}
	}
	return true
}