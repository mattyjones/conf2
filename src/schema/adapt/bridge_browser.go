package adapt
import (
	"schema/browse"
	"schema"
	"schema/yang"
)

type BridgeBrowser struct {
	Meta *schema.Module
	Bridges map[string]*Bridge
}

func NewBridgeBrowser() *BridgeBrowser {
	meta, err := yang.LoadModuleFromByteArray([]byte(bridgeBrowserYang), nil)
	if err != nil {
		panic(err.Error())
	}
	return &BridgeBrowser{
		Bridges : make(map[string]*Bridge, 5),
		Meta : meta,
	}
}

func (bb *BridgeBrowser) Module() *schema.Module {
	return bb.Meta
}

func (bb *BridgeBrowser) AddBridge(name string, bridge *Bridge) {
	bb.Bridges[name] = bridge
}

func (bb *BridgeBrowser) Selector(path *browse.Path, strategy browse.Strategy) (browse.Selection, *browse.WalkState, error) {
	s := &browse.MySelection{}
	s.OnSelect = func (state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
		switch meta.GetIdent() {
			case "bridges":
				return bb.SelectBridges(bb.Bridges)
		}
		return nil, nil
	}
	return browse.WalkPath(browse.NewWalkState(bb.Meta), s, path)
}

func (bb *BridgeBrowser) SelectBridges(bridges map[string]*Bridge) (browse.Selection, error) {
	s := &browse.MySelection{}
	index := newBridgeIndex(bridges)
	s.OnNext = index.Index.OnNext
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
		internal := index.Selected.internal.Module()
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
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
			case "name":
				return &browse.Value{Str : index.Index.CurrentKey()}, nil
		}
		return browse.ReadField(meta, index.Selected)
	}

	return s, nil
}


func (bb *BridgeBrowser) selectMapping(mapping *BridgeMapping, external schema.MetaList, internal schema.MetaList) (browse.Selection, error) {
	s := &browse.MySelection{}
	index := newMappingIndex(mapping.Children)
	s.OnNext = index.Index.OnNext
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
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
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "externalIdent":
			return &browse.Value{Str:index.Index.CurrentKey()}, nil
		}
		return browse.ReadField(meta, index.Selected)
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, val *browse.Value) error {
		switch op {
		case browse.CREATE_LIST_ITEM:
			index.Selected = NewBridgeMapping("")
		case browse.UPDATE_VALUE:
			switch meta.GetIdent() {
			case "externalIdent":
				mapping.Children[val.Str] = index.Selected
			default:
				// case "internalIdent":
				err := browse.WriteField(meta.(schema.HasDataType), index.Selected, val)
				return err
			}
		}
		return nil
	}
	return s, nil
}

func (bb *BridgeBrowser) findMetaList(parent schema.MetaList, ident string) (child schema.MetaList) {
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

func (bb *BridgeBrowser) selectFieldOptions(field schema.MetaList) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
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
				if ! schema.IsLeaf(m) && ! schema.IsList(m) {
					v.Strlist = append(v.Strlist, m.GetIdent())
				}
			}
		}
		return v, nil
	}

	return s, nil
}

type bridgeIndex struct {
	Index browse.StringIndex
	Data map[string]*Bridge
	Selected *Bridge
}

func newBridgeIndex(data map[string]*Bridge) *bridgeIndex {
	ndx := &bridgeIndex{Data:data}
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
	Index browse.StringIndex
	Data map[string]*BridgeMapping
	Selected *BridgeMapping
}

func newMappingIndex(data map[string]*BridgeMapping) *mappingIndex {
	ndx := &mappingIndex{Data:data}
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

var bridgeBrowserYang = `
module bridge {
    prefix "bridge";
    namespace "conf2.org/bridge";
    revision 0000-00-00 {
    	description "Bridges transform one schema into another given a mapping";
    }
    grouping field-options {
		leaf-list leafs {
			type string;
		}
		leaf-list containers {
			type string;
		}
		leaf-list lists {
			type string;
		}
    }

    grouping meta-mapping {
        list mapping {
            key "externalIdent";
            leaf externalIdent {
                type string;
            }
            leaf internalIdent {
                type string;
            }
            uses meta-mapping;
        }
		container externalOptions {
			uses field-options;
		}
		container internalOptions {
			uses field-options;
		}
    }

    list bridges {
        key "name";
        leaf name {
            type string;
        }
        uses meta-mapping;
    }
}
`