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

func (self *SchemaData) Select() *Selection {
	return NewSelection(self.meta, self.Node())
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
	s.OnSelect = func(state *Selection, r ContainerRequest) (Node, error) {
		switch r.Meta.GetIdent() {
		case "module":
			if r.New {
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
	return &Extend{
		Label:"Module",
		Node:self.selectMetaList(module),
		OnSelect : func(parent Node, sel *Selection, r ContainerRequest) (child Node, err error) {
			switch r.Meta.GetIdent() {
			case "revision":
				if r.New {
					module.Revision = &schema.Revision{}
				}
				if module.Revision != nil {
					return self.selectRevision(module.Revision), nil
				}
			case "notifications":
				return self.selectNotifications(module.GetNotifications()), nil
			default:
				return parent.Select(sel, r)
			}
			return nil, nil
		},
	}
}

func (self *SchemaData) selectRevision(rev *schema.Revision) (Node) {
	s := &MyNode{}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		switch meta.GetIdent() {
		case "rev-date":
			return &Value{Str: rev.Ident, Type:meta.GetDataType()}, nil
		default:
			return ReadField(meta, rev)
		}
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *Value) error {
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
	return &MyNode{
		OnRead: func(sel *Selection, meta schema.HasDataType) (*Value, error) {
			switch meta.GetIdent()  {
			case "ident":
				return SetValue(meta.GetDataType(), typeData.Ident)
			case "minLength":
				if self.resolve || typeData.MinLengthPtr != nil {
					return SetValue(meta.GetDataType(), typeData.MinLength())
				}
			case "maxLength":
				if self.resolve || typeData.MaxLengthPtr != nil {
					return SetValue(meta.GetDataType(), typeData.MaxLength())
				}
			case "path":
				if self.resolve || typeData.PathPtr != nil {
					return SetValue(meta.GetDataType(), typeData.Path())
				}
			case "enumeration":
				if self.resolve || len(typeData.EnumerationRef) > 0 {
					return SetValue(meta.GetDataType(), typeData.Enumeration())
				}
			}
			return nil, nil
		},
		OnWrite: func(state *Selection, meta schema.HasDataType, val *Value) error {
			switch meta.GetIdent() {
			case "ident":
				typeData.Ident = val.Str
				typeData.SetFormat(schema.DataTypeImplicitFormat(val.Str))
			case "minLength":
				typeData.SetMinLength(val.Int)
			case "maxLength":
				typeData.SetMaxLength(val.Int)
			case "path":
				typeData.SetPath(val.Str)
			case "enumeration":
				typeData.SetEnumeration(val.Strlist)
			}
			return nil
		},
	}
}

func (self *SchemaData) selectGroupings(groupings schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: groupings, resolve: self.resolve}
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		var key = r.Key
		var group *schema.Grouping
		if r.New {
			group = &schema.Grouping{Ident:r.Key[0].Str}
			groupings.AddMeta(group)
		} else {
			if i.iterate(sel, r.Meta, r.Key, r.First) {
				group = i.data.(*schema.Grouping)
				if len(key) == 0 {
					key = SetValues(r.Meta.KeyMeta(), group.Ident)
				}
			}
		}
		if group != nil {
			return self.selectMetaList(group), key, nil
		}
		return nil, nil, nil
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
	case "anyxml":
		child = &schema.Any{}
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
	return &Extend{
		Label:"rpc",
		Node: MarshalContainer(rpc),
		OnSelect: func(parent Node, sel *Selection, r ContainerRequest) (Node, error) {
			switch r.Meta.GetIdent() {
			case "input":
				if r.New {
					rpc.AddMeta(&schema.RpcInput{})
				}
				if rpc.Input != nil {
					return self.selectRpcIO(rpc.Input, nil), nil
				}
			case "output":
				if r.New {
					rpc.AddMeta(&schema.RpcOutput{})
				}
				if rpc.Output != nil {
					return self.selectRpcIO(nil, rpc.Output), nil
				}
			}
			return nil, nil
		},
	}
}

func (self *SchemaData) selectTypedefs(typedefs schema.MetaList) (Node) {
	s := &MyNode{}
	i := listIterator{dataList: typedefs, resolve: self.resolve}
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		var key = r.Key
		var typedef *schema.Typedef
		if r.New {
			typedef = &schema.Typedef{Ident:r.Key[0].Str}
			typedefs.AddMeta(typedef)
		} else {
			if i.iterate(sel, r.Meta, r.Key, r.First) {
				typedef = i.data.(*schema.Typedef)
			}
			if len(key) == 0 {
				key = SetValues(r.Meta.KeyMeta(), typedef.Ident)
			}
		}
		if typedef != nil {
			return self.selectTypedef(typedef), key, nil
		}
		return nil, nil, nil
	}
	return s
}

func (self *SchemaData) selectTypedef(typedef *schema.Typedef) Node {
	return &Extend{
		Label:"Typedef",
		Node: MarshalContainer(typedef),
		OnSelect :func(parent Node, sel *Selection, r ContainerRequest) (Node, error) {
			switch r.Meta.GetIdent() {
			case "type":
				if r.New {
					typedef.SetDataType(&schema.DataType{})
				}
				if typedef.DataType != nil {
					return self.selectType(typedef.DataType), nil
				}
			}
			return nil, nil
		},
	}
}

func (self *SchemaData) selectMetaList(data schema.MetaList) (Node) {
	return &Extend{
		Label: "MetaList",
		Node: MarshalContainer(data),
		OnSelect : func(parent Node, sel *Selection, r ContainerRequest) (Node, error) {
			hasGroupings, implementsHasGroupings := data.(schema.HasGroupings)
			hasTypedefs, implementsHasTypedefs := data.(schema.HasTypedefs)
			switch r.Meta.GetIdent() {
			case "groupings":
				if ! self.resolve && implementsHasGroupings {
					groupings := hasGroupings.GetGroupings()
					if r.New || ! schema.ListEmpty(groupings) {
						return self.selectGroupings(groupings), nil
					}
				}
			case "typedefs":
				if ! self.resolve && implementsHasTypedefs {
					typedefs := hasTypedefs.GetTypedefs()
					if r.New || ! schema.ListEmpty(typedefs) {
						return self.selectTypedefs(typedefs), nil
					}
				}
			case "definitions":
				defs := data.(schema.MetaList)
				if r.New || ! schema.ListEmpty(defs) {
					return self.SelectDefinitionsList(defs), nil
				}
			}
			return nil, nil
		},
	}
}

func (self *SchemaData) selectNotifications(notifications schema.MetaList) (Node) {
	s := &MyNode{
		Peekables: map[string]interface{} {"internal" : notifications},
	}
	i := listIterator{dataList: notifications, resolve: self.resolve}
	s.OnNext = func(state *Selection, r ListRequest) (Node, []*Value, error) {
		key := r.Key
		var notif *schema.Notification
		if r.New {
			notif = &schema.Notification{}
			notifications.AddMeta(notif)
		} else {
			if i.iterate(state, r.Meta, r.Key, r.First) {
				notif = i.data.(*schema.Notification)
				if len(key) == 0 {
					key = SetValues(r.Meta.KeyMeta(), notif.Ident)
				}
			}
		}
		if notif != nil {
			return self.selectMetaList(notif), key, nil
		}
		return nil, nil, nil
	}
	return s
}

func (self *SchemaData) selectMetaLeafy(leaf *schema.Leaf, leafList *schema.LeafList, any *schema.Any) (Node) {
	var leafy schema.HasDataType
	if leaf != nil {
		leafy = leaf
	} else if leafList != nil {
		leafy = leafList
	} else {
		leafy = any
	}
	s := &MyNode{
		Peekables: map[string]interface{} {"internal" : leafy},
	}
	details := leafy.(schema.HasDetails).Details()
	s.OnSelect = func(state *Selection, r ContainerRequest) (Node, error) {
		switch r.Meta.GetIdent() {
		case "type":
			if r.New {
				leafy.SetDataType(&schema.DataType{})
			}
			if leafy.GetDataType() != nil {
				return self.selectType(leafy.GetDataType()), nil
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		switch meta.GetIdent() {
		case "config":
			if self.resolve || details.ConfigPtr != nil {
				return &Value{Bool: details.Config(state.Path()), Type:meta.GetDataType()}, nil
			}
		case "mandatory":
			if self.resolve || details.MandatoryPtr != nil {
				return &Value{Bool: details.Mandatory(), Type:meta.GetDataType()}, nil
			}
		default:
			return ReadField(meta, leafy)
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *Value) error {
		switch meta.GetIdent() {
		case "config":
			details.SetConfig(val.Bool)
		case "mandatory":
			details.SetMandatory(val.Bool)
		default:
			return WriteField(meta, leafy, val)
		}
		return nil
	}
	return s
}

func (self *SchemaData) selectMetaUses(data *schema.Uses) (Node) {
	// TODO: uses has refine container(s)
	return MarshalContainer(data)
}

func (self *SchemaData) selectMetaCases(choice *schema.Choice) (Node) {
	s := &MyNode{
		Peekables: map[string]interface{} {"internal" : choice},
	}
	i := listIterator{dataList: choice, resolve: self.resolve}
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		key := r.Key
		var choiceCase *schema.ChoiceCase
		if r.New {
			choiceCase = &schema.ChoiceCase{}
			choice.AddMeta(choiceCase)
		} else {
			if i.iterate(sel, r.Meta, r.Key, r.First) {
				choiceCase = i.data.(*schema.ChoiceCase)
				if len(key) == 0 {
					key = SetValues(r.Meta.KeyMeta(), choiceCase.Ident)
				}
			}
		}
		if choiceCase != nil {
			return self.selectMetaList(choiceCase), key, nil
		}
		return nil, nil, nil
	}
	return s
}

func (self *SchemaData) selectMetaChoice(data *schema.Choice) (Node) {
	return &Extend{
		Label:"Choice",
		Node: MarshalContainer(data),
		OnSelect: func(parent Node, sel *Selection, r ContainerRequest) (Node, error) {
			switch r.Meta.GetIdent() {
			case "cases":
				// TODO: Not sure how to do create w/o what type to create
				return self.selectMetaCases(data), nil
			}
			return nil, nil
		},
	}
}

type listIterator struct {
	data     schema.Meta
	dataList schema.MetaList
	iterator schema.MetaIterator
	resolve  bool
	temp     int
}

func (i *listIterator) iterate(sel *Selection, meta *schema.List, key []*Value, first bool) bool {
	i.data = nil
	if i.dataList == nil {
		return false
	}
	if len(key) > 0 {
		sel.path.key = key
		if first {
			i.data = schema.FindByIdent2(i.dataList, key[0].Str)
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
			sel.path.key = SetValues(meta.KeyMeta(), i.data.GetIdent())
		}
		i.temp++
	}
	return i.data != nil
}

func (self *SchemaData) SelectDefinition(parent schema.MetaList, data schema.Meta) (Node) {
	s := &MyNode{
		Peekables: map[string]interface{} {"internal" : data},
	}
	s.OnChoose = func(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, data)
	}
	s.OnSelect = func(state *Selection, r ContainerRequest) (Node, error) {
		if r.New {
			data = self.createGroupingsTypedefsDefinitions(parent, r.Meta)
		}
		if data == nil {
			return nil, nil
		}
		switch r.Meta.GetIdent() {
		case "anyxml":
			return self.selectMetaLeafy(nil, nil, data.(*schema.Any)), nil
		case "leaf":
			return self.selectMetaLeafy(data.(*schema.Leaf), nil, nil), nil
		case "leaf-list":
			return self.selectMetaLeafy(nil, data.(*schema.LeafList), nil), nil
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
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.HasDataType, val *Value) (err error) {
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
	s := &MyNode{
		Peekables: map[string]interface{} {"internal" : dataList},
	}
	i := listIterator{dataList: dataList, resolve: self.resolve}
	s.OnNext = func(sel *Selection, r ListRequest) (Node, []*Value, error) {
		key := r.Key
		if r.New {
			return self.SelectDefinition(dataList, nil), key, nil
		} else {
			if i.iterate(sel, r.Meta, r.Key, r.First) {
				if len(key) == 0 {
					key = SetValues(r.Meta.KeyMeta(), i.data.GetIdent())
				}
				return self.SelectDefinition(dataList, i.data), key, nil
			}
		}
		return nil, nil, nil
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
	case *schema.Any:
		return "anyxml"
	case *schema.Leaf:
		return "leaf"
	case *schema.LeafList:
		return "leaf-list"
	default:
		return "container"
	}
}
