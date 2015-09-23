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

func (b *Bridge) RootSelector() (browse.Selection, error) {
	internalRoot, err := b.internal.RootSelector()
	if err != nil {
		return nil, err
	}
	var bridged browse.Selection
	bridged, err = b.selectBridge(internalRoot, b.Mapping)
	bridged.WalkState().Meta = b.external
	return bridged, err
}

func (b *Bridge) selectBridge(internalSelection browse.Selection, mapping *BridgeMapping) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func(externalMeta schema.MetaList) (child browse.Selection, err error) {
		internalState := internalSelection.WalkState()
		var childMapping *BridgeMapping
		if internalState.Position, childMapping = mapping.SelectMap(externalMeta, internalState.Meta); internalState.Position != nil {
			var internalChild browse.Selection
			if internalChild, err = internalSelection.Select(internalState.Position.(schema.MetaList)); err == nil {
				if internalChild != nil {
					internalChild.WalkState().Meta = internalState.Position.(schema.MetaList)
					return b.selectBridge(internalChild, childMapping)
				}
			}
		}
		return
	}
	s.OnWrite = func(externalMeta schema.Meta, op browse.Operation, val *browse.Value) (err error) {
		internalState := internalSelection.WalkState()
		internalState.Position, _ = mapping.SelectMap(s.State.Position, internalState.Meta)
		if internalState.Position == nil && op == browse.UPDATE_VALUE {
			return nil
		}
		return internalSelection.Write(internalState.Position, op, val)
	}
	s.OnRead = func(externalMeta schema.HasDataType) (*browse.Value, error) {
		internalState := internalSelection.WalkState()
		internalPosition, _ := mapping.SelectMap(s.State.Position, internalState.Meta)
		if internalPosition != nil {
			internalSelection.WalkState().Position = internalPosition
			// TODO: translate val
			return internalSelection.Read(internalPosition.(schema.HasDataType))
		}
		return nil, nil
	}
	s.OnNext = func(key []*browse.Value, first bool) (bool, error) {
		// TODO: need to translate keys?
		return internalSelection.Next(key, first)
	}
	return s, nil
}

func (b *Bridge) Module() *schema.Module {
	return b.external
}



