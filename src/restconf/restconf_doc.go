package restconf

import (
	"schema"
	"schema/yang"
	"schema/browse"
)

type Document struct {
	Service *Service
	Meta *schema.Module
}

func NewDoc(restconf *Service) (rcb *Document, err error) {
	var module *schema.Module
	module, err = yang.LoadModule(yang.YangPath(), "restconf.yang")
	if err == nil {
		parent := schema.FindByPath(module, "modules").(*schema.List)
		placeholder := schema.FindByPath(parent, "module")
		targetParent := browse.GetSchemaSchema()
		targetMaster := schema.FindByPath(targetParent, "module").(*schema.Container)
		// shallow clone target otherwise we alter browser's schema
		target := *targetMaster
		parent.ReplaceMeta(placeholder, &target);
		rcb = &Document{Meta:module, Service:restconf}
	}
	return
}

func (rcb *Document) Schema() (schema.MetaList) {
	return rcb.Meta
}

func (rcb *Document) Selector(path *browse.Path) (*browse.Selection, error) {
	s := &browse.MyNode{}
	s.OnSelect = func (state *browse.Selection, meta schema.MetaList) (browse.Node, error) {
		switch meta.GetIdent() {
			case "modules":
				return enterRegistrations(rcb.Service.registrations)
		}
		return nil, nil
	}
	return browse.WalkPath(browse.NewSelection(s, rcb.Meta), path)
}

func enterRegistrations(registrations map[string]*registration) (browse.Node, error) {
	s := &browse.MyNode{}
	index := newRegIndex(registrations)
	s.OnNext = func(state *browse.Selection, meta *schema.List, keys []*browse.Value, isFirst bool) (browse.Node, error) {
		if hasMore, err := index.Index.OnNext(state, meta, keys, isFirst); hasMore {
			return s, err
		}
		return nil, nil
	}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList) (browse.Node, error) {
		switch meta.GetIdent() {
		case "module":
			if index.Selected != nil {
				// TODO: support browsing schema at any point, not assume module
				data := index.Selected.browser.Schema().(*schema.Module)
				browser := browse.NewSchemaBrowser(data, true)
				return browser.SelectModule(data)
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
			case "name":
				return &browse.Value{Str : index.Index.CurrentKey()}, nil
		}
		return nil, nil
	}
	return s, nil
}

type regIndex struct {
	Index browse.StringIndex
	Selected *registration
	Data map[string]*registration
}

func newRegIndex(registrations map[string]*registration) *regIndex {
	ndx := &regIndex{Data:registrations}
	ndx.Index.Builder = ndx
	return ndx
}

func (ndx *regIndex) Select(key string) (found bool) {
	ndx.Selected, found = ndx.Data[key]
	return
}

func (ndx *regIndex) Build() []string {
	names := make([]string, len(ndx.Data))
	var j int
	for name, _ := range ndx.Data {
		names[j] = name
		j++
	}
	return names
}
