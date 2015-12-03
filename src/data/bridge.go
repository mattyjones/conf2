package data

import (
	"schema"
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

func (b *Bridge) Node() (Node) {
	return b.selectBridge(b.internal, b.Mapping)
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

//func (b *Bridge) internalPath(externalPath *schema.Path, internalMeta schema.MetaList) *schema.Path {
//	mapping := b.Mapping
//	var found bool
//	internalPath := &schema.Path{
//		Info : externalPath.Info,
//		Meta: internalMeta,
//	}
//	i := internalPath
//	xNext := externalPath
//	for xNext != nil {
//		iNext := &schema.Path{
//			Info : i.Info,
//			Parent: i,
//		}
//		xIdent := xNext.Meta.GetIdent()
//		mapping, found = mapping.Children[xIdent]
//		if !found {
//			panic(xIdent + " path unmapped")
//		}
//		iNext.Meta = schema.FindByIdentExpandChoices(i.Meta.(schema.MetaList), mapping.InternalIdent)
//		if iNext.Meta == nil {
//			panic(mapping.InternalIdent + " not found in schema")
//		}
//		// TODO : txlate keys
//		i.Next.Key = xNext.Key
//
//		xNext = xNext.Next
//		i = i.Next
//	}
//	return internalPath
//}

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
	s.OnWrite = func(state *Selection, externalMeta schema.HasDataType, val *schema.Value) error {
		if _, ok := b.updateInternalPosition(externalMeta, internal, mapping); ok {
			return internal.Node.Write(internal, internal.State.Position().(schema.HasDataType), val)
		}
		return nil
	}
	s.OnRead = func(state *Selection, externalMeta schema.HasDataType) (*schema.Value, error) {
		if _, ok := b.updateInternalPosition(externalMeta, internal, mapping); ok {
			// TODO: translate val
			return internal.Node.Read(internal, internal.State.Position().(schema.HasDataType))
		}
		return nil, nil
	}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, key []*schema.Value, first bool) (next Node, err error) {
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
