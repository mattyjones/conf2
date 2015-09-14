package adapt
import (
	"schema"
	"schema/browse"
)

type BridgeBrowser struct {
	To browse.Browser
	From *schema.Module
	Map *Map
}

type Bridge struct {
	Map *Map
}

type Mapping struct {
	To string
}

type Map struct {
	Mappings map[string]Mapping
}

func NewMap() *Map {
	return &Map{
		Mappings : make(map[string]Mapping, 10),
	}
}

func (m *Map) MapMeta(from schema.Meta, toParent schema.MetaList) (schema.Meta, error) {
	if from == nil {
		return nil, nil
	}
	toIdent := from.GetIdent()
	if mapping, found := m.Mappings[from.GetIdent()]; found {
		toIdent = mapping.To
	}
	to := schema.FindByIdent2(toParent, toIdent)
	return to, nil
}

func (bb *BridgeBrowser) RootSelector() (browse.Selection, error) {
	root, err := bb.To.RootSelector()
	if err != nil {
		return nil, err
	}
	bridge := &Bridge{Map: bb.Map}
	return bridge.enterBridge(root)
}

func (b *Bridge) enterBridge(to browse.Selection) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func() (child browse.Selection, err error) {
		toState := to.WalkState()
		if toState.Position, err = b.Map.MapMeta(s.State.Position, toState.Meta); err == nil {
			var toChild browse.Selection
			if toChild, err = to.Select(); err == nil {
				s.WalkState().Found = to.WalkState().Found
				if toChild != nil {
					toChild.WalkState().Meta = toState.Position.(schema.MetaList)
					return b.enterBridge(toChild)
				}
			}
		}
		return
	}
	s.OnWrite = func(op browse.Operation, val *browse.Value) (err error) {
		toState := to.WalkState()
		if toState.Position, err = b.Map.MapMeta(s.State.Position, toState.Meta); err == nil {
			return to.Write(op, val)
		}
		return
	}
	s.OnRead = func(val *browse.Value) (err error) {
		toState := to.WalkState()
		if toState.Position, err = b.Map.MapMeta(s.State.Position, toState.Meta); err == nil {
			// TODO: txlate val
			return to.Write(browse.UPDATE_VALUE, val)
		}
		return
	}
	return s, nil
}

func (b *BridgeBrowser) Module() *schema.Module {
	return b.From
}



