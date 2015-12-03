package data

import (
	"schema"
)

type Parameters struct {
	Ignores  map[string]struct{}
	Collected map[string]*schema.Value
}

func (p *Parameters) Record(ident string) {
	if p.Ignores == nil {
		p.Ignores = make(map[string]struct{})
	}
	p.Ignores[ident] = struct{}{}
}

func (p *Parameters) Collect(ident string, val *schema.Value) {
	if p.Collected == nil {
		p.Collected = make(map[string]*schema.Value)
	}
	p.Collected[ident] = val
}

func (p *Parameters) Finish(sel *Selection, node Node) (err error) {
	i := schema.NewMetaListIterator(sel.State.SelectedMeta(), true)
	for i.HasNextMeta() {
		m := i.NextMeta()
		if _, ignore := p.Ignores[m.GetIdent()]; ignore {
			continue
		}
		t, hasType := m.(schema.HasDataType)
		if ! hasType {
			continue
		}
		var v *schema.Value
		var found bool
		v, found = p.Collected[m.GetIdent()]
		if !found {
			def := t.GetDataType().Default
			if len(def) == 0 {
				continue
			}
			v = &schema.Value{Type:t.GetDataType()}
			if err = v.CoerseStrValue(def); err != nil {
				return err
			}
		}
		if err = node.Write(sel, t, v); err != nil {
			return err
		}
	}
	return nil
}
