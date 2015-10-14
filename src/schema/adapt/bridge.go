package adapt
import (
	"schema"
	"schema/browse"
	"fmt"
	"strings"
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

func (b *Bridge) Selector(externalPath *browse.Path, strategy browse.Strategy) (browse.Selection, *browse.WalkState, error) {
	internalPath, externalState := b.internalPath(externalPath)
	internalRoot, internalState, err := b.internal.Selector(internalPath, strategy)
	if err != nil {
		return nil, nil, err
	}
	bridged, _ := b.selectBridge(internalRoot, internalState, b.Mapping)
	return bridged, externalState, nil
}

func (b *Bridge) internalPath(p *browse.Path) (*browse.Path, *browse.WalkState) {
	mapping := b.Mapping
	var found bool
	internalPath := make([]string, len(p.Segments))
	state := browse.NewWalkState(b.external)
	for i, seg := range p.Segments {
		mapping, found = mapping.Children[seg.Ident]
		if !found {
			panic("path unmappable")
		}
		internalPath[i] = mapping.InternalIdent
		if len(seg.Keys) > 0 {
			internalPath[i] = fmt.Sprint(internalPath[i], "=", seg.Keys[0])
		}
		position := schema.FindByIdent2(state.SelectedMeta(), seg.Ident)
		state.SetPosition(position)
		state = state.Select()
	}
	return browse.NewPath(strings.Join(internalPath, "/")), state
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
	if internalState == nil {
		panic("STOP")
	}
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



