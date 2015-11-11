package browse

import (
	"fmt"
	"reflect"
	"schema"
	"strconv"
	"strings"
	"conf2"
)

// Stores stuff in memory according to a given schema.  Useful in testing or store of
// generic settings.
type BucketBrowser struct {
	Meta      *schema.Module
	Bucket    map[string]interface{}
	PathDelim string
}

func NewBucketBrowser(module *schema.Module) (bb *BucketBrowser) {
	bb = &BucketBrowser{Meta: module, PathDelim: "."}
	bb.Bucket = make(map[string]interface{}, 10)
	return bb
}

func (bb *BucketBrowser) Selector(path *Path) (*Selection, error) {
	if root, err := bb.selectContainer(bb.Bucket); err == nil {
		return WalkPath(NewSelection(root, bb.Meta), path)
	} else {
		return nil, err
	}
}

func (bb *BucketBrowser) Schema() schema.MetaList {
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
			return nil, &browseError{Msg: fmt.Sprintf("Bad type %s on %s", reflect.TypeOf(v), seg)}
		}
		if v == nil {
			return nil, &browseError{Msg: fmt.Sprintf("%s not found", seg)}
		}
	}
	return v, nil
}

func (bb *BucketBrowser) selectContainer(container map[string]interface{}) (Node, error) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		if data, found := container[meta.GetIdent()]; found {
			if schema.IsList(meta) {
				return bb.selectList(container, data.([]map[string]interface{}))
			}
			return bb.selectContainer(data.(map[string]interface{}))
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			bb.updateLeaf(meta.(schema.HasDataType), container, val)
		case CREATE_CONTAINER:
			container[meta.GetIdent()] = make(map[string]interface{}, 10)
		case CREATE_LIST:
			container[meta.GetIdent()] = make([]map[string]interface{}, 0, 10)
		}
		return nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return bb.readLeaf(meta, container)
	}
	return s, nil
}

func (bb *BucketBrowser) readKey(meta *schema.List, container map[string]interface{}) (key []*Value, err error) {
	keyMeta := meta.KeyMeta()
	key = make([]*Value, len(keyMeta))
	for i, m := range keyMeta {
		if key[i], err = bb.readLeaf(m, container); err != nil {
			return nil, err
		}
	}
	return
}

func (bb *BucketBrowser) selectList(parent map[string]interface{}, initialList []map[string]interface{}) (Node, error) {
	var i int
	list := initialList
	s := &MyNode{}
	var next Node
	var selection map[string]interface{}
	s.OnNext = func(state *Selection, meta *schema.List, key []*Value, isFirst bool) (Node, error) {
		if next != nil {
			s := next
			next = nil
			return s, nil
		}
		selection = nil
		if len(key) > 0 {
			if !isFirst {
				return nil, nil
			}

			// looping not very efficient, but we do not have an index
			for _, candidate := range list {
				// TODO: Support compound keys
				if candidate[meta.Keys[0]] == key[0].Value() {
					selection = candidate
					state.SetKey(key)
					break
				}
			}
		} else {
			if isFirst {
				i = 0
			} else {
				i++
			}
			if i < len(initialList) {
				selection = list[i]
			}
			if key, err := bb.readKey(meta, selection); err != nil {
				return nil, err
			} else {
				state.SetKey(key)
			}
		}
		if selection != nil {
			return bb.selectContainer(selection)
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			selection = make(map[string]interface{}, 10)
			list = append(list, selection)
			// list reference may have changed so update parent
			parent[meta.GetIdent()] = list
			var err error
			if next, err = bb.selectContainer(selection); err != nil {
				return err
			}
		}
		return nil
	}

	return s, nil
}

func (bb *BucketBrowser) readLeaf(m schema.HasDataType, container map[string]interface{}) (*Value, error) {
	return SetValue(m.GetDataType(), container[m.GetIdent()])
}

func (bb *BucketBrowser) updateLeaf(m schema.HasDataType, container map[string]interface{}, v *Value) error {
	container[m.GetIdent()] = v.Value()
	return nil
}
