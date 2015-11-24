package data

import (
	"fmt"
	"schema"
	"strings"
)

type Bridge struct {
	internal *Selection
	path     string
	external schema.MetaList
	Mapping  *BridgeMapping
}

func NewBridge(internal *Selection, external schema.MetaList) *Bridge {
	bridge := &Bridge{
		internal: internal,
		external: external,
		Mapping:  NewBridgeMapping(external.GetIdent()),
	}
	return bridge
}

func (b *Bridge) Selector(externalPath *Path) (s *Selection, err error) {
	root := b.selectBridge(b.internal, b.Mapping)
	return WalkPath(NewSelection(root, b.external), externalPath)
}

type BridgeMapping struct {
	InternalIdent string
	Children      map[string]*BridgeMapping
}

func (m *BridgeMapping) AddMapping(externalIdent string, internalIdent string) *BridgeMapping {
	mapping := NewBridgeMapping(internalIdent)
	m.Children[externalIdent] = mapping
	return mapping
}

func NewBridgeMapping(internalIdent string) *BridgeMapping {
	return &BridgeMapping{
		InternalIdent: internalIdent,
		Children:      make(map[string]*BridgeMapping, 0),
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

func (b *Bridge) internalPath(p *Path, meta schema.Meta) *Path {
	mapping := b.Mapping
	var found bool
	internalPath := make([]string, len(p.Segments))
	m := meta
	for i, seg := range p.Segments {
		mapping, found = mapping.Children[seg.Ident]
		if !found {
			panic("path unmappable")
		}
		internalPath[i] = mapping.InternalIdent
		if len(seg.Keys) > 0 {
			internalPath[i] = fmt.Sprint(internalPath[i], "=", seg.Keys[0])
		}
		m = schema.FindByIdent2(m.(schema.MetaList), seg.Ident)
		if m == nil {
			panic("Mapping invalid")
		}
	}
	return NewPath(strings.Join(internalPath, "/"))
}

func (b *Bridge) updateInternalPosition(externalMeta schema.Meta, internalState *Selection, mapping *BridgeMapping) (*BridgeMapping, bool) {
	var childMapping *BridgeMapping
	var internalPosition schema.Meta
	if internalPosition, childMapping = mapping.SelectMap(externalMeta, internalState.State.SelectedMeta()); internalPosition != nil {
		internalState.State.SetPosition(internalPosition)
		return childMapping, true
	}
	return nil, false
}

func (b *Bridge) selectBridge(internal *Selection, mapping *BridgeMapping) Node {
	s := &MyNode{OnEvent: internal.Node.Event}
	s.OnSelect = func(state *Selection, externalMeta schema.MetaList, new bool) (child Node, err error) {
		if childMapping, ok := b.updateInternalPosition(externalMeta, internal, mapping); ok {
			var internalChild Node
			if internalChild, err = internal.Node.Select(internal, internal.State.Position().(schema.MetaList), new); err != nil {
				return nil, err
			} else if internalChild == nil {
				return nil, nil
			}
			return b.selectBridge(internal.Select(internalChild), childMapping), nil
		}
		return
	}
	s.OnWrite = func(state *Selection, externalMeta schema.HasDataType, val *Value) error {
		if _, ok := b.updateInternalPosition(externalMeta, internal, mapping); ok {
			return internal.Node.Write(internal, internal.State.Position().(schema.HasDataType), val)
		}
		return nil
	}
	s.OnRead = func(state *Selection, externalMeta schema.HasDataType) (*Value, error) {
		if _, ok := b.updateInternalPosition(externalMeta, internal, mapping); ok {
			// TODO: translate val
			return internal.Node.Read(internal, internal.State.Position().(schema.HasDataType))
		}
		return nil, nil
	}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		var internalNextNode Node
		// TODO: translate keys?
		internalNextNode, err = internal.Node.Next(internal, meta, new, key, first)
		if internalNextNode != nil && err == nil {
			internalNext := internal.SelectListItem(internalNextNode, internal.State.Key())
			next = b.selectBridge(internalNext, mapping)
		}
		return
	}
	return s
}

func (b *Bridge) Schema() schema.MetaList {
	return b.external
}
