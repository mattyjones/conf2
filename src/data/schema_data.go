package data

import (
	"fmt"
	"schema"
	"schema/yang"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  Note: It can be confusing which is the data and which is the
 * meta.
 */
type SchemaData struct {
	data *schema.Module // read: data
	meta *schema.Module // read: meta-data

	// resolve all uses, groups and typedefs.  if this is false, then depth must be
	// used to avoid infinite recursion
	resolve bool
}

func (self *SchemaData) Schema() schema.MetaList {
	return self.meta
}

func NewSchemaData(data *schema.Module, resolve bool) *SchemaData {
	// TODO: Not require data to be a module
	browser := &SchemaData{data: data, meta: GetSchemaSchema(), resolve: resolve}
	return browser
}

var yang1_0 *schema.Module

func GetSchemaSchema() *schema.Module {
	if yang1_0 == nil {
		var err error
		yang1_0, err = yang.LoadModule(yang.YangPath(), "yang.yang")
		if err != nil {
			msg := fmt.Sprintf("Error parsing yang-1.0 yang, %s", err.Error())
			panic(msg)
		}
	}
	return yang1_0
}

type MetaListSelector func(m schema.Meta) (Node, error)

func (self *SchemaData) Node() Node {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		switch meta.GetIdent() {
		case "module":
			if new {
				self.data = &schema.Module{}
			}
			if self.data != nil {
				return self.SelectModule(self.data), nil
			}
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) SelectModule(module *schema.Module) (Node) {
	s := &MyNode{}
	delegate := self.selectMetaList(module)
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (child Node, err error) {
		switch meta.GetIdent() {
		case "revision":
			if new {
				module.Revision = &schema.Revision{}
			}
			if module.Revision != nil {
				return self.selectRevision(module.Revision), nil
			}
		case "notifications":
			return self.selectNotifications(module.GetNotifications()), nil
		default:
			return delegate.Select(state, meta, new)
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, module)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, module, val)
	}
	return s
}

func (self *SchemaData) selectRevision(rev *schema.Revision) (Node) {
	s := &MyNode{}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		switch meta.GetIdent() {
		case "rev-date":
			return &schema.Value{Str: rev.Ident}, nil
		default:
			return ReadField(meta, rev)
		}
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		switch meta.GetIdent() {
		case "rev-date":
			rev.Ident = val.Str
		default:
			return WriteField(meta, rev, val)
		}
		return nil
	}
	return s
}

func (self *SchemaData) selectType(typeData *schema.DataType) (Node) {
	s := &MyNode{}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, typeData)
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType,val *schema.Value) error {
		return WriteField(meta.(schema.HasDataType), typeData, val)
	}
	return s
}

func (self *SchemaData) selectGroupings(groupings schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: groupings, resolve: self.resolve}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (Node, error) {
		var group *schema.Grouping
		if new {
			group = &schema.Grouping{}
		} else {
			if i.iterate(sel, meta, keys, first) {
				group = i.data.(*schema.Grouping)
			}
		}
		if group != nil {
			sel.OnChild(NEW, meta, func() error {
				groupings.AddMeta(group)
				return nil
			})
			return self.selectMetaList(group), nil
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) selectRpcIO(i *schema.RpcInput, o *schema.RpcOutput) (Node) {
	var io schema.MetaList
	if i != nil {
		io = i
	} else {
		io = o
	}
	return self.selectMetaList(io)
}

func (self *SchemaData) createGroupingsTypedefsDefinitions(parent schema.MetaList, childMeta schema.Meta) (schema.Meta) {
	var child schema.Meta
	switch childMeta.GetIdent() {
	case "leaf":
		child = &schema.Leaf{}
	case "leaf-list":
		child = &schema.LeafList{}
	case "container":
		child = &schema.Container{}
	case "list":
		child = &schema.List{}
	case "uses":
		child = &schema.Uses{}
	case "grouping":
		child = &schema.Grouping{}
	case "typedef":
		child = &schema.Typedef{}
	case "rpc":
		child = &schema.Rpc{}
	default:
		panic("Unknown type")
	}
	parent.AddMeta(child)
	return child
}

func (self *SchemaData) selectRpc(rpc *schema.Rpc) (Node) {
	s := &MyNode{}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
		switch meta.GetIdent() {
		case "input":
			if new {
				rpc.AddMeta(&schema.RpcInput{})
			}
			if rpc.Input != nil {
				return self.selectRpcIO(rpc.Input, nil), nil
			}
		case "output":
			if new {
				rpc.AddMeta(&schema.RpcOutput{})
			}
			if rpc.Output != nil {
				return self.selectRpcIO(nil, rpc.Output), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, rpc)
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, rpc, val)
	}
	return s
}

func (self *SchemaData) selectTypedefs(typedefs schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: typedefs, resolve: self.resolve}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (Node, error) {
		var typedef *schema.Typedef
		if new {
			typedef = &schema.Typedef{}
		} else {
			if i.iterate(sel, meta, keys, first) {
				typedef = i.data.(*schema.Typedef)
			}
		}
		if typedef != nil {
			return self.selectTypedef(typedef), nil
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) selectTypedef(typedef *schema.Typedef) Node {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		switch meta.GetIdent() {
		case "type":
			if new {
				typedef.SetDataType(&schema.DataType{})
			}
			if typedef.DataType != nil {
				return self.selectType(typedef.DataType), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, typedef)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, typedef, val)
	}
	return s
}

func (self *SchemaData) selectMetaList(data schema.MetaList) (Node) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		hasGroupings, implementsHasGroupings := data.(schema.HasGroupings)
		hasTypedefs, implementsHasTypedefs := data.(schema.HasTypedefs)
		switch meta.GetIdent() {
		case "groupings":
			if !self.resolve && implementsHasGroupings {
				groupings := hasGroupings.GetGroupings()
				return self.selectGroupings(groupings), nil
			}
		case "typedefs":
			if !self.resolve && implementsHasTypedefs {
				typedefs := hasTypedefs.GetTypedefs()
				return self.selectTypedefs(typedefs), nil
			}
		case "definitions":
			defs := data.(schema.MetaList)
			return self.SelectDefinitionsList(defs), nil
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, data, val)
	}
	return s
}

func (self *SchemaData) selectNotifications(notifications schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: notifications, resolve: self.resolve}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (Node, error) {
		var notif *schema.Notification
		if new {
			notif = &schema.Notification{}
			notifications.AddMeta(notif)
		} else {
			if i.iterate(state, meta, keys, first) {
				notif = i.data.(*schema.Notification)
			}
		}
		if notif != nil {
			return self.selectMetaList(notif), nil
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) selectMetaLeafy(leaf *schema.Leaf, leafList *schema.LeafList) (Node) {
	s := &MyNode{}
	var leafy schema.HasDataType
	if leaf != nil {
		leafy = leaf
	} else {
		leafy = leafList
	}
	details := leafy.(schema.HasDetails).Details()
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		switch meta.GetIdent() {
		case "type":
			if new {
				leafy.SetDataType(&schema.DataType{})
			}
			if leafy.GetDataType() != nil {
				return self.selectType(leafy.GetDataType()), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		switch meta.GetIdent() {
		case "config":
			if details.ConfigFlag.IsSet() {
				return &schema.Value{Bool: details.ConfigFlag.Bool()}, nil
			}
		case "mandatory":
			if details.MandatoryFlag.IsSet() {
				return &schema.Value{Bool: details.MandatoryFlag.Bool()}, nil
			}
		default:
			return ReadField(meta, leafy)
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		switch meta.GetIdent() {
		case "config":
			details.ConfigFlag.Set(val.Bool)
		case "mandatory":
			details.MandatoryFlag.Set(val.Bool)
		default:
			return WriteField(meta, leafy, val)
		}
		return nil
	}
	return s
}

func (self *SchemaData) selectMetaUses(data *schema.Uses) (Node) {
	s := &MyNode{}
	// TODO: uses has refine container(s)
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, data, val)
	}
	return s
}

func (self *SchemaData) selectMetaCases(choice *schema.Choice) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: choice, resolve: self.resolve}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (Node, error) {
		var choiceCase *schema.ChoiceCase
		if new {
			choiceCase = &schema.ChoiceCase{}
			choice.AddMeta(choiceCase)
		} else {
			if i.iterate(state, meta, keys, first) {
				choiceCase = i.data.(*schema.ChoiceCase)
			}
		}
		if choiceCase != nil {
			return self.selectMetaList(choiceCase), nil
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) selectMetaChoice(data *schema.Choice) (Node) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		switch meta.GetIdent() {
		case "cases":
			// TODO: Not sure how to do create w/o what type to create
			return self.selectMetaCases(data), nil
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) error {
		return WriteField(meta, data, val)
	}
	return s
}

type listIterator struct {
	data     schema.Meta
	dataList schema.MetaList
	iterator schema.MetaIterator
	resolve  bool
	temp     int
}

func (i *listIterator) iterate(sel *Selection, meta *schema.List, keys []*schema.Value, first bool) bool {
	i.data = nil
	if i.dataList == nil {
		return false
	}
	if len(keys) > 0 {
		sel.State.SetKey(keys)
		if first {
			i.data = schema.FindByIdent2(i.dataList, keys[0].Str)
		}
	} else {
		if first {
			i.iterator = schema.NewMetaListIterator(i.dataList, i.resolve)
		}
		if i.iterator.HasNextMeta() {
			i.data = i.iterator.NextMeta()
			if i.data == nil {
				panic(fmt.Sprintf("Bad iterator at %s, item number %d", sel.String(), i.temp))
			}
			sel.State.SetKey([]*schema.Value{
				&schema.Value{
					Str:  i.data.GetIdent(),
					Type: &schema.DataType{Format: schema.FMT_STRING},
				},
			})
		}
		i.temp++
	}
	return i.data != nil
}

func (self *SchemaData) SelectDefinition(parent schema.MetaList, data schema.Meta) (Node) {
	s := &MyNode{}
	s.OnChoose = func(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, data)
	}
	s.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (Node, error) {
		if new {
			data = self.createGroupingsTypedefsDefinitions(parent, meta)
		}
		if data == nil {
			return nil, nil
		}
		switch meta.GetIdent() {
		case "leaf":
			return self.selectMetaLeafy(data.(*schema.Leaf), nil), nil
		case "leaf-list":
			return self.selectMetaLeafy(nil, data.(*schema.LeafList)), nil
		case "uses":
			return self.selectMetaUses(data.(*schema.Uses)), nil
		case "choice":
			return self.selectMetaChoice(data.(*schema.Choice)), nil
		case "rpc", "action":
			return self.selectRpc(data.(*schema.Rpc)), nil
		default:
			return self.selectMetaList(data.(schema.MetaList)), nil
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*schema.Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *schema.Value) (err error) {
		switch meta.GetIdent() {
		case "ident":
			// if data is nil then we're creating a def and we'll get name again

			if data != nil {
				return WriteField(meta, data, val)
			}
		default:
			return WriteField(meta, data, val)
		}
		//
		return nil
	}
	return s
}

func (self *SchemaData) SelectDefinitionsList(dataList schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: dataList, resolve: self.resolve}
	s.OnNext = func(state *Selection, meta *schema.List, new bool, keys []*schema.Value, first bool) (Node, error) {
		if new {
			return self.SelectDefinition(dataList, nil), nil
		} else {
			if i.iterate(state, meta, keys, first) {
				return self.SelectDefinition(dataList, i.data), nil
			}
		}
		return nil, nil
	}
	return s
}

func (self *SchemaData) resolveDefinitionCase(choice *schema.Choice, data schema.Meta) (caseMeta schema.MetaList, err error) {
	caseType := self.definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*schema.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg: msg}
	} else {
		return caseMeta, nil
	}
}

func (self *SchemaData) definitionType(data schema.Meta) string {
	switch data.(type) {
	case *schema.List:
		return "list"
	case *schema.Uses:
		return "uses"
	case *schema.Choice:
		return "choice"
	case *schema.Leaf:
		return "leaf"
	case *schema.LeafList:
		return "leaf-list"
	default:
		return "container"
	}
}
