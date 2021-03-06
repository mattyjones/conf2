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

func (bb *BucketData) Select() *Selection {
	return Select(bb.Meta, bb.Node())
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
	s.OnSelect = func(sel *Selection, r ContainerRequest) (Node, error) {
		var data interface{}
		if r.New {
			if schema.IsList(r.Meta) {
				data = make([]map[string]interface{}, 0, 10)
			} else {
				data = make(map[string]interface{})
			}
			container[b.KeyMap(sel, r.Meta.GetIdent())] = data
		} else {
			data = container[b.KeyMap(sel, r.Meta.GetIdent())]
		}
		if data != nil {
			if schema.IsList(r.Meta) {
				return b.List(container, data.([]map[string]interface{})), nil
			} else {
				// TODO: Silently ignoring unexpected format. We should be *less*
				// tolerant and fail here otherwise we silently ignore bad data.
				if c, valid := data.(map[string]interface{}); valid {
					return b.Container(c), nil
				}
			}
		}
		return nil, nil
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *Value) error {
		return b.UpdateLeaf(sel, container, meta.(schema.HasDataType), val)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*Value, error) {
		return b.ReadLeaf(sel, container, meta)
	}
	return s
}

func (b *Bucket) ReadKey(sel *Selection, container map[string]interface{}, meta *schema.List) (key []*Value, err error) {
	keyMeta := meta.KeyMeta()
	key = make([]*Value, len(keyMeta))
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
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		selected = nil
		if r.New {
			selection := make(map[string]interface{})
			list = append(list, selection)
			parent[b.KeyMap(sel, r.Meta.GetIdent())] = list
			return b.Container(selection), r.Key, nil
		} else {
			if len(r.Key) > 0 {
				if !r.First {
					return nil, nil, nil
				}
				// looping not very efficient, but we do not have an index
				for _, candidate := range list {
					// TODO: Support compound keys
					if candidate[b.KeyMap(sel, r.Meta.Key[0])] == r.Key[0].Value() {
						selected = candidate
						break
					}
				}
			} else {
				if r.First {
					i = 0
				} else {
					i++
				}
				if i < len(list) {
					selected = list[i]
				}
				var err error
				if r.Key, err = b.ReadKey(sel, selected, r.Meta); err != nil {
					return nil, nil, err
				}
			}
		}
		if selected != nil {
			return b.Container(selected), r.Key, nil
		}
		return nil, nil, nil
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *Value) error {
		return b.UpdateLeaf(sel, selected, meta.(schema.HasDataType), val)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*Value, error) {
		return b.ReadLeaf(sel, selected, meta)
	}

	return s
}

func (b *Bucket) ReadLeaf(sel *Selection, container map[string]interface{}, m schema.HasDataType) (*Value, error) {
	return SetValue(m.GetDataType(), container[b.KeyMap(sel, m.GetIdent())])
}

func (b *Bucket) UpdateLeaf(sel *Selection, container map[string]interface{}, m schema.HasDataType, v *Value) error {
	container[b.KeyMap(sel, m.GetIdent())] = v.Value()
	return nil
}
