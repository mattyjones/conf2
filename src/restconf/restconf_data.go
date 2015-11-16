package restconf

import (
	"schema"
	"schema/browse"
	"schema/yang"
)

type Data struct {
	Service *Service
	Meta    *schema.Module
}

func NewData(restconf *Service) (rcb *Data, err error) {
	var module *schema.Module
	module, err = yang.LoadModule(yang.YangPath(), "restconf.yang")
	if err == nil {
		parent := schema.FindByPath(module, "modules").(*schema.List)
		placeholder := schema.FindByPath(parent, "module")
		targetParent := browse.GetSchemaSchema()
		targetMaster := schema.FindByPath(targetParent, "module").(*schema.Container)
		// shallow clone target otherwise we alter browser's schema
		target := *targetMaster
		parent.ReplaceMeta(placeholder, &target)
		rcb = &Data{Meta: module, Service: restconf}
	}
	return
}

func (rcb *Data) Schema() schema.MetaList {
	return rcb.Meta
}

func (rcb *Data) Selector(path *browse.Path) (*browse.Selection, error) {
	s := &browse.MyNode{}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (browse.Node, error) {
		switch meta.GetIdent() {
		case "modules":
			return SelectModules(rcb.Service.registrations), nil
		}
		return nil, nil
	}
	return browse.WalkPath(browse.NewSelection(s, rcb.Meta), path)
}

func SelectManagement(service *Service) browse.Node {
	s := &browse.MyNode{}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (child browse.Node, err error) {
		switch meta.GetIdent() {
		case "modules":
			return SelectModules(service.registrations), nil
		}
		return
	}
	s.OnWrite = func(sel *browse.Selection, meta schema.HasDataType, val *browse.Value) (err error) {
		switch meta.GetIdent() {
		case "docRoot":
			service.SetDocRoot(&schema.FileStreamSource{Root: val.Str})
		default:
			return browse.WriteField(meta, service, val)
		}
		return
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) (err error) {
		switch e {
		case browse.NEW:
			go service.Listen()
			var data *Data
			if data, err = NewData(service); err != nil {
				return err
			}
			err = service.RegisterBrowser(data)
		}
		return
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		default:
			return browse.ReadField(meta, service)
		}
	}
	return s
}

func SelectModule(name string, reg *registration) browse.Node {
	s := &browse.MyNode{}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, new bool) (browse.Node, error) {
		switch meta.GetIdent() {
		case "module":
			// TODO: support browsing schema at any point, not assume module
			data := reg.browser.Schema().(*schema.Module)
			browser := browse.NewSchemaData(data, true)
			return browser.SelectModule(data), nil
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "name":
			return &browse.Value{Str: name}, nil
		}
		return nil, nil
	}
	return s
}

func SelectModules(registrations map[string]*registration) (browse.Node) {
	s := &browse.MyNode{}
	index := newRegIndex(registrations)
	s.OnNext = func(state *browse.Selection, meta *schema.List, new bool, keys []*browse.Value, isFirst bool) (browse.Node, error) {
		if hasMore, err := index.Index.OnNext(state, meta, keys, isFirst); hasMore {
			return SelectModule(index.Index.CurrentKey(), index.Selected), err
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "name":
			return &browse.Value{Str: index.Index.CurrentKey()}, nil
		}
		return nil, nil
	}
	return s
}

type regIndex struct {
	Index    browse.StringIndex
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
