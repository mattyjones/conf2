package adapt
import (
	"schema"
	"schema/browse"
	"errors"
	"fmt"
)

type Bridge struct {
	Actual browse.Browser
	Emulate *schema.Module
	Mapping *MetaListMapping
}

type BridgeMapping interface {
	ToIdent() string
}

type MetaMapping struct {
	To string
}

func (mm *MetaMapping) ToIdent() string {
	return mm.To
}

type MetaListMapping struct {
	To string
	Mapping map[string]BridgeMapping
}

func (mlm *MetaListMapping) ToIdent() string {
	return mlm.To
}

func (m *MetaListMapping) AddMetaMapping(from string, to string) *MetaMapping {
	mapping := &MetaMapping{To : to}
	m.Mapping[from] = mapping
	return mapping
}

func (m *MetaListMapping) AddMetaListMapping(from string, to string) *MetaListMapping {
	mapping := NewMetaListMapping(to)
	m.Mapping[from] = mapping
	return mapping
}

func NewMetaListMapping(to string) *MetaListMapping {
	return &MetaListMapping{
		To: to,
		Mapping : make(map[string]BridgeMapping, 10),
	}
}

func (m *MetaListMapping) MapMetaList(from schema.Meta, toParent schema.MetaList) (schema.Meta, *MetaListMapping, error) {
	if from == nil {
		return nil, nil, nil
	}
	ident := from.GetIdent()
	var listMapping *MetaListMapping
	if m != nil {
		if mapping, found := m.Mapping[ident]; found {
			ident = mapping.ToIdent()
			listMapping = mapping.(*MetaListMapping)
		}
	}
	to := schema.FindByIdent2(toParent, ident)
	if to == nil {
		return nil, nil, errors.New(fmt.Sprint("No meta list mapping found for ", ident))
	}
	return to, listMapping, nil
}

func (m *MetaListMapping) MapMeta(from schema.Meta, toParent schema.MetaList) (schema.Meta, error) {
	if from == nil {
		return nil, nil
	}
	ident := from.GetIdent()
	if m != nil {
		if mapping, found := m.Mapping[ident]; found {
			ident = mapping.ToIdent()
		}
	}
	to := schema.FindByIdent2(toParent, ident)
	if to == nil {
		return nil, errors.New(fmt.Sprint("No meta mapping found for ", ident))
	}
	return to, nil
}

func (b *Bridge) RootSelector() (browse.Selection, error) {
	root, err := b.Actual.RootSelector()
	if err != nil {
		return nil, err
	}
	return b.selectBridge(root, b.Mapping)
}

func (b *Bridge) selectBridge(to browse.Selection, mapping *MetaListMapping) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func() (child browse.Selection, err error) {
		toState := to.WalkState()
		var childMapping *MetaListMapping
		if toState.Position, childMapping, err = mapping.MapMetaList(s.State.Position, toState.Meta); err == nil {
			var toChild browse.Selection
			if toChild, err = to.Select(); err == nil {
				s.WalkState().Found = to.WalkState().Found
				if toChild != nil {
					toChild.WalkState().Meta = toState.Position.(schema.MetaList)
					return b.selectBridge(toChild, childMapping)
				}
			}
		}
		return
	}
	s.OnWrite = func(op browse.Operation, val *browse.Value) (err error) {
		toState := to.WalkState()
		if toState.Position, err = mapping.MapMeta(s.State.Position, toState.Meta); err == nil {
			return to.Write(op, val)
		}
		return
	}
	s.OnRead = func(val *browse.Value) (err error) {
		toState := to.WalkState()
		if toState.Position, err = mapping.MapMeta(s.State.Position, toState.Meta); err == nil {
			// TODO: txlate val
			return to.Write(browse.UPDATE_VALUE, val)
		}
		return
	}
	s.OnNext = func(key []browse.Value, first bool) (bool, error) {
		return to.Next(key, first)
	}
	return s, nil
}

func (b *Bridge) Module() *schema.Module {
	return b.Emulate
}



