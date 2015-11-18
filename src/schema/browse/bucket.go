package browse

import (
	"fmt"
	"reflect"
	"schema"
	"strconv"
	"strings"
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
	root := bb.selectContainer(bb.Bucket)
	return WalkPath(NewSelection(root, bb.Meta), path)
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

func (bb *BucketBrowser) selectContainer(container map[string]interface{}) (Node) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		var data interface{}
		if new {
			if schema.IsList(meta) {
				data = make([]map[string]interface{}, 0, 10)
			} else {
				data = make(map[string]interface{}, 10)
			}
			container[meta.GetIdent()] = data
		} else {
			data = container[meta.GetIdent()]
		}
		if data != nil {
			if schema.IsList(meta) {
				return bb.selectList(container, data.([]map[string]interface{})), nil
			} else {
				return bb.selectContainer(data.(map[string]interface{})), nil
			}
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *Value) error {
		return bb.updateLeaf(meta.(schema.HasDataType), container, val)
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return bb.readLeaf(meta, container)
	}
	return s
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

func (bb *BucketBrowser) selectList(parent map[string]interface{}, initialList []map[string]interface{}) (Node) {
	var i int
	list := initialList
	s := &MyNode{}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, isFirst bool) (Node, error) {
		var selection map[string]interface{}
		if new {
			selection = make(map[string]interface{}, 10)
			list = append(list, selection)
			parent[meta.GetIdent()] = list
			return bb.selectContainer(selection), nil
		} else {
			if len(key) > 0 {
				if !isFirst {
					return nil, nil
				}

				// looping not very efficient, but we do not have an index
				for _, candidate := range list {
					// TODO: Support compound keys
					if candidate[meta.Keys[0]] == key[0].Value() {
						selection = candidate
						sel.State.SetKey(key)
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
					sel.State.SetKey(key)
				}
			}
		}
		if selection != nil {
			return bb.selectContainer(selection), nil
		}
		return nil, nil
	}

	return s
}

func (bb *BucketBrowser) readLeaf(m schema.HasDataType, container map[string]interface{}) (*Value, error) {
	return SetValue(m.GetDataType(), container[m.GetIdent()])
}

func (bb *BucketBrowser) updateLeaf(m schema.HasDataType, container map[string]interface{}, v *Value) error {
	container[m.GetIdent()] = v.Value()
	return nil
}
