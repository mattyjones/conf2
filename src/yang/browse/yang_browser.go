package browse
import (
	"yang"
	"fmt"
	"reflect"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type YangBrowser struct {
	Module *yang.Module // read: meta
	Meta *yang.Module // read: meta
}

var yang1_0 *yang.Module
func NewYangBrowser(module *yang.Module) *YangBrowser {
	if yang1_0 == nil {
		var err error
		yang1_0, err = yang.LoadModuleFromByteArray([]byte(YANG_1_0))
		if err != nil {
			msg := fmt.Sprintf("Error parsing yang-1.0 yang, %s", err.Error())
			panic(msg)
		}
	}
	browser := &YangBrowser{Module:module, Meta:yang1_0}
	return browser
}

type MetaListSelector func(m yang.Meta) (*Selection, error)

func (self *YangBrowser) RootSelector() (s *Selection, err error) {
	s = &Selection{Meta:self.Meta}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "module" :
			return selectModule(self.Module)
		}
		return nil, nil
	}
	return
}

func selectModule(module *yang.Module) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (child *Selection, err error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "revision":
			return selectRevision(module.Revision)
		case "rpcs":
			return selectRpcs(module.GetRpcs())
		case "notifications":
			s.Found = yang.ListLen(module.GetNotifications()) > 0
			return selectNotifications(module.GetNotifications())
		default:
			return GroupingsTypedefsDefinitions(s, s.Position, module)
		}
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), module, val)
	}
	return
}

func selectRevision(rev *yang.Revision) (*Selection, error) {
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

func selectType(typeData *yang.DataType) (s *Selection, err error) {
	s = &Selection{}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), typeData, val)
	}
	return
}

func selectGroupings(groupings yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:groupings}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	return
}

func selectRpcInput(rpc *yang.RpcInput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), rpc, val)
	}
	return
}

func selectRpcOutput(rpc *yang.RpcOutput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), rpc, val)
	}
	return
}

func selectRpcs(rpcs yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:rpcs}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "input":
			return selectRpcInput(i.data.(*yang.Rpc).Input)
		case "output":
			return selectRpcOutput(i.data.(*yang.Rpc).Output)
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func selectTypedefs(typedefs yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:typedefs}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return selectType(i.data.(*yang.Typedef).GetDataType())
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func GroupingsTypedefsDefinitions(s *Selection, meta yang.Meta, data yang.Meta) (*Selection, error) {
	switch meta.GetIdent() {
	case "groupings":
		groupings := data.(yang.HasGroupings).GetGroupings()
		s.Found = yang.ListLen(groupings) > 0
		return selectGroupings(groupings)
	case "typedefs":
		typedefs := data.(yang.HasTypedefs).GetTypedefs()
		s.Found = yang.ListLen(typedefs) > 0
		return selectTypedefs(typedefs)
	case "definitions":
		defs := data.(yang.MetaList)
		s.Found = yang.ListLen(defs) > 0
		return selectDefinitionsList(defs)
	}
	return nil, nil
}

func selectNotifications(notifications yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:notifications}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), i.data, val)
	}
	return
}

func selectMetaList(data *yang.List) (s *Selection, err error) {
	s = &Selection{}
	s.Found = true
	s.Enter = func() (*Selection, error) {
		return GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaContainer(data *yang.Container) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaLeaf(data *yang.Leaf) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return selectType(data.DataType)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaLeafList(data *yang.LeafList) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return selectType(data.DataType)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaUses(data *yang.Uses) (s *Selection, err error) {
	s = &Selection{}
	// TODO: uses has refine container(s)
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaCases(data *yang.Choice) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "definitions":
			return selectDefinitionsList(data)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(yang.HasDataType), data, val)
	}
	return
}

func selectMetaChoice(data *yang.Choice) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "cases":
			return selectMetaCases(data);
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
			i.iterator = yang.NewMetaListIterator(i.dataList, false)
		}
		if i.iterator.HasNextMeta() {
			i.data = i.iterator.NextMeta()
		}
	}
	return i.data != nil, nil
}

func selectDefinitionsList(dataList yang.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:dataList}
	s.Choose = func(choice *yang.Choice) (m yang.Meta, err error) {
		return resolveDefinitionCase(choice, i.data)
	}
	s.Enter = func() (s2 *Selection, e error) {
		choice := s.Meta.GetFirstMeta().(*yang.Choice)
		if s.Position, e = resolveDefinitionCase(choice, i.data); e != nil {
			return nil, e
		}
		s.Found = true
		switch s.Position.GetIdent() {
		case "list":
			return selectMetaList(i.data.(*yang.List))
		case "container":
			return selectMetaContainer(i.data.(*yang.Container))
		case "leaf":
			return selectMetaLeaf(i.data.(*yang.Leaf))
		case "leaf-list":
			return selectMetaLeafList(i.data.(*yang.LeafList))
		case "uses":
			return selectMetaUses(i.data.(*yang.Uses))
		case "choice":
			return selectMetaChoice(i.data.(*yang.Choice))
		}
		return nil, nil
	}
	s.Iterate = i.Iterate
	return

}

func resolveDefinitionCase(choice *yang.Choice, data yang.Meta) (caseMeta yang.MetaList, err error) {
	caseType := definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*yang.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &browseError{Msg:msg}
	} else {
		return caseMeta, nil
	}
}

func definitionType(data yang.Meta) string {
	switch data.(type) {
	case *yang.List:
		return "list"
	case *yang.Container:
		return "container"
	case *yang.Uses:
		return "uses"
	case *yang.Choice:
		return "choice"
	case *yang.Leaf:
		return "leaf"
	case *yang.LeafList:
		return "leaf-list"
	default:
		msg := fmt.Sprint("unknown type", reflect.TypeOf(data))
		panic(msg)
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
                            type string;
                        }
                        leaf mandatory {
                            type string;
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