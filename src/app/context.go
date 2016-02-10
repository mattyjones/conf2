package app

import (
	"data"
	"schema"
)

type Context struct {
	Nodes map[string]data.Node
	Module *schema.Module
}

func (self *Context) Select() *data.Selection {
	n := &data.MyNode{}
	n.OnSelect = func(sel *data.Selection, module schema.MetaList, new bool) (data.Node, error) {
		switch module.GetIdent() {
		case "modules":
			return self.manageModules(), nil
		}
		return self.Nodes[module.GetIdent()], nil
	}
	return data.NewSelection(self.Schema(), n)
}

func (self *Context) manageModules() data.Node {
	n := &data.MyNode{}
	index := data.NewIndex(self.Nodes)
	n.OnNext = func(sel *data.Selection, meta *schema.List, new bool, key []*data.Value, first bool) (data.Node, error) {
		var m *schema.Module
		if len(key) > 0 {
			m = schema.FindByIdent2(self.Module.DataDefs(), key[0].Str).(*schema.Module)
		} else {
			if kval := index.NextKey(first); kval.IsValid() {
				if candidate := schema.FindByIdent2(self.Module.DataDefs(), kval.String()); candidate != nil {
					m = candidate.(*schema.Module)
				}
			}
		}
		if m != nil {
			sel.SetKey(data.SetValues(meta.KeyMeta(), m.GetIdent()))
			return data.NewSchemaData(m, false).SelectModule(m), nil
		}
		return nil, nil
	}
	return n
}

func (self *Context) Schema() *schema.Module {
	return self.Module
}

func (self *Context) Register(module *schema.Module, root data.Node) {
	if self.Nodes == nil {
		self.Nodes = make(map[string]data.Node)
	}
	self.Nodes[module.GetIdent()] = root
	self.Module.AddMeta(module)
}
