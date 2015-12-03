package data

import (
	"schema"
)

// Stores stuff in memory according to a given schema.  Useful in testing or store of
// generic settings.
type BucketData struct {
	Meta      schema.MetaList
	Root      map[string]interface{}
}

func (bb *BucketData) Node() (Node) {
	if bb.Root == nil {
		bb.Root = make(map[string]interface{})
	}
	b := &Bucket{
		KeyMap : DefaultKeyMapper,
	}
	return b.Container(bb.Root)
}

func (bb *BucketData) Schema() schema.MetaList {
	return bb.Meta
}

type Bucket struct {
	KeyMap KeyMapper
}

func NewBucket() *Bucket {
	return &Bucket{
		KeyMap : DefaultKeyMapper,
	}
}

type KeyMapper func(sel *Selection, ident string) (key string)

func DefaultKeyMapper(sel *Selection, ident string) (key string) {
	return ident
}

func (b *Bucket) Container(container map[string]interface{}) (Node) {
	s := &MyNode{}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
		var data interface{}
		if new {
			if schema.IsList(meta) {
				data = make([]map[string]interface{}, 0, 10)
			} else {
				data = make(map[string]interface{})
			}
			container[b.KeyMap(sel, meta.GetIdent())] = data
		} else {
			data = container[b.KeyMap(sel, meta.GetIdent())]
		}
		if data != nil {
			if schema.IsList(meta) {
				return b.List(container, data.([]map[string]interface{})), nil
			} else {
				return b.Container(data.(map[string]interface{})), nil
			}
		}
		return nil, nil
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *schema.Value) error {
		return b.UpdateLeaf(sel, container, meta.(schema.HasDataType), val)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return b.ReadLeaf(sel, container, meta)
	}
	return s
}

func (b *Bucket) ReadKey(sel *Selection, container map[string]interface{}, meta *schema.List) (key []*schema.Value, err error) {
	keyMeta := meta.KeyMeta()
	key = make([]*schema.Value, len(keyMeta))
	for i, m := range keyMeta {
		if key[i], err = b.ReadLeaf(sel, container, m); err != nil {
			return nil, err
		}
	}
	return
}

func (b *Bucket) List(parent map[string]interface{}, initialList []map[string]interface{}) (Node) {
	var i int
	list := initialList
	s := &MyNode{}
	var selected map[string]interface{}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*schema.Value, isFirst bool) (Node, error) {
		selected = nil
		if new {
			selection := make(map[string]interface{})
			list = append(list, selection)
			parent[b.KeyMap(sel, meta.GetIdent())] = list
			return b.Container(selection), nil
		} else {
			if len(key) > 0 {
				if !isFirst {
					return nil, nil
				}
				// looping not very efficient, but we do not have an index
				for _, candidate := range list {
					// TODO: Support compound keys
					if candidate[b.KeyMap(sel, meta.Keys[0])] == key[0].Value() {
						selected = candidate
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
				if i < len(list) {
					selected = list[i]
				}
				if key, err := b.ReadKey(sel, selected, meta); err != nil {
					return nil, err
				} else {
					sel.State.SetKey(key)
				}
			}
		}
		if selected != nil {
			return b.Container(selected), nil
		}
		return nil, nil
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *schema.Value) error {
		return b.UpdateLeaf(sel, selected, meta.(schema.HasDataType), val)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return b.ReadLeaf(sel, selected, meta)
	}

	return s
}

func (b *Bucket) ReadLeaf(sel *Selection, container map[string]interface{}, m schema.HasDataType) (*schema.Value, error) {
	return schema.SetValue(m.GetDataType(), container[b.KeyMap(sel, m.GetIdent())])
}

func (b *Bucket) UpdateLeaf(sel *Selection, container map[string]interface{}, m schema.HasDataType, v *schema.Value) error {
	container[b.KeyMap(sel, m.GetIdent())] = v.Value()
	return nil
}
