package restconf

import (
	"yang"
	"yang/browse"
	"sort"
)

type RestconfBrowser struct {
	Service *serviceImpl
	Meta *yang.Module
}

func NewBrowser(restconf *serviceImpl) (rcb *RestconfBrowser, err error) {
	var module *yang.Module
	module, err = yang.LoadModuleFromByteArray([]byte(restconfYang))
	if err == nil {
		parent := yang.FindByPath(module, "modules").(*yang.List)
		placeholder := yang.FindByPath(parent, "module")
		targetParent := browse.GetYangModule()
		targetMaster := yang.FindByPath(targetParent, "module").(*yang.Container)
		// shallow clone target otherwise we alter browser's schema
		target := *targetMaster
		parent.ReplaceMeta(placeholder, &target);
		rcb = &RestconfBrowser{Meta:module, Service:restconf}
	}
	return
}

func (rcb *RestconfBrowser) Module() (*yang.Module) {
	return rcb.Meta
}

func (rcb *RestconfBrowser) Close() error {
	return nil
}

func (rcb *RestconfBrowser) RootSelector() (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Meta = rcb.Meta
	s.Enter = func () (*browse.Selection, error) {
		ident := s.Position.GetIdent()
		switch ident {
			case "modules":
				s.Found = len(rcb.Service.registrations) > 0
				return enterRegistrations(rcb.Service.registrations)
		}
		return nil, nil
	}
	return
}

func enterRegistrations(registrations map[string]registration) (*browse.Selection, error) {
	var i int
	var names []string
	s := &browse.Selection{}
	s.Iterate = func(keys []string, isFirst bool) (bool, error) {
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
	s.Enter = func() (*browse.Selection, error) {
		ident := s.Position.GetIdent()
		switch ident {
		case "module":
			var reg registration
			reg, s.Found = registrations[names[i]]
			if s.Found {
				browser := browse.NewYangBrowser(reg.browser.Module(), true)
				return browser.SelectModule(reg.browser.Module())
			}
		}
		return nil, nil
	}
	s.ReadValue = func(val *browse.Value) error {
		ident := s.Position.GetIdent()
		switch ident {
			case "name":
				val.Type = s.Position.(yang.HasDataType).GetDataType()
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