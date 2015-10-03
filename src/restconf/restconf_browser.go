package restconf

import (
	"schema"
	"schema/yang"
	"schema/browse"
)

type RestconfBrowser struct {
	Service *serviceImpl
	Meta *schema.Module
}

func NewBrowser(restconf *serviceImpl) (rcb *RestconfBrowser, err error) {
	var module *schema.Module
	module, err = yang.LoadModuleFromByteArray([]byte(restconfYang), nil)
	if err == nil {
		parent := schema.FindByPath(module, "modules").(*schema.List)
		placeholder := schema.FindByPath(parent, "module")
		targetParent := browse.GetSchemaSchema()
		targetMaster := schema.FindByPath(targetParent, "module").(*schema.Container)
		// shallow clone target otherwise we alter browser's schema
		target := *targetMaster
		parent.ReplaceMeta(placeholder, &target);
		rcb = &RestconfBrowser{Meta:module, Service:restconf}
	}
	return
}

func (rcb *RestconfBrowser) Module() (*schema.Module) {
	return rcb.Meta
}

func (rcb *RestconfBrowser) Close() error {
	return nil
}

func (rcb *RestconfBrowser) RootSelector() (browse.Selection, *browse.WalkState, error) {
	s := &browse.MySelection{}
	s.OnSelect = func (state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
		switch meta.GetIdent() {
			case "modules":
				return enterRegistrations(rcb.Service.registrations)
		}
		return nil, nil
	}
	return s, browse.NewWalkState(rcb.Meta), nil
}

func enterRegistrations(registrations map[string]*registration) (browse.Selection, error) {
	s := &browse.MySelection{}
	index := newRegIndex(registrations)
	s.OnNext = index.Index.OnNext
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
		switch meta.GetIdent() {
		case "module":
			if index.Selected != nil {
				browser := browse.NewSchemaBrowser(index.Selected.browser.Module(), true)
				return browser.SelectModule(index.Selected.browser.Module())
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
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

var restconfYang = `
module restconf {
	namespace "http://org.conf2/ns/modules";
	prefix "modules";
	revision 0000-00-00 {
		description "initial ver";
	}
	list modules {
		key "name";
		leaf name {
			type string;
		}
		container module {
			/* replace with YANG-1.0 meta */
			leaf nop {
				type string;
			}
		}
	}
}`