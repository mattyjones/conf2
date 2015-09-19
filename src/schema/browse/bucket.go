package browse
import (
	"schema"
	"strings"
	"strconv"
	"fmt"
	"reflect"
)

// Stores stuff in memory according to a given schema.  Useful in testing or store of
// generic settings.
type BucketBrowser struct {
	Module *schema.Module
	Bucket map[string]interface{}
	PathDelim string
}

func NewBucketBrowser(module *schema.Module) (bb *BucketBrowser) {
	bb = &BucketBrowser{Module:module, PathDelim:"."}
	bb.Bucket = make(map[string]interface{}, 10)
	return bb
}

func (bb *BucketBrowser) RootSelector() (s Selection, err error) {
	s, err = bb.selectContainer(bb.Bucket)
	s.WalkState().Meta = bb.Module
	return
}

//
// Example:
//  bb.Read("foo.10.bar.blah.0")
func (bb *BucketBrowser) Read(path string) (interface{}, error) {
	segments := strings.Split(path, bb.PathDelim)
	var v interface{}
	v = bb.Bucket
	for _, seg := range segments {
		n, _ := strconv.Atoi(seg)
		switch x := v.(type) {
		case []interface{}:
			v = x[n]
		case []map[string]interface{}:
			v = x[n]
		case map[string]interface{}:
			v = x[seg]
		default:
			return nil, &browseError{Msg:fmt.Sprintf("Bad type %s on %s", reflect.TypeOf(v), seg)}
		}
		if v == nil {
			return nil, &browseError{Msg:fmt.Sprintf("%s not found", seg)}
		}
	}
	return v, nil
}

func (bb *BucketBrowser) selectContainer(container map[string]interface{}) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		var data interface{}
		data, s.State.Found = container[s.State.Position.GetIdent()]
		if s.State.Found {
			if schema.IsList(s.State.Position) {
				return bb.enterList(container, data.([]map[string]interface{}))
			}
			return bb.selectContainer(data.(map[string]interface{}))
		}
		return nil, nil
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			bb.updateLeaf(s.State.Position.(schema.HasDataType), container, val)
		case CREATE_CHILD:
			container[s.State.Position.GetIdent()] = make(map[string]interface{}, 10)
		case CREATE_LIST:
			container[s.State.Position.GetIdent()] = make([]map[string]interface{}, 0, 10)
		}
		return nil
	}
	s.OnRead = func(val *Value) error {
		return bb.readLeaf(s.State.Position.(schema.HasDataType), container, val)
	}
	return s, nil
}

func (bb *BucketBrowser) enterList(parent map[string]interface{}, initialList []map[string]interface{}) (Selection, error) {
	var i int
	list := initialList
	s := &MySelection{}
	s.OnNext = func(keys []Value, isFirst bool) (bool, error) {
		if len(keys) > 0 {
fmt.Printf("bucket - OnNext keys=%v\n", keys)
			keyFieldNames := s.State.Meta.(*schema.List).Keys
			//var candidate map[string]interface{}
			// looping not very efficient, but we do not have an index
			for j, candidate := range list {
				for k, keyName := range keyFieldNames {
					if candidate[keyName] != keys[k].Value() {
						break
					} else {
						lastKey := (k == len(keyFieldNames) - 1)
						if lastKey {
							i = j
fmt.Printf("bucket - OnNext MATCH\n")
							return true, nil
						}
					}
				}
			}
fmt.Printf("bucket - OnNext NO MATCH\n")
			return false, nil
		} else {
			if isFirst {
				i = 0
			} else {
				i++
			}
		}
		return len(list) > i, nil
	}
	s.OnSelect = func() (Selection, error) {
		var data interface{}
		data, s.State.Found = list[i][s.State.Position.GetIdent()]
		if s.State.Found {
			if schema.IsList(s.State.Position) {
				return bb.enterList(list[i], data.([]map[string]interface{}))
			}
			return bb.selectContainer(data.(map[string]interface{}))
		}
		return nil, nil
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			bb.updateLeaf(s.State.Position.(schema.HasDataType), list[i], val)
		case CREATE_LIST_ITEM:
			created := make(map[string]interface{}, 10)
			list = append(list, created)
			i = len(list) - 1
			// list reference may have changed so update parent
			parent[s.State.Meta.GetIdent()] = list
		case CREATE_LIST:
			list[i][s.State.Meta.GetIdent()] = make([]map[string]interface{}, 0, 10)
		case CREATE_CHILD:
			child := make(map[string]interface{}, 10)
			list[i][s.State.Position.GetIdent()] = child
		}
		return nil
	}
	s.OnRead = func(val *Value) error {
		return bb.readLeaf(s.State.Position.(schema.HasDataType), list[i], val)
	}

	return s, nil
}

func (bb *BucketBrowser) readLeaf(m schema.HasDataType, container map[string]interface{}, v *Value) error {
	v.IsList = ! schema.IsLeaf(m)
	v.Type = m.GetDataType()
	switch v.Type.Format {
	case schema.FMT_STRING:
		if v.IsList {
			v.Strlist = container[m.GetIdent()].([]string)
		} else {
			v.Str = container[m.GetIdent()].(string)
		}
	case schema.FMT_INT32:
		if v.IsList {
			v.Intlist = container[m.GetIdent()].([]int)
		} else {
			v.Int = container[m.GetIdent()].(int)
		}
	}
	return nil
}

func (bb *BucketBrowser) updateLeaf(m schema.HasDataType, container map[string]interface{}, v *Value) (error) {
	switch m.GetDataType().Format {
	case schema.FMT_STRING:
		if v.IsList {
			container[m.GetIdent()] = v.Strlist
		} else {
			container[m.GetIdent()] = v.Str
		}
	case schema.FMT_INT32:
		if v.IsList {
			container[m.GetIdent()] = v.Intlist
		} else {
			container[m.GetIdent()] = v.Int
		}
	}
	return nil
}