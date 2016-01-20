package data

import (
	"schema"
)

// Pulls defaults from yang after explicit values have been set
//
// Mostly deprecated, Upsert and Insert will automatically r/w default
// values from yang

type Parameters struct {
	Ignores  map[string]struct{}
	Collected map[string]*Value
}

func (p *Parameters) Record(ident string) {
	if p.Ignores == nil {
		p.Ignores = make(map[string]struct{})
	}
	p.Ignores[ident] = struct{}{}
}

func (p *Parameters) Collect(ident string, val *Value) {
	if p.Collected == nil {
		p.Collected = make(map[string]*Value)
	}
	p.Collected[ident] = val
}

func (p *Parameters) Finish(sel *Selection, node Node) (err error) {
	i := schema.NewMetaListIterator(sel.path.meta, true)
	for i.HasNextMeta() {
		m := i.NextMeta()
		if _, ignore := p.Ignores[m.GetIdent()]; ignore {
			continue
		}
		t, hasType := m.(schema.HasDataType)
		if ! hasType {
			continue
		}
		var v *Value
		var found bool
		v, found = p.Collected[m.GetIdent()]
		if !found {
			if ! t.GetDataType().HasDefault() {
				continue
			}
			v = &Value{Type:t.GetDataType()}
			if err = v.CoerseStrValue(t.GetDataType().Default()); err != nil {
				return err
			}
		}
		if err = node.Write(sel, t, v); err != nil {
			return err
		}
	}
	return nil
}
