package browse

import (
	"io"
	"yang"
	"encoding/json"
)

type JsonTransmitter struct {
	in  io.Reader
}

func NewJsonTransmitter(in io.Reader) *JsonTransmitter {
	return &JsonTransmitter{in:in}
}

func (self *JsonTransmitter) GetSelector(meta yang.MetaList) (s *Selection, err error) {
	var values map[string]interface{}
	d := json.NewDecoder(self.in)
	if err = d.Decode(&values); err != nil {
		return
	}
	if yang.IsList(meta) {
		singleton := make([]interface{}, 1)
		singleton[0] = values
		s, err = selectJsonList(singleton)
		s.Iterate([]string{}, true)
	} else {
		s, err = selectJsonContainer(values)
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

func defaultSelector(meta yang.Meta, data interface{}) (s *Selection, err error) {
	switch meta.(type) {
	case *yang.List:
		// This doesn't compile event though it's true.  must be go reflection limitation
		//   return selectJsonList(data.([]map[string]interface{}))
		return selectJsonList(data.([]interface{}))
	case *yang.Container:
		return selectJsonContainer(data.(map[string]interface{}))
	}
	return
}

func selectJsonContainer(values map[string]interface{}) (s *Selection, err error) {
	s = &Selection{}
	var found bool
	var data interface{}
	s.Select = func(ident string) (child *Selection, e error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		data, found = values[ident]
		if !found {
			s.Position = nil
		} else if !yang.IsLeaf(s.Position) {
			return defaultSelector(s.Position, data)
		}
		return
	}
	s.Read = func (val *Value) (err error) {
		if !found {
			return
		}
		return readLeafOrLeafList(s.Position, data, val)
	}
	return
}

func selectJsonList(list []interface{}) (s *Selection, err error) {
	var i int
	s = &Selection{}
	var values map[string]interface{}
	var data interface{}
	var found bool
	s.Iterate = func(keys []string, first bool) (bool, error) {
		/* ignoring keys, cannot see the use case */
		if (first) {
			i = 0
		} else {
			i += 1
		}
		if i < len(list) {
			values = list[i].(map[string]interface{})
			return true, nil
		}
		return false, nil
	}
	s.Select = func(ident string) (child *Selection, e error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		data, found = values[ident]
		if !found {
			s.Position = nil
		} else if !yang.IsLeaf(s.Position) {
			return selectJsonContainer(data.(map[string]interface{}))
		}
		return
	}
	return
}
