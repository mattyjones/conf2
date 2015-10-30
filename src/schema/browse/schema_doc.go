package browse
import (
	"schema"
	"schema/yang"
	"fmt"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  Note: It can be confusing which is the data and which is the
 * meta.
 */
type SchemaDoc struct {
	data *schema.Module  // read: data
	meta *schema.Module    // read: meta-data

	// resolve all uses, groups and typedefs.  if this is false, then depth must be
	// used to avoid infinite recursion
	resolve bool
}

func (self *SchemaDoc) Schema() schema.MetaList {
	return self.meta
}

func NewSchemaBrowser(data *schema.Module, resolve bool) *SchemaDoc {
	// TODO: Not require data to be a module
	browser := &SchemaDoc{data:data, meta:GetSchemaSchema(), resolve:resolve}
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

func (self *SchemaDoc) Selector(p *Path) (*Selection, error) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		switch meta.GetIdent() {
		case "module" :
			if self.data != nil {
				return self.SelectModule(self.data)
			}
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
			case CREATE_CONTAINER:
				self.data = &schema.Module{}
		}
		return nil
	}
	return WalkPath(NewSelection(s, self.meta), p)
}

func (self *SchemaDoc) SelectModule(module *schema.Module) (Node, error) {
	s := &MyNode{}
	delegate, _ := self.selectMetaList(module)
	s.OnSelect = func(state *Selection, meta schema.MetaList) (child Node, err error) {
		switch meta.GetIdent() {
		case "revision":
			if module.Revision != nil {
				return self.selectRevision(module.Revision)
			}
		case "notifications":
			return self.selectNotifications(module.GetNotifications())
		default:
			return delegate.Select(state, meta)
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, module)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
			case CREATE_CONTAINER:
				switch meta.GetIdent() {
					case "revision":
						module.Revision = &schema.Revision{}
					default:
						return EditNotImplemented(meta)
				}
			case UPDATE_VALUE:
				return WriteField(meta.(schema.HasDataType), module, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectRevision(rev *schema.Revision) (Node, error) {
	s := &MyNode{}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		switch meta.GetIdent() {
		case "rev-date":
			return &Value{Str:rev.Ident}, nil
		default:
			return ReadField(meta, rev)
		}
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
			case UPDATE_VALUE:
				switch meta.GetIdent() {
					case "rev-date":
						rev.Ident = val.Str
					default:
						return WriteField(meta.(schema.HasDataType), rev, val)
				}
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectType(typeData *schema.DataType) (Node, error) {
	s := &MyNode{}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, typeData)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(meta.(schema.HasDataType), typeData, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectGroupings(groupings schema.MetaList) (Node, error) {
	s := &MyNode{}
	i := listIterator{dataList:groupings, resolve:self.resolve}
	var created *schema.Grouping
	s.OnNext = func(state *Selection, meta *schema.List, keys []*Value, first bool) (Node, error) {
		if created != nil {
			return self.selectMetaList(created)
		}
		if i.iterate(state, meta, keys, first) {
			return self.selectMetaList(i.data.(schema.MetaList))
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			created = &schema.Grouping{}
		case POST_CREATE_LIST_ITEM:
			groupings.AddMeta(created)
			created = nil
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectRpcIO(i *schema.RpcInput, o *schema.RpcOutput) (Node, error) {
	var io schema.MetaList
	if i != nil {
		io = i
	} else {
		io = o
	}
	return self.selectMetaList(io);
}

func (self *SchemaDoc) createGroupingsTypedefsDefinitions(parent schema.MetaList, childMeta schema.Meta) (schema.Meta, error) {
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
			return nil, NotImplemented(childMeta)
	}
	parent.AddMeta(child)
	return child, nil
}

func (self *SchemaDoc) selectRpc(rpc *schema.Rpc) (Node, error) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		switch meta.GetIdent() {
		case "input":
			if rpc.Input != nil {
				return self.selectRpcIO(rpc.Input, nil)
			}
		case "output":
			if rpc.Output != nil {
				return self.selectRpcIO(nil, rpc.Output)
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, rpc)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
			case CREATE_CONTAINER:
				switch meta.GetIdent() {
				case "input":
					rpc.AddMeta(&schema.RpcInput{})
				case "output":
					rpc.AddMeta(&schema.RpcOutput{})
				}
			case UPDATE_VALUE:
				return WriteField(meta.(schema.HasDataType), rpc, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectTypedefs(typedefs schema.MetaList) (Node, error) {
	s := &MyNode{}
	i := listIterator{dataList:typedefs, resolve:self.resolve}
	var created *schema.Typedef
	var selected *schema.Typedef
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		if selected == nil {
			return nil, nil
		}
		switch meta.GetIdent() {
		case "type":
			if selected.DataType != nil {
				return self.selectType(selected.DataType)
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, selected)
	}
	s.OnNext = func(state *Selection, meta *schema.List, keys []*Value, first bool) (Node, error) {
		if created != nil {
			selected = created
			return s, nil
		}
		if i.iterate(state, meta, keys, first) {
			selected = i.data.(*schema.Typedef)
			return s, nil
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_CONTAINER:
			switch meta.GetIdent() {
			case "type":
				selected.SetDataType(&schema.DataType{})
			default:
				return NotImplemented(meta)
			}
		case CREATE_LIST_ITEM:
			created = &schema.Typedef{}
		case POST_CREATE_LIST_ITEM:
			typedefs.AddMeta(created)
			created = nil
		case UPDATE_VALUE:
			return WriteField(meta.(schema.HasDataType), selected, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectMetaList(data schema.MetaList) (Node, error) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		hasGroupings, implementsHasGroupings := data.(schema.HasGroupings)
		hasTypedefs, implementsHasTypedefs := data.(schema.HasTypedefs)
		switch meta.GetIdent() {
		case "groupings":
			if !self.resolve && implementsHasGroupings {
				groupings := hasGroupings.GetGroupings()
				return self.selectGroupings(groupings)
			}
		case "typedefs":
			if !self.resolve && implementsHasTypedefs {
				typedefs := hasTypedefs.GetTypedefs()
				return self.selectTypedefs(typedefs)
			}
		case "definitions":
			defs := data.(schema.MetaList)
			return self.SelectDefinitionsList(defs)
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(meta.(schema.HasDataType), data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectNotifications(notifications schema.MetaList) (Node, error) {
	s := &MyNode{}
	i := listIterator{dataList:notifications, resolve:self.resolve}
	var created *schema.Notification
	s.OnNext = func(state *Selection, meta *schema.List, keys []*Value, first bool) (Node, error) {
		if created != nil {
			return self.selectMetaList(created)
		}
		if i.iterate(state, meta, keys, first) {
			return self.selectMetaList(i.data.(schema.MetaList))
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			created = &schema.Notification{}
		case POST_CREATE_LIST_ITEM:
			notifications.AddMeta(created)
			created = nil
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectMetaLeafy(leaf *schema.Leaf, leafList *schema.LeafList) (Node, error) {
	s := &MyNode{}
	var leafy schema.HasDataType
	if leaf != nil {
		leafy = leaf
	} else {
		leafy = leafList
	}
	details := leafy.(schema.HasDetails).Details()
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		switch meta.GetIdent() {
		case "type":
			if leafy.GetDataType() != nil {
				return self.selectType(leafy.GetDataType())
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		switch meta.GetIdent() {
			case "config":
				if details.ConfigFlag.IsSet() {
					return &Value{Bool:details.ConfigFlag.Bool()}, nil
				}
			case "mandatory":
				if details.MandatoryFlag.IsSet() {
					return &Value{Bool:details.MandatoryFlag.Bool()}, nil
				}
			default:
				return ReadField(meta, leafy)
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_CONTAINER:
			switch meta.GetIdent() {
			case "type":
				leafy.SetDataType(&schema.DataType{})
			default:
				return NotImplemented(meta)
			}
		case UPDATE_VALUE:
			switch meta.GetIdent() {
			case "config":
				details.ConfigFlag.Set(val.Bool)
			case "mandatory":
				details.MandatoryFlag.Set(val.Bool)
			default:
				return WriteField(meta.(schema.HasDataType), leafy, val)
			}
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectMetaUses(data *schema.Uses) (Node, error) {
	s := &MyNode{}
	// TODO: uses has refine container(s)
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(meta.(schema.HasDataType), data, val)
		}
		return nil
	}
	return s, nil
}


func (self *SchemaDoc) selectMetaCases(choice *schema.Choice) (Node, error) {
	s := &MyNode{}
	i := listIterator{dataList:choice, resolve:self.resolve}
	var created *schema.ChoiceCase
	s.OnNext = func(state *Selection, meta *schema.List, keys []*Value, first bool) (Node, error) {
		if i.iterate(state, meta, keys, first) {
			return self.selectMetaList(i.data.(schema.MetaList))
		}
		return nil, nil
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			created = &schema.ChoiceCase{}
		case POST_CREATE_LIST_ITEM:
			choice.AddMeta(created)
			created = nil
		}
		return nil
	}
	return s, nil
}

func (self *SchemaDoc) selectMetaChoice(data *schema.Choice) (Node, error) {
	s := &MyNode{}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		switch meta.GetIdent() {
		case "cases":
			return self.selectMetaCases(data);
		}
		return nil, nil
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, data)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(meta.(schema.HasDataType), data, val)
		}
		return nil
	}
	return s, nil
}

type listIterator struct {
	data schema.Meta
	dataList schema.MetaList
	iterator schema.MetaIterator
	resolve bool
	temp int
}

func (i *listIterator) iterate(state *Selection, meta *schema.List, keys []*Value, first bool) (bool) {
	i.data = nil
	if i.dataList == nil {
		return false
	}
	if len(keys) > 0 {
		state.SetKey(keys)
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
				panic(fmt.Sprintf("Bad iterator at %s, item number %d", state.String(), i.temp))
			}
			state.SetKey([]*Value{
				&Value{
					Str:i.data.GetIdent(),
					Type:&schema.DataType{Format:schema.FMT_STRING},
				},
			})
		}
		i.temp++
	}
	return i.data != nil
}

func (self *SchemaDoc) SelectDefinition(parent schema.MetaList, data schema.Meta) (Node, error) {
	s := &MyNode{}
	selected := data
	s.OnChoose = func(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, selected)
	}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (Node, error) {
		if selected == nil {
			return nil, nil
		}
		switch meta.GetIdent() {
		case "leaf":
			return self.selectMetaLeafy(selected.(*schema.Leaf), nil)
		case "leaf-list":
			return self.selectMetaLeafy(nil, selected.(*schema.LeafList))
		case "uses":
			return self.selectMetaUses(selected.(*schema.Uses))
		case "choice":
			return self.selectMetaChoice(selected.(*schema.Choice))
		case "rpc", "action":
			return self.selectRpc(selected.(*schema.Rpc))
		default:
			return self.selectMetaList(selected.(schema.MetaList))
		}
	}
	s.OnRead = func(state *Selection, meta schema.HasDataType) (*Value, error) {
		return ReadField(meta, selected)
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) (err error) {
		switch op {
		case CREATE_CONTAINER:
			selected, err = self.createGroupingsTypedefsDefinitions(parent, meta)
		case POST_CREATE_CONTAINER:
			selected = nil
		}
		return err
	}
	return s, nil
}


func (self *SchemaDoc) SelectDefinitionsList(dataList schema.MetaList) (Node, error) {
	s := &MyNode{}
	i := listIterator{dataList:dataList, resolve:self.resolve}
	var selected Node
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, val *Value) (err error) {
		switch op {
		case CREATE_LIST_ITEM:
			selected, err = self.SelectDefinition(dataList, nil)
		case POST_CREATE_LIST_ITEM:
			selected = nil
		}
		return err
	}
	s.OnNext = func(state *Selection, meta *schema.List, keys []*Value, first bool) (Node, error) {
		if selected != nil {
			return selected, nil
		}
		if i.iterate(state, meta, keys, first) {
			return self.SelectDefinition(dataList, i.data)
		}
		return nil, nil
	}
	return s, nil
}

func (self *SchemaDoc) resolveDefinitionCase(choice *schema.Choice, data schema.Meta) (caseMeta schema.MetaList, err error) {
	caseType := self.definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*schema.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg:msg}
	} else {
		return caseMeta, nil
	}
}

func (self *SchemaDoc) definitionType(data schema.Meta) string {
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
