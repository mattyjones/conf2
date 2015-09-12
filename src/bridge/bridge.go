package bridge
import (
	"yang"
	"yang/browse"
	"fmt"
)

type BridgeBrowser struct {
	To browse.Browser
	From *yang.Module
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

func (m *Map) MapMeta(from yang.Meta, toParent yang.MetaList) (yang.Meta, error) {
	if from == nil {
		return nil, nil
	}
	toIdent := from.GetIdent()
	if mapping, found := m.Mappings[from.GetIdent()]; found {
		toIdent = mapping.To
	}
fmt.Printf("bridge.go toParent.ident=%s toIdent=%s\n", toParent.GetIdent(), toIdent)
	to := yang.FindByIdent2(toParent, toIdent)
	return to, nil
}

func (bb *BridgeBrowser) RootSelector() (*browse.Selection, error) {
	root, err := bb.To.RootSelector()
	if err != nil {
		return nil, err
	}
	bridge := &Bridge{Map: bb.Map}
	return bridge.enterBridge(root)
}

func (b *Bridge) enterBridge(to *browse.Selection) (*browse.Selection, error) {
	s := &browse.Selection{}
	s.Enter = func() (child *browse.Selection, err error) {
		if to.Position, err = b.Map.MapMeta(s.Position, to.Meta); err == nil {
			var toChild *browse.Selection
			if toChild, err = to.Enter(); err == nil {
				s.Found = to.Found
				if toChild != nil {
					toChild.Meta = to.Position.(yang.MetaList)
					return b.enterBridge(toChild)
				}
			}
		}
		return
	}
	s.Edit = func(op browse.Operation, val *browse.Value) (err error) {
		if to.Position, err = b.Map.MapMeta(s.Position, to.Meta); err == nil {
			return to.Edit(op, val)
		}
		return
	}
	s.ReadValue = func(val *browse.Value) (err error) {
		if to.Position, err = b.Map.MapMeta(s.Position, to.Meta); err == nil {
			// TODO: txlate val
			return to.Edit(browse.UPDATE_VALUE, val)
		}
		return
	}
	return s, nil
}

func (b *BridgeBrowser) Module() *yang.Module {
	return b.From
}



