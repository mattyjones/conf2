package browse
import (
	"yang"
	"fmt"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type YangBrowser struct {
	module *yang.Module // read: meta
	meta *yang.Module // read: meta

	// resolve all uses, groups and typedefs.  if this is false, then depth must be
	// used to avoid infinite recursion
	resolve bool
}

func (self *YangBrowser) Module() *yang.Module {
	return self.meta
}

func NewYangBrowser(module *yang.Module, resolve bool) *YangBrowser {
	browser := &YangBrowser{module:module, meta:GetYangModule(), resolve:resolve}
	return browser
}

var yang1_0 *yang.Module
func GetYangModule() *yang.Module {
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

type MetaListSelector func(m yang.Meta) (*Selection, error)

func (self *YangBrowser) RootSelector() (s *Selection, err error) {
	s = &Selection{Meta:self.meta}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "module" :
			return self.SelectModule(self.module)
		}
		return nil, nil
	}
	return
}

func (self *YangBrowser) SelectModule(module *yang.Module) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (child *Selection, err error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "revision":
			return self.selectRevision(module.Revision)
		case "rpcs":
			return self.selectRpcs(module.GetRpcs())
		case "notifications":
			s.Found = yang.ListLen(module.GetNotifications()) > 0
			return self.selectNotifications(module.GetNotifications())
		default:
			return self.GroupingsTypedefsDefinitions(s, s.Position, module)
		}
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), module, val)
	}
	return
}

func (self *YangBrowser) selectRevision(rev *yang.Revision) (*Selection, error) {
	s := &Selection{}
	s.ReadValue = func(val *Value) (err error) {
		switch s.Position.GetIdent() {
		case "rev-date":
			return ReadFieldWithFieldName("Ident", s.Position.(yang.HasDataType), rev, val)
		default:
			return ReadField(s.Position.(yang.HasDataType), rev, val)
		}
	}
	return s, nil
}

func (self *YangBrowser) selectType(typeData *yang.DataType) (s *Selection, err error) {
	s = &Selection{}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), typeData, val)
	}
	return
}

func (self *YangBrowser) selectGroupings(groupings yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:groupings, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	return
}

func (self *YangBrowser) selectRpcInput(rpc *yang.RpcInput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return self.GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), rpc, val)
	}
	return
}

func (self *YangBrowser) selectRpcOutput(rpc *yang.RpcOutput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return self.GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), rpc, val)
	}
	return
}

func (self *YangBrowser) selectRpcs(rpcs yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:rpcs, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "input":
			return self.selectRpcInput(i.data.(*yang.Rpc).Input)
		case "output":
			return self.selectRpcOutput(i.data.(*yang.Rpc).Output)
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func (self *YangBrowser) selectTypedefs(typedefs yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:typedefs, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return self.selectType(i.data.(*yang.Typedef).GetDataType())
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func (self *YangBrowser) GroupingsTypedefsDefinitions(s *Selection, meta yang.Meta, data yang.Meta) (*Selection, error) {
	switch meta.GetIdent() {
	case "groupings":
		groupings := data.(yang.HasGroupings).GetGroupings()
		s.Found = !self.resolve && yang.ListLen(groupings) > 0
		return self.selectGroupings(groupings)
	case "typedefs":
		typedefs := data.(yang.HasTypedefs).GetTypedefs()
		s.Found = !self.resolve && yang.ListLen(typedefs) > 0
		return self.selectTypedefs(typedefs)
	case "definitions":
		defs := data.(yang.MetaList)
		s.Found = yang.ListLen(defs) > 0
		return self.selectDefinitionsList(defs)
	}
	return nil, nil
}

func (self *YangBrowser) selectNotifications(notifications yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:notifications, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	return
}

func (self *YangBrowser) selectMetaList(data *yang.List) (s *Selection, err error) {
	s = &Selection{}
	s.Found = true
	s.Enter = func() (*Selection, error) {
		return self.GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaContainer(data yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaLeaf(data *yang.Leaf) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return self.selectType(data.DataType)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaLeafList(data *yang.LeafList) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return self.selectType(data.DataType)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaUses(data *yang.Uses) (s *Selection, err error) {
	s = &Selection{}
	// TODO: uses has refine container(s)
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaCases(choice *yang.Choice) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:choice, resolve:self.resolve}
	s.Iterate = i.Iterate
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "definitions":
			return self.selectDefinitionsList(choice)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), choice, val)
	}
	return
}

func (self *YangBrowser) selectMetaChoice(data *yang.Choice) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "cases":
			return self.selectMetaCases(data);
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

type listIterator struct {
	data yang.Meta
	dataList yang.MetaList
	iterator yang.MetaIterator
	resolve bool
}

func (i *listIterator) Iterate(keys []string, first bool) (bool, error) {
	i.data = nil
	if i.dataList == nil {
		return false, nil
	}
	if len(keys) > 0 {
		if first {
			i.data = yang.FindByIdent2(i.dataList, keys[0])
		}
	} else {
		if first {
			i.iterator = yang.NewMetaListIterator(i.dataList, i.resolve)
		}
		if i.iterator.HasNextMeta() {
			i.data = i.iterator.NextMeta()
		}
	}
	return i.data != nil, nil
}

func (self *YangBrowser) selectDefinitionsList(dataList yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:dataList, resolve:self.resolve}
	s.Choose = func(choice *yang.Choice) (m yang.Meta, err error) {
		return self.resolveDefinitionCase(choice, i.data)
	}
	s.Enter = func() (s2 *Selection, e error) {
		choice := s.Meta.GetFirstMeta().(*yang.Choice)
		if s.Position, e = self.resolveDefinitionCase(choice, i.data); e != nil {
			return nil, e
		}
		s.Found = true
		switch s.Position.GetIdent() {
		case "list":
			return self.selectMetaList(i.data.(*yang.List))
		case "leaf":
			return self.selectMetaLeaf(i.data.(*yang.Leaf))
		case "leaf-list":
			return self.selectMetaLeafList(i.data.(*yang.LeafList))
		case "uses":
			return self.selectMetaUses(i.data.(*yang.Uses))
		case "choice":
			return self.selectMetaChoice(i.data.(*yang.Choice))
		default:
			return self.selectMetaContainer(i.data.(yang.MetaList))
		}
		return nil, nil
	}
	s.Iterate = i.Iterate
	return

}

func (self *YangBrowser) resolveDefinitionCase(choice *yang.Choice, data yang.Meta) (caseMeta yang.MetaList, err error) {
	caseType := self.definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*yang.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg:msg}
	} else {
		return caseMeta, nil
	}
}

func (self *YangBrowser) definitionType(data yang.Meta) string {
	switch data.(type) {
	case *yang.List:
		return "list"
	case *yang.Uses:
		return "uses"
	case *yang.Choice:
		return "choice"
	case *yang.Leaf:
		return "leaf"
	case *yang.LeafList:
		return "leaf-list"
	default:
		return "container"
	}
}

const YANG_1_0 = `module yang {
    namespace "http://yang.org/yang";
    prefix "yang";
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