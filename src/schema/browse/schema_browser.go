package browse
import (
	"schema"
	"schema/yang"
	"fmt"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type SchemaBrowser struct {
	module *schema.Module  // read: data
	meta *schema.Module    // read: meta-data

	// resolve all uses, groups and typedefs.  if this is false, then depth must be
	// used to avoid infinite recursion
	resolve bool
}

func (self *SchemaBrowser) Module() *schema.Module {
	return self.meta
}

func NewSchemaBrowser(module *schema.Module, resolve bool) *SchemaBrowser {
	browser := &SchemaBrowser{module:module, meta:GetSchemaSchema(), resolve:resolve}
	return browser
}

var yang1_0 *schema.Module
func GetSchemaSchema() *schema.Module {
	if yang1_0 == nil {
		var err error
		yang1_0, err = yang.LoadModuleFromByteArray([]byte(YANG_1_0), nil)
		if err != nil {
			msg := fmt.Sprintf("Error parsing yang-1.0 yang, %s", err.Error())
			panic(msg)
		}
	}
	return yang1_0
}

type MetaListSelector func(m schema.Meta) (Selection, error)

func (self *SchemaBrowser) RootSelector() (Selection, error) {
	s := &MySelection{}
	s.State.Meta = self.meta
	s.OnSelect = func() (Selection, error) {
		s.State.Found = self.module != nil
		switch s.State.Position.GetIdent() {
		case "module" :
			if self.module != nil {
				return self.SelectModule(self.module)
			}
		}
		return nil, nil
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
			case CREATE_CHILD:
				self.module = &schema.Module{}
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) SelectModule(module *schema.Module) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (child Selection, err error) {
		switch s.State.Position.GetIdent() {
		case "revision":
			s.State.Found = module.Revision != nil
			if s.State.Found {
				return self.selectRevision(module.Revision)
			}
		case "rpcs":
			s.State.Found = module.GetRpcs().GetFirstMeta() != nil
			return self.selectRpcs(module.GetRpcs())
		case "notifications":
			s.State.Found = module.GetNotifications().GetFirstMeta() != nil
			return self.selectNotifications(module.GetNotifications())
		default:
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, module)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), module, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
			case CREATE_CHILD:
				switch s.State.Position.GetIdent() {
					case "revision":
						module.Revision = &schema.Revision{}
					default:
						return EditNotImplemented(s.State.Position)
				}
			case UPDATE_VALUE:
				return WriteField(s.State.Position, module, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectRevision(rev *schema.Revision) (Selection, error) {
	s := &MySelection{}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		switch s.State.Position.GetIdent() {
		case "rev-date":
			return ReadFieldWithFieldName("Ident", s.State.Position.(schema.HasDataType), rev, val)
		default:
			return ReadField(s.State.Position.(schema.HasDataType), rev, val)
		}
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
			case UPDATE_VALUE:
				switch s.State.Position.GetIdent() {
					case "rev-date":
						rev.Ident = val.Str
					default:
						return WriteField(s.State.Position, rev, val)
				}
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectType(typeData *schema.DataType) (Selection, error) {
	s := &MySelection{}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), typeData, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, typeData, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectGroupings(groupings schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:groupings, resolve:self.resolve}
	var created *schema.Grouping
	s.OnSelect = func() (Selection, error) {
		s.State.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, i.data)
	}
	s.OnNext = i.Iterate
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			created = &schema.Grouping{}
			i.data = created
		case POST_CREATE_LIST_ITEM:
			groupings.AddMeta(created)
		case UPDATE_VALUE:
			return WriteField(s.State.Position, i.data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectRpcIO(i *schema.RpcInput, o *schema.RpcOutput) (Selection, error) {
	s := &MySelection{}
	var io schema.MetaList
	if i != nil {
		io = i
	} else {
		io = o
	}
	s.OnSelect = func() (Selection, error) {
		state := s.WalkState()
		state.Found = true
		if state.Found {
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, io)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), io, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, io, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) CreateGroupingsTypedefsDefinitions(parent schema.MetaList, childMeta schema.Meta) (schema.Meta, error) {
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
		default:
			return nil, NotImplemented(childMeta)
	}
	parent.AddMeta(child)
	return child, nil
}

func (self *SchemaBrowser) selectRpcs(rpcs schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:rpcs, resolve:self.resolve}
	var created *schema.Rpc
	s.OnSelect = func() (Selection, error) {
		rpc := i.data.(*schema.Rpc)
		switch s.State.Position.GetIdent() {
		case "input":
			s.State.Found = rpc.Input != nil
			if s.State.Found {
				return self.selectRpcIO(rpc.Input, nil)
			}
		case "output":
			s.State.Found = rpc.Output != nil
			if s.State.Found {
				return self.selectRpcIO(nil, rpc.Output)
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
			case CREATE_LIST_ITEM:
				created = &schema.Rpc{}
				i.data = created
			case POST_CREATE_LIST_ITEM:
				rpcs.AddMeta(created)
			case CREATE_CHILD:
				rpc := i.data.(*schema.Rpc)
				switch s.State.Position.GetIdent() {
				case "input":
					rpc.AddMeta(&schema.RpcInput{})
				case "output":
					rpc.AddMeta(&schema.RpcOutput{})
				}
			case UPDATE_VALUE:
				return WriteField(s.State.Position, i.data, val)
		}
		return nil
	}
	s.OnNext = i.Iterate
	return s, nil
}

func (self *SchemaBrowser) selectTypedefs(typedefs schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:typedefs, resolve:self.resolve}
	var created *schema.Typedef
	s.OnSelect = func() (Selection, error) {
		tdef := i.data.(*schema.Typedef)
		switch s.State.Position.GetIdent() {
		case "type":
			s.State.Found = tdef.DataType != nil
			if s.State.Found {
				return self.selectType(tdef.DataType)
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnNext = i.Iterate
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case CREATE_CHILD:
			tdef := i.data.(*schema.Typedef)
			switch s.State.Position.GetIdent() {
			case "type":
				tdef.SetDataType(&schema.DataType{})
			default:
				return NotImplemented(s.State.Position)
			}
		case CREATE_LIST_ITEM:
			created = &schema.Typedef{}
			i.data = created
		case POST_CREATE_LIST_ITEM:
			typedefs.AddMeta(created)
		case UPDATE_VALUE:
			return WriteField(s.State.Position, i.data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) GroupingsTypedefsDefinitions(s Selection, meta schema.Meta, data schema.Meta) (Selection, error) {
	state := s.WalkState()
	switch meta.GetIdent() {
	case "groupings":
		groupings := data.(schema.HasGroupings).GetGroupings()
		state.Found = !self.resolve && schema.ListLen(groupings) > 0
		return self.selectGroupings(groupings)
	case "typedefs":
		typedefs := data.(schema.HasTypedefs).GetTypedefs()
		state.Found = !self.resolve && schema.ListLen(typedefs) > 0
		return self.selectTypedefs(typedefs)
	case "definitions":
		defs := data.(schema.MetaList)
		state.Found = schema.ListLen(defs) > 0
		return self.selectDefinitionsList(defs)
	}
	return nil, nil
}

func (self *SchemaBrowser) selectNotifications(notifications schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:notifications, resolve:self.resolve}
	var created *schema.Notification
	s.OnSelect = func() (Selection, error) {
		s.State.Found = i.data != nil
		if s.State.Found {
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, i.data)
		}
		return nil, nil
	}
	s.OnNext = i.Iterate
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case CREATE_LIST_ITEM:
			created = &schema.Notification{}
			i.data = created
		case POST_CREATE_LIST_ITEM:
			notifications.AddMeta(created)
		case UPDATE_VALUE:
			return WriteField(s.State.Position, i.data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaList(data *schema.List) (Selection, error) {
	s := &MySelection{}
	s.WalkState().Found = true
	s.OnSelect = func() (Selection, error) {
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, data)
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaContainer(data schema.MetaList) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, data)
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaLeafy(leaf *schema.Leaf, leafList *schema.LeafList) (Selection, error) {
	s := &MySelection{}
	var leafy schema.HasDataType
	if leaf != nil {
		leafy = leaf
	} else {
		leafy = leafList
	}
	s.OnSelect = func() (Selection, error) {
		switch s.State.Position.GetIdent() {
		case "type":
			s.State.Found = leafy.GetDataType() != nil
			if s.State.Found {
				return self.selectType(leafy.GetDataType())
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), leafy, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case CREATE_CHILD:
			switch s.State.Position.GetIdent() {
			case "type":
				leafy.SetDataType(&schema.DataType{})
			default:
				return NotImplemented(s.State.Position)
			}
		case UPDATE_VALUE:
			return WriteField(s.State.Position, leafy, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaUses(data *schema.Uses) (Selection, error) {
	s := &MySelection{}
	// TODO: uses has refine container(s)
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, data, val)
		}
		return nil
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaCases(choice *schema.Choice) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:choice, resolve:self.resolve}
	s.OnNext = i.Iterate
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "definitions":
			return self.selectDefinitionsList(choice)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), choice, val)
	}
	return s, nil
}

func (self *SchemaBrowser) selectMetaChoice(data *schema.Choice) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "cases":
			return self.selectMetaCases(data);
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case UPDATE_VALUE:
			return WriteField(s.State.Position, data, val)
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
}

func (i *listIterator) Iterate(keys []Value, first bool) (bool, error) {
	i.data = nil
	if i.dataList == nil {
		return false, nil
	}
	if len(keys) > 0 {
		if first {
			i.data = schema.FindByIdent2(i.dataList, keys[0].Str)
		}
	} else {
		if first {
			i.iterator = schema.NewMetaListIterator(i.dataList, i.resolve)
		}
		if i.iterator.HasNextMeta() {
			i.data = i.iterator.NextMeta()
		}
	}
	return i.data != nil, nil
}

func (self *SchemaBrowser) selectDefinitionsList(dataList schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:dataList, resolve:self.resolve}
	s.OnChoose = func(choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, i.data)
	}
	s.OnSelect = func() (Selection, error) {
		s.State.Found = i.data != nil
		if s.State.Found {
			switch s.State.Position.GetIdent() {
			case "list":
				return self.selectMetaList(i.data.(*schema.List))
			case "leaf":
				return self.selectMetaLeafy(i.data.(*schema.Leaf), nil)
			case "leaf-list":
				return self.selectMetaLeafy(nil, i.data.(*schema.LeafList))
			case "uses":
				return self.selectMetaUses(i.data.(*schema.Uses))
			case "choice":
				return self.selectMetaChoice(i.data.(*schema.Choice))
			default:
				return self.selectMetaContainer(i.data.(schema.MetaList))
			}
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) error {
		s.State.Found = true
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnWrite = func(op Operation, val *Value) error {
		switch op {
		case CREATE_CHILD:
			var err error
			if i.data, err = self.CreateGroupingsTypedefsDefinitions(dataList, s.State.Position); err != nil {
				return err
			}
		}
		return nil
		return nil
	}
	s.OnNext = i.Iterate
	return s, nil

}

func (self *SchemaBrowser) resolveDefinitionCase(choice *schema.Choice, data schema.Meta) (caseMeta schema.MetaList, err error) {
	caseType := self.definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*schema.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg:msg}
	} else {
		return caseMeta, nil
	}
}

func (self *SchemaBrowser) definitionType(data schema.Meta) string {
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

const YANG_1_0 = `module yang {
    namespace "http://schema.org/yang";
    prefix "schema";
    description "Yang definition of yang";
    revision 2015-07-11 {
        description "Yang 1.0";
    }

    grouping def-header {
        leaf ident {
            type string;
        }
        leaf description {
            type string;
        }
    }

    grouping type {
        container type {
            leaf ident {
                type string;
            }
            leaf range {
                type string;
            }
            leaf-list enumeration {
                type string;
            }
        }
    }

    grouping groupings-typedefs {
        list groupings {
            key "ident";
            uses def-header;

            /*
              !! CIRCULAR
            */
            uses groupings-typedefs;
            uses containers-lists-leafs-uses-choice;
        }
        list typedefs {
            key "ident";
            uses def-header;
            uses type;
        }
    }

    grouping containers-lists-leafs-uses-choice {
        list definitions {
            key "ident";
            leaf ident {
            	type string;
            }
            choice body-stmt {
                case container {
                    container container {
                        uses def-header;
                        uses groupings-typedefs;
                        uses containers-lists-leafs-uses-choice;
                    }
                }
                case list {
                    container list {
                        uses def-header;
                        leaf-list keys {
                            type string;
                        }
                        uses groupings-typedefs;
                        uses containers-lists-leafs-uses-choice;
                    }
                }
                case leaf {
                    container leaf {
                        uses def-header;
                        leaf config {
                            type boolean;
                        }
                        leaf mandatory {
                            type boolean;
                        }
                        uses type;
                    }
                }
                case leaf-list {
                    container leaf-list {
                        uses def-header;
                        leaf config {
                            type boolean;
                        }
                        leaf mandatory {
                            type boolean;
                        }
                        uses type;
                    }
                }
                case uses {
                    container uses {
                        uses def-header;
                        /* need to expand this to use refine */
                    }
                }
                case choice {
                    container choice {
                        uses def-header;
                        list cases {
                            key "ident";
                            leaf ident {
                                type string;
                            }
                            /*
                             !! CIRCULAR
                            */
                            uses containers-lists-leafs-uses-choice;
                        }
                    }
                }
            }
        }
    }

    container module {
        uses def-header;
        leaf namespace {
            type string;
        }
        leaf prefix {
            type string;
        }
        container revision {
            leaf rev-date {
                type string;
            }
            leaf description {
                type string;
            }
        }
        list rpcs {
            key "ident";
            uses def-header;
            container input {
                uses groupings-typedefs;
                uses containers-lists-leafs-uses-choice;
            }
            container output {
                uses groupings-typedefs;
                uses containers-lists-leafs-uses-choice;
            }
        }
        list notifications {
            key "ident";
            uses def-header;
            uses groupings-typedefs;
            uses containers-lists-leafs-uses-choice;
        }
        uses groupings-typedefs;
        uses containers-lists-leafs-uses-choice;
    }
}`