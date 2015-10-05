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
	Meta *schema.Module
	Bucket map[string]interface{}
	PathDelim string
}

func NewBucketBrowser(module *schema.Module) (bb *BucketBrowser) {
	bb = &BucketBrowser{Meta:module, PathDelim:"."}
	bb.Bucket = make(map[string]interface{}, 10)
	return bb
}

func (bb *BucketBrowser) Selector(path *Path, strategy Strategy) (s Selection, state *WalkState, err error) {
	if s, err = bb.selectContainer(bb.Bucket); err == nil {
		return WalkPath(NewWalkState(bb.Meta), s, path)
	}
	return
}

func (bb *BucketBrowser) Module() *schema.Module {
	return bb.Meta
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
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (Selection, error) {
		if data, found := container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return bb.enterList(container, data.([]map[string]interface{}))
			}
			return bb.selectContainer(data.(map[string]interface{}))
		}
		return nil, nil
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			bb.updateLeaf(meta.(schema.HasDataType), container, val)
		case CREATE_CHILD:
			container[meta.GetIdent()] = make(map[string]interface{}, 10)
		case CREATE_LIST:
			container[meta.GetIdent()] = make([]map[string]interface{}, 0, 10)
		}
		return nil
	}
	s.OnRead = func(state *WalkState, meta schema.HasDataType) (*Value, error) {
		return bb.readLeaf(meta, container)
	}
	return s, nil
}

func (bb *BucketBrowser) enterList(parent map[string]interface{}, initialList []map[string]interface{}) (Selection, error) {
	var i int
	list := initialList
	s := &MySelection{}
	var selection map[string]interface{}
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, isFirst bool) (bool, error) {
		selection = nil
		if len(key) > 0 {
			if !isFirst {
				return false, nil
			}

			// looping not very efficient, but we do not have an index
			for _, candidate := range list {
				// TODO: Support compound keys
				if candidate[meta.Keys[0]] == key[0].Value() {
					selection = candidate
					return true, nil
				}
			}
			return false, nil
		} else {
			if isFirst {
				i = 0
			} else {
				i++
			}
			if i < len(initialList) {
				selection = list[i]
			}
		}
		return selection != nil, nil
	}
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (Selection, error) {
		if data, found := selection[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return bb.enterList(selection, data.([]map[string]interface{}))
			}
			return bb.selectContainer(data.(map[string]interface{}))
		}
		return nil, nil
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			bb.updateLeaf(meta.(schema.HasDataType), selection, val)
		case CREATE_LIST_ITEM:
			selection = make(map[string]interface{}, 10)
			list = append(list, selection)
			// list reference may have changed so update parent
			parent[meta.GetIdent()] = list
		case CREATE_LIST:
			selection[meta.GetIdent()] = make([]map[string]interface{}, 0, 10)
		case CREATE_CHILD:
			child := make(map[string]interface{}, 10)
			selection[meta.GetIdent()] = child
		}
		return nil
	}
	s.OnRead = func(state *WalkState, meta schema.HasDataType) (*Value, error) {
		return bb.readLeaf(meta, selection)
	}

	return s, nil
}

func (bb *BucketBrowser) readLeaf(m schema.HasDataType, container map[string]interface{}) (*Value, error) {
	return SetValue(m.GetDataType(), container[m.GetIdent()])
}

func (bb *BucketBrowser) updateLeaf(m schema.HasDataType, container map[string]interface{}, v *Value) (error) {
	container[m.GetIdent()] = v.Value()
	return nil
}