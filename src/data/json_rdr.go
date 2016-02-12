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

func (self *JsonReader) Node() (Node) {
	var err error
	if self.values == nil {
		self.values, err = self.decode()
		if err != nil {
			return ErrorNode{Err:err}
		}
	}
	return JsonContainerReader(self.values)
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

func leafOrLeafListJsonReader(meta schema.HasDataType, data interface{}) (v *Value, err error) {
	// TODO: Consider using CoerseValue
	v = &Value{Type: meta.GetDataType()}
	switch v.Type.Format() {
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
	case schema.FMT_ANYDATA:
		v.Data = &AnyJson{
			container: data.(map[string]interface{}),
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

func JsonListReader(list []interface{}) Node {
	s := &MyNode{Label: "JSON Read List"}
	var i int
	s.OnNext = func(sel *Selection, r ListRequest) (next Node, key []*Value, err error) {
		key = r.Key
		if r.New {
			panic("Cannot write to JSON reader")
		}
		if len(r.Key) > 0 {
			if r.First {
				keyFields := r.Meta.Key
				for ; i < len(list); i++ {
					candidate := list[i].(map[string]interface{})
					if jsonKeyMatches(keyFields, candidate, key) {
						sel.path.key = key
						return JsonContainerReader(candidate), r.Key, nil
					}
				}
			}
		} else {
			if r.First {
				i = 0
			} else {
				i++
			}
			if i < len(list) {
				container := list[i].(map[string]interface{})
				if len(r.Meta.Key) > 0 {
					// TODO: compound keys
					if keyData, hasKey := container[r.Meta.Key[0]]; hasKey {
						// Key may legitimately not exist when inserting new data
						key = SetValues(r.Meta.KeyMeta(), keyData)
					}
				}
				return JsonContainerReader(container), key, nil
			}
		}
		return nil, nil, nil
	}
	return s
}

func JsonContainerReader(container map[string]interface{}) Node {
	s := &MyNode{Label: "JSON Read Container"}
	var divertedList Node
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
	s.OnSelect = func(state *Selection, r ContainerRequest) (child Node, e error) {
		if r.New {
			panic("Cannot write to JSON reader")
		}
		if value, found := container[r.Meta.GetIdent()]; found {
			if schema.IsList(r.Meta) {
				return JsonListReader(value.([]interface{})), nil
			} else {
				return JsonContainerReader(value.(map[string]interface{})), nil
			}
		}
		return
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (val *Value, err error) {
		if value, found := container[meta.GetIdent()]; found {
			return leafOrLeafListJsonReader(meta, value)
		}
		return
	}
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		if divertedList != nil {
			return nil, nil, nil
		}
		// divert to list handler
		foundValues, found := container[r.Meta.GetIdent()]
		list, ok := foundValues.([]interface{})
		if len(container) != 1 || !found || !ok {
			msg := fmt.Sprintf("Expected { %s: [] }", r.Meta.GetIdent())
			return nil, nil, errors.New(msg)
		}
		divertedList = JsonListReader(list)
		return divertedList.Next(sel, r)
	}
	return s
}

func jsonKeyMatches(keyFields []string, candidate map[string]interface{}, key []*Value) bool {
	for i, field := range keyFields {
		if candidate[field] != key[i].String() {
			return false
		}
	}
	return true
}
