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
type YangBrowser struct {
	module *schema.Module // read: meta
	meta *schema.Module // read: meta

	// resolve all uses, groups and typedefs.  if this is false, then depth must be
	// used to avoid infinite recursion
	resolve bool
}

func (self *YangBrowser) Module() *schema.Module {
	return self.meta
}

func NewYangBrowser(module *schema.Module, resolve bool) *YangBrowser {
	browser := &YangBrowser{module:module, meta:GetYangModule(), resolve:resolve}
	return browser
}

var yang1_0 *schema.Module
func GetYangModule() *schema.Module {
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

func (self *YangBrowser) RootSelector() (Selection, error) {
	s := &MySelection{}
	s.State.Meta = self.meta
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "module" :
			return self.SelectModule(self.module)
		}
		return nil, nil
	}
	return s, nil
}

func (self *YangBrowser) SelectModule(module *schema.Module) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (child Selection, err error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "revision":
			return self.selectRevision(module.Revision)
		case "rpcs":
			return self.selectRpcs(module.GetRpcs())
		case "notifications":
			s.WalkState().Found = schema.ListLen(module.GetNotifications()) > 0
			return self.selectNotifications(module.GetNotifications())
		default:
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, module)
		}
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), module, val)
	}
	return s, nil
}

func (self *YangBrowser) selectRevision(rev *schema.Revision) (Selection, error) {
	s := &MySelection{}
	s.OnRead = func(val *Value) (err error) {
		switch s.State.Position.GetIdent() {
		case "rev-date":
			return ReadFieldWithFieldName("Ident", s.State.Position.(schema.HasDataType), rev, val)
		default:
			return ReadField(s.State.Position.(schema.HasDataType), rev, val)
		}
	}
	return s, nil
}

func (self *YangBrowser) selectType(typeData *schema.DataType) (Selection, error) {
	s := &MySelection{}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), typeData, val)
	}
	return s, nil
}

func (self *YangBrowser) selectGroupings(groupings schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:groupings, resolve:self.resolve}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, i.data)
	}
	s.OnNext = i.Iterate
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectRpcInput(rpc *schema.RpcInput) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		state := s.WalkState()
		state.Found = rpc != nil
		if state.Found {
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, rpc)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), rpc, val)
	}
	return s, nil
}

func (self *YangBrowser) selectRpcOutput(rpc *schema.RpcOutput) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		state := s.WalkState()
		state.Found = rpc != nil
		if state.Found {
			return self.GroupingsTypedefsDefinitions(s, s.State.Position, rpc)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), rpc, val)
	}
	return s, nil
}

func (self *YangBrowser) selectRpcs(rpcs schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:rpcs, resolve:self.resolve}
	s.OnSelect = func() (Selection, error) {
		state := s.WalkState()
		state.Found = true
		switch s.State.Position.GetIdent() {
		case "input":
			return self.selectRpcInput(i.data.(*schema.Rpc).Input)
		case "output":
			return self.selectRpcOutput(i.data.(*schema.Rpc).Output)
		}

		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnNext = i.Iterate
	return s, nil
}

func (self *YangBrowser) selectTypedefs(typedefs schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:typedefs, resolve:self.resolve}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "type":
			return self.selectType(i.data.(*schema.Typedef).GetDataType())
		}

		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	s.OnNext = i.Iterate
	return s, nil
}

func (self *YangBrowser) GroupingsTypedefsDefinitions(s Selection, meta schema.Meta, data schema.Meta) (Selection, error) {
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

func (self *YangBrowser) selectNotifications(notifications schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:notifications, resolve:self.resolve}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, i.data)
	}
	s.OnNext = i.Iterate
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), i.data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaList(data *schema.List) (Selection, error) {
	s := &MySelection{}
	s.WalkState().Found = true
	s.OnSelect = func() (Selection, error) {
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, data)
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaContainer(data schema.MetaList) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		return self.GroupingsTypedefsDefinitions(s, s.State.Position, data)
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaLeaf(data *schema.Leaf) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "type":
			return self.selectType(data.DataType)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaLeafList(data *schema.LeafList) (Selection, error) {
	s := &MySelection{}
	s.OnSelect = func() (Selection, error) {
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "type":
			return self.selectType(data.DataType)
		}
		return nil, nil
	}
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaUses(data *schema.Uses) (Selection, error) {
	s := &MySelection{}
	// TODO: uses has refine container(s)
	s.OnRead = func(val *Value) (err error) {
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaCases(choice *schema.Choice) (Selection, error) {
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
		return ReadField(s.State.Position.(schema.HasDataType), choice, val)
	}
	return s, nil
}

func (self *YangBrowser) selectMetaChoice(data *schema.Choice) (Selection, error) {
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
		return ReadField(s.State.Position.(schema.HasDataType), data, val)
	}
	return s, nil
}

type listIterator struct {
	data schema.Meta
	dataList schema.MetaList
	iterator schema.MetaIterator
	resolve bool
}

func (i *listIterator) Iterate(keys []string, first bool) (bool, error) {
	i.data = nil
	if i.dataList == nil {
		return false, nil
	}
	if len(keys) > 0 {
		if first {
			i.data = schema.FindByIdent2(i.dataList, keys[0])
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

func (self *YangBrowser) selectDefinitionsList(dataList schema.MetaList) (Selection, error) {
	s := &MySelection{}
	i := listIterator{dataList:dataList, resolve:self.resolve}
	s.OnChoose = func(choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, i.data)
	}
	s.OnSelect = func() (Selection, error) {
		var e error
		choice := s.State.Meta.GetFirstMeta().(*schema.Choice)
		if s.State.Position, e = self.resolveDefinitionCase(choice, i.data); e != nil {
			return nil, e
		}
		s.WalkState().Found = true
		switch s.State.Position.GetIdent() {
		case "list":
			return self.selectMetaList(i.data.(*schema.List))
		case "leaf":
			return self.selectMetaLeaf(i.data.(*schema.Leaf))
		case "leaf-list":
			return self.selectMetaLeafList(i.data.(*schema.LeafList))
		case "uses":
			return self.selectMetaUses(i.data.(*schema.Uses))
		case "choice":
			return self.selectMetaChoice(i.data.(*schema.Choice))
		default:
			return self.selectMetaContainer(i.data.(schema.MetaList))
		}
		return nil, nil
	}
	s.OnNext = i.Iterate
	return s, nil

}

func (self *YangBrowser) resolveDefinitionCase(choice *schema.Choice, data schema.Meta) (caseMeta schema.MetaList, err error) {
	caseType := self.definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*schema.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg:msg}
	} else {
		return caseMeta, nil
	}
}

func (self *YangBrowser) definitionType(data schema.Meta) string {
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