package restconf

import (
	"schema"
	"schema/yang"
	"schema/browse"
	"sort"
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

func (rcb *RestconfBrowser) RootSelector() (browse.Selection, error) {
	s := &browse.MySelection{}
	s.WalkState().Meta = rcb.Meta
	s.OnSelect = func () (browse.Selection, error) {
		ident := s.State.Position.GetIdent()
		switch ident {
			case "modules":
				s.WalkState().Found = len(rcb.Service.registrations) > 0
				return enterRegistrations(rcb.Service.registrations)
		}
		return nil, nil
	}
	return s, nil
}

func enterRegistrations(registrations map[string]registration) (browse.Selection, error) {
	var i int
	var names []string
	s := &browse.MySelection{}
	s.OnNext = func(keys []browse.Value, isFirst bool) (bool, error) {
		if isFirst {
			i = 0
			names = make([]string, len(registrations))
			var j int
			for name, _ := range registrations {
				names[j] = name
				j++
			}
			sort.Strings(names)
		} else {
			i++
		}
		return i < len(names), nil
	}
	s.OnSelect = func() (browse.Selection, error) {
		ident := s.State.Position.GetIdent()
		switch ident {
		case "module":
			var reg registration
			state := s.WalkState()
			reg, state.Found = registrations[names[i]]
			if state.Found {
				browser := browse.NewSchemaBrowser(reg.browser.Module(), true)
				return browser.SelectModule(reg.browser.Module())
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *browse.Value) error {
		ident := s.State.Position.GetIdent()
		switch ident {
			case "name":
				val.Type = s.State.Position.(schema.HasDataType).GetDataType()
				val.Str = names[i]
		}
		return nil
	}
	return s, nil
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