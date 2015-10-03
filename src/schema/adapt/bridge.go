package adapt
import (
	"schema"
	"schema/browse"
)

type Bridge struct {
	internal browse.Browser
	path string
	external *schema.Module
	Mapping *BridgeMapping
}

func NewBridge(internal browse.Browser, external *schema.Module) *Bridge {
	bridge := &Bridge{
		internal: internal,
		external: external,
		Mapping: NewBridgeMapping(external.GetIdent()),
	}
	return bridge
}

type BridgeMapping struct {
	InternalIdent string
	Children map[string]*BridgeMapping
}

func (m *BridgeMapping) AddMapping(externalIdent string, internalIdent string) *BridgeMapping {
	mapping := NewBridgeMapping(internalIdent)
	m.Children[externalIdent] = mapping
	return mapping
}

func NewBridgeMapping(internalIdent string) *BridgeMapping {
	return &BridgeMapping{
		InternalIdent: internalIdent,
		Children : make(map[string]*BridgeMapping, 0),
	}
}

func (m *BridgeMapping) SelectMap(externalMeta schema.Meta, internalParentMeta schema.MetaList) (schema.Meta, *BridgeMapping) {
	if externalMeta == nil {
		return nil, nil
	}
	ident := externalMeta.GetIdent()
	var mapping *BridgeMapping
	var found bool
	if m != nil {
		if mapping, found = m.Children[ident]; found {
			ident = mapping.InternalIdent
		}
	}
	internalMeta := schema.FindByIdent2(internalParentMeta, ident)
	return internalMeta, mapping
}

func (b *Bridge) RootSelector() (browse.Selection, *browse.WalkState, error) {
	internalRoot, internalState, err := b.internal.RootSelector()
	if err != nil {
		return nil, nil, err
	}
	var bridged browse.Selection
	bridged, err = b.selectBridge(internalRoot, internalState, b.Mapping)
	return bridged, browse.NewWalkState(b.external), err
}

func (b *Bridge) updateInternalPosition(externalMeta schema.Meta, internalState *browse.WalkState, mapping *BridgeMapping) (*BridgeMapping, bool) {
	var childMapping *BridgeMapping
	var internalPosition schema.Meta
	if internalPosition, childMapping = mapping.SelectMap(externalMeta, internalState.SelectedMeta()); internalPosition != nil {
		internalState.SetPosition(internalPosition)
		return childMapping, true
	}
	return nil, false
}

func (b *Bridge) selectBridge(internalSelection browse.Selection, internalState *browse.WalkState, mapping *BridgeMapping) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func(state *browse.WalkState, externalMeta schema.MetaList) (child browse.Selection, err error) {
		if childMapping, ok := b.updateInternalPosition(externalMeta, internalState, mapping); ok {
			var internalChild browse.Selection
			if internalChild, err = internalSelection.Select(internalState, internalState.Position().(schema.MetaList)); err != nil {
				return nil, err
			} else if internalChild == nil {
				return nil, nil
			}
			return b.selectBridge(internalChild, internalState.Select(), childMapping)
		}
		return
	}
	s.OnWrite = func(state *browse.WalkState, externalMeta schema.Meta, op browse.Operation, val *browse.Value) error {
		if op == browse.BEGIN_EDIT || op == browse.END_EDIT {
			return internalSelection.Write(internalState, internalState.SelectedMeta(), op, val)
		}
		if _, ok := b.updateInternalPosition(externalMeta, internalState, mapping); ok {
			return internalSelection.Write(internalState, internalState.Position(), op, val)
		}
		return nil
	}
	s.OnRead = func(state *browse.WalkState, externalMeta schema.HasDataType) (*browse.Value, error) {
		if _, ok := b.updateInternalPosition(externalMeta, internalState, mapping); ok {
			// TODO: translate val
			return internalSelection.Read(internalState, internalState.Position().(schema.HasDataType))
		}
		return nil, nil
	}
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (bool, error) {
		// TODO: translate keys?
		return internalSelection.Next(internalState, meta, key, first)
	}
	return s, nil
}

func (b *Bridge) Module() *schema.Module {
	return b.external
}



