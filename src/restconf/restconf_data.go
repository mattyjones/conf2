package restconf

import (
	"schema"
	"data"
	"schema/yang"
)

type Data struct {
	Service *Service
	Meta    *schema.Module
}

func NewData(restconf *Service) (rcb *Data, err error) {
	var module *schema.Module
	module, err = yang.LoadModule(yang.YangPath(), "restconf.yang")
	if err != nil {
		return nil, err
	}
	parent := schema.FindByPath(module, "modules").(*schema.List)
	placeholder := schema.FindByPath(parent, "module")
	targetParent := data.GetSchemaSchema()
	targetMaster := schema.FindByPath(targetParent, "module").(*schema.Container)
	// shallow clone target otherwise we alter browser's schema
	target := *targetMaster
	parent.ReplaceMeta(placeholder, &target)
	rcb = &Data{Meta: module, Service: restconf}
	return
}

func (rcb *Data) Schema() schema.MetaList {
	return rcb.Meta
}

func (rcb *Data) Node() data.Node {
	return SelectManagement(rcb.Service)
}

func SelectManagement(service *Service) data.Node {
	s := &data.MyNode{}
	s.OnSelect = func(sel *data.Selection, meta schema.MetaList, new bool) (child data.Node, err error) {
		switch meta.GetIdent() {
		case "modules":
			return SelectModules(service.registrations), nil
		}
		return
	}
	s.OnWrite = func(sel *data.Selection, meta schema.HasDataType, val *schema.Value) (err error) {
		switch meta.GetIdent() {
		case "docRoot":
			service.SetDocRoot(&schema.FileStreamSource{Root: val.Str})
		default:
			return data.WriteField(meta, service, val)
		}
		return
	}
	s.OnEvent = func(sel *data.Selection, e data.Event) (err error) {
		switch e {
		case data.NEW:
			go service.Listen()
		}
		return
	}
	s.OnRead = func(sel *data.Selection, meta schema.HasDataType) (*schema.Value, error) {
		switch meta.GetIdent() {
		default:
			return data.ReadField(meta, service)
		}
	}
	return s
}

func SelectModule(name string, reg *registration) data.Node {
	s := &data.MyNode{}
	s.OnSelect = func(sel *data.Selection, meta schema.MetaList, new bool) (data.Node, error) {
		switch meta.GetIdent() {
		case "module":
			// TODO: support browsing schema at any point, not assume module
			d := reg.browser.Schema().(*schema.Module)
			browser := data.NewSchemaData(d, true)
			return browser.SelectModule(d), nil
		}
		return nil, nil
	}
	s.OnRead = func(sel *data.Selection, meta schema.HasDataType) (*schema.Value, error) {
		switch meta.GetIdent() {
		case "name":
			return &schema.Value{Str: name}, nil
		}
		return nil, nil
	}
	return s
}

func SelectModules(registrations map[string]*registration) (data.Node) {
	s := &data.MyNode{}
	index := newRegIndex(registrations)
	s.OnNext = func(sel *data.Selection, meta *schema.List, new bool, keys []*schema.Value, isFirst bool) (data.Node, error) {
		if hasMore, err := index.Index.OnNext(sel, meta, keys, isFirst); hasMore {
			return SelectModule(index.Index.CurrentKey(), index.Selected), err
		}
		return nil, nil
	}
	s.OnRead = func(sel *data.Selection, meta schema.HasDataType) (*schema.Value, error) {
		switch meta.GetIdent() {
		case "name":
			return &schema.Value{Str: index.Index.CurrentKey()}, nil
		}
		return nil, nil
	}
	return s
}

type regIndex struct {
	Index    data.StringIndex
	Selected *registration
	Data     map[string]*registration
}

func newRegIndex(registrations map[string]*registration) *regIndex {
	ndx := &regIndex{Data: registrations}
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
