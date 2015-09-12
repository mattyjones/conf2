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

type MetaListSelector func(m schema.Meta) (*Selection, error)

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

func (self *YangBrowser) SelectModule(module *schema.Module) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (child *Selection, err error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "revision":
			return self.selectRevision(module.Revision)
		case "rpcs":
			return self.selectRpcs(module.GetRpcs())
		case "notifications":
			s.Found = schema.ListLen(module.GetNotifications()) > 0
			return self.selectNotifications(module.GetNotifications())
		default:
			return self.GroupingsTypedefsDefinitions(s, s.Position, module)
		}
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), module, val)
	}
	return
}

func (self *YangBrowser) selectRevision(rev *schema.Revision) (*Selection, error) {
	s := &Selection{}
	s.ReadValue = func(val *Value) (err error) {
		switch s.Position.GetIdent() {
		case "rev-date":
			return ReadFieldWithFieldName("Ident", s.Position.(schema.HasDataType), rev, val)
		default:
			return ReadField(s.Position.(schema.HasDataType), rev, val)
		}
	}
	return s, nil
}

func (self *YangBrowser) selectType(typeData *schema.DataType) (s *Selection, err error) {
	s = &Selection{}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), typeData, val)
	}
	return
}

func (self *YangBrowser) selectGroupings(groupings schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:groupings, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), i.data, val)
	}
	return
}

func (self *YangBrowser) selectRpcInput(rpc *schema.RpcInput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return self.GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), rpc, val)
	}
	return
}

func (self *YangBrowser) selectRpcOutput(rpc *schema.RpcOutput) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = rpc != nil
		if s.Found {
			return self.GroupingsTypedefsDefinitions(s, s.Position, rpc)
		}
		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), rpc, val)
	}
	return
}

func (self *YangBrowser) selectRpcs(rpcs schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:rpcs, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "input":
			return self.selectRpcInput(i.data.(*schema.Rpc).Input)
		case "output":
			return self.selectRpcOutput(i.data.(*schema.Rpc).Output)
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func (self *YangBrowser) selectTypedefs(typedefs schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:typedefs, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		switch s.Position.GetIdent() {
		case "type":
			return self.selectType(i.data.(*schema.Typedef).GetDataType())
		}

		return nil, nil
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), i.data, val)
	}
	s.Iterate = i.Iterate
	return
}

func (self *YangBrowser) GroupingsTypedefsDefinitions(s *Selection, meta schema.Meta, data schema.Meta) (*Selection, error) {
	switch meta.GetIdent() {
	case "groupings":
		groupings := data.(schema.HasGroupings).GetGroupings()
		s.Found = !self.resolve && schema.ListLen(groupings) > 0
		return self.selectGroupings(groupings)
	case "typedefs":
		typedefs := data.(schema.HasTypedefs).GetTypedefs()
		s.Found = !self.resolve && schema.ListLen(typedefs) > 0
		return self.selectTypedefs(typedefs)
	case "definitions":
		defs := data.(schema.MetaList)
		s.Found = schema.ListLen(defs) > 0
		return self.selectDefinitionsList(defs)
	}
	return nil, nil
}

func (self *YangBrowser) selectNotifications(notifications schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:notifications, resolve:self.resolve}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, i.data)
	}
	s.Iterate = i.Iterate
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), i.data, val)
	}
	return
}

func (self *YangBrowser) selectMetaList(data *schema.List) (s *Selection, err error) {
	s = &Selection{}
	s.Found = true
	s.Enter = func() (*Selection, error) {
		return self.GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaContainer(data schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	s.Enter = func() (*Selection, error) {
		s.Found = true
		return self.GroupingsTypedefsDefinitions(s, s.Position, data)
	}
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaLeaf(data *schema.Leaf) (s *Selection, err error) {
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
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaLeafList(data *schema.LeafList) (s *Selection, err error) {
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
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaUses(data *schema.Uses) (s *Selection, err error) {
	s = &Selection{}
	// TODO: uses has refine container(s)
	s.ReadValue = func(val *Value) (err error) {
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
}

func (self *YangBrowser) selectMetaCases(choice *schema.Choice) (s *Selection, err error) {
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
		return ReadField(s.Position.(schema.HasDataType), choice, val)
	}
	return
}

func (self *YangBrowser) selectMetaChoice(data *schema.Choice) (s *Selection, err error) {
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
		return ReadField(s.Position.(schema.HasDataType), data, val)
	}
	return
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

func (self *YangBrowser) selectDefinitionsList(dataList schema.MetaList) (s *Selection, err error) {
	s = &Selection{}
	i := listIterator{dataList:dataList, resolve:self.resolve}
	s.Choose = func(choice *schema.Choice) (m schema.Meta, err error) {
		return self.resolveDefinitionCase(choice, i.data)
	}
	s.Enter = func() (s2 *Selection, e error) {
		choice := s.Meta.GetFirstMeta().(*schema.Choice)
		if s.Position, e = self.resolveDefinitionCase(choice, i.data); e != nil {
			return nil, e
		}
		s.Found = true
		switch s.Position.GetIdent() {
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
	s.Iterate = i.Iterate
	return

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