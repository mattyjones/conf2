package adapt

import (
	"schema"
	"schema/browse"
	"schema/yang"
)

type BridgeData struct {
	Meta    schema.MetaList
	Bridges map[string]*Bridge
}

func NewBridgeData() *BridgeData {
	meta, err := yang.LoadModule(yang.YangPath(), "bridge.yang")
	if err != nil {
		panic(err.Error())
	}
	return &BridgeData{
		Bridges: make(map[string]*Bridge, 5),
		Meta:    meta,
	}
}

func (bb *BridgeData) Schema() schema.MetaList {
	return bb.Meta
}

func (bb *BridgeData) AddBridge(name string, bridge *Bridge) {
	bb.Bridges[name] = bridge
}

func (bb *BridgeData) Selector(path *browse.Path) (*browse.Selection, error) {
	s := &browse.MyNode{}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (browse.Node, error) {
		switch meta.GetIdent() {
		case "bridges":
			return bb.SelectBridges(bb.Bridges)
		}
		return nil, nil
	}
	return browse.WalkPath(browse.NewSelection(s, bb.Meta), path)
}

func (bb *BridgeData) SelectBridges(bridges map[string]*Bridge) (browse.Node, error) {
	s := &browse.MyNode{}
	index := newBridgeIndex(bridges)
	s.OnNext = func(state *browse.Selection, meta *schema.List, new bool, key []*browse.Value, first bool) (next browse.Node, err error) {
		if new {
			return nil, nil
		}
		var hasNext bool
		if hasNext, err = index.Index.OnNext(state, meta, key, first); hasNext {
			return s, err
		}
		return nil, nil
	}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (browse.Node, error) {
		internal := index.Selected.internal.State.SelectedMeta()
		external := index.Selected.external
		switch meta.GetIdent() {
		case "mapping":
			return bb.selectMapping(index.Selected.Mapping, external, internal)
		case "externalOptions":
			return bb.selectFieldOptions(external)
		case "internalOptions":
			return bb.selectFieldOptions(internal)
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "name":
			return &browse.Value{Str: index.Index.CurrentKey()}, nil
		}
		return browse.ReadField(meta, index.Selected)
	}

	return s, nil
}

func (bb *BridgeData) selectMapping(mapping *BridgeMapping, external schema.MetaList, internal schema.MetaList) (browse.Node, error) {
	s := &browse.MyNode{}
	index := newMappingIndex(mapping.Children)
	s.OnNext = func(state *browse.Selection, meta *schema.List, new bool, key []*browse.Value, first bool) (next browse.Node, err error) {
		if new {
			index.Selected = NewBridgeMapping("")
			return s, nil
		} else {
			var hasNext bool
			if hasNext, err = index.Index.OnNext(state, meta, key, first); hasNext {
				return s, err
			}
		}
		return nil, nil
	}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (browse.Node, error) {
		externalChild := bb.findMetaList(external, index.Index.CurrentKey())
		if externalChild == nil {
			return nil, nil
		}
		internalChild := bb.findMetaList(internal, index.Selected.InternalIdent)
		switch meta.GetIdent() {
		case "mapping":
			return bb.selectMapping(index.Selected, externalChild, internalChild)
		case "externalOptions":
			return bb.selectFieldOptions(externalChild)
		case "internalOptions":
			return bb.selectFieldOptions(internalChild)
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "externalIdent":
			return &browse.Value{Str: index.Index.CurrentKey()}, nil
		}
		return browse.ReadField(meta, index.Selected)
	}
	s.OnWrite = func(state *browse.Selection, meta schema.HasDataType, val *browse.Value) error {
		switch meta.GetIdent() {
		case "externalIdent":
			mapping.Children[val.Str] = index.Selected
		default:
			// case "internalIdent":
			return browse.WriteField(meta.(schema.HasDataType), index.Selected, val)
		}
		return nil
	}
	return s, nil
}

func (bb *BridgeData) findMetaList(parent schema.MetaList, ident string) (child schema.MetaList) {
	childMeta := schema.FindByIdent2(parent, ident)
	if childMeta != nil {
		var isList bool
		child, isList = childMeta.(schema.MetaList)
		if isList {
			return child
		}
	}
	return nil
}

func (bb *BridgeData) selectFieldOptions(field schema.MetaList) (browse.Node, error) {
	s := &browse.MyNode{}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		i := schema.NewMetaListIterator(field, true)
		v := &browse.Value{}
		v.Strlist = make([]string, 0, 10)
		ident := meta.GetIdent()
		for i.HasNextMeta() {
			m := i.NextMeta()
			switch ident {
			case "leafs":
				if schema.IsLeaf(m) {
					v.Strlist = append(v.Strlist, m.GetIdent())
				}
			case "lists":
				if schema.IsList(m) {
					v.Strlist = append(v.Strlist, m.GetIdent())
				}
			default:
				if !schema.IsLeaf(m) && !schema.IsList(m) {
					v.Strlist = append(v.Strlist, m.GetIdent())
				}
			}
		}
		return v, nil
	}

	return s, nil
}

type bridgeIndex struct {
	Index    browse.StringIndex
	Data     map[string]*Bridge
	Selected *Bridge
}

func newBridgeIndex(data map[string]*Bridge) *bridgeIndex {
	ndx := &bridgeIndex{Data: data}
	ndx.Index.Builder = ndx
	return ndx
}

func (impl *bridgeIndex) Select(key string) (found bool) {
	impl.Selected, found = impl.Data[key]
	return
}

func (impl *bridgeIndex) Build() []string {
	index := make([]string, len(impl.Data))
	j := 0
	for key, _ := range impl.Data {
		index[j] = key
		j++
	}
	return index
}

type mappingIndex struct {
	Index    browse.StringIndex
	Data     map[string]*BridgeMapping
	Selected *BridgeMapping
}

func newMappingIndex(data map[string]*BridgeMapping) *mappingIndex {
	ndx := &mappingIndex{Data: data}
	ndx.Index.Builder = ndx
	return ndx
}

func (impl *mappingIndex) Select(key string) (found bool) {
	impl.Selected, found = impl.Data[key]
	return
}

func (impl *mappingIndex) Build() []string {
	index := make([]string, len(impl.Data))
	j := 0
	for key, _ := range impl.Data {
		index[j] = key
		j++
	}
	return index
}
