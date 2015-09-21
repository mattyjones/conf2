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

func (bb *BridgeBrowser) RootSelector() (browse.Selection, error) {
	s := &browse.MySelection{}
	s.State.Meta = bb.Meta
	s.OnSelect = func () (browse.Selection, error) {
		switch s.State.Position.GetIdent() {
			case "bridges":
				s.State.Found = true
				return bb.selectBridges(bb.Bridges)
		}
		return nil, nil
	}
	return s, nil
}

func (bb *BridgeBrowser) selectBridges(bridges map[string]*Bridge) (browse.Selection, error) {
	s := &browse.MySelection{}
	index := &bridgeIndex{Data:bridges}
	s.OnNext = index.Index.OnNext
	s.OnSelect = func() (browse.Selection, error) {
		switch s.State.Position.GetIdent() {
		case "mapping":
			s.State.Found = index.Selected != nil
			if s.State.Found {
				return bb.selectMapping(index.Selected.Mapping)
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *browse.Value) error {
		return browse.ReadField(s.State.Position.(schema.HasDataType), index.Selected, val)
	}
	return s, nil
}


func (bb *BridgeBrowser) selectMapping(mapping *MetaListMapping) (browse.Selection, error) {
	s := &browse.MySelection{}
	var index = &mappingIndex{Data:mapping.Mapping}
	s.OnNext = index.Index.OnNext
	s.OnSelect = func() (browse.Selection, error) {
		switch s.State.Position.GetIdent() {
		case "mapping":
			s.State.Found = index.Selected != nil
			if s.State.Found {
				return bb.selectMapping(index.Selected.(*MetaListMapping))
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *browse.Value) error {
		return browse.ReadField(s.State.Position.(schema.HasDataType), index.Selected, val)
	}
	return s, nil
}

type bridgeIndex struct {
	Index browse.StringIndex
	Data map[string]*Bridge
	Selected *Bridge
}

func (impl bridgeIndex) Select(key string) (found bool) {
	impl.Selected, found = impl.Data[key]
	return
}

func (impl bridgeIndex) Build() []string {
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
	Data map[string]BridgeMapping
	Selected BridgeMapping
}

func (impl mappingIndex) Select(key string) (found bool) {
	impl.Selected, found = impl.Data[key]
	return
}

func (impl mappingIndex) Build() []string {
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

    grouping meta-mapping {
        list mapping {
            key "from";
            leaf from {
                type string;
            }
            leaf to {
                type string;
            }
            uses meta-mapping;
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