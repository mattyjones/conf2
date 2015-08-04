package comm
import (
	"yang"
	"yang/browse"
	"fmt"
	"reflect"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type MetaTransmitter struct {
	module *yang.Module // read: meta
	meta *yang.Module
}

func (self *MetaTransmitter) RootSelector() (s *browse.Selection, err error) {
	s = &browse.Selection{Meta:self.meta.DataDefs().GetFirstMeta().(yang.MetaList)}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "revision":
			return selectRevision(self.module.Revision)
		case "rpcs":
			return selectRpcs(self.module.GetRpcs())
		case "notifications":
			return selectNotifications(self.module.GetNotifications())
		default:
			return GroupingsTypedefsDefinitions(s.Position, self.module)
		}
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, self.module, val)
	}
	return
}

func selectRevision(rev *yang.Revision) (s *browse.Selection, serr error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		switch s.Position.GetIdent() {
		case "rev-date":
			return browse.ReadFieldWithFieldName("Ident", s.Position, rev, val)
		default:
			return browse.ReadField(s.Position, rev, val)
		}
	}
	return
}

func selectType(typeData *yang.DataType) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return  nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, typeData, val)
	}
	return
}

func selectGroupings(groupings yang.MetaList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	i := &defsIterator{dataList:groupings}
	s.ListIterator = i.ListIterator
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "input":
			return selectDefinitionsList(i.data.(*yang.Rpc).Input)
		case "output":
			return selectDefinitionsList(i.data.(*yang.Rpc).Output)
		}

		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, i.data, val)
	}
	return
}

func selectRpcs(rpcs yang.MetaList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	i := &defsIterator{dataList:rpcs}
	s.ListIterator = i.ListIterator
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "input":
			return selectDefinitionsList(i.data.(*yang.Rpc).Input)
		case "output":
			return selectDefinitionsList(i.data.(*yang.Rpc).Output)
		}

		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, i.data, val)
	}
	return
}

func selectTypedefs(typedefs yang.MetaList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	i := &defsIterator{dataList:typedefs}
	s.ListIterator = i.ListIterator
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "type":
			return selectType(i.data.(*yang.Typedef).GetDataType())
		}

		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, i.data, val)
	}
	return
}

func GroupingsTypedefsDefinitions(meta yang.Meta, data yang.Meta) (*browse.Selection, error) {
	switch meta.GetIdent() {
	case "groupings":
		return selectGroupings(data.(yang.HasGroupings).GetGroupings())
	case "typedefs":
		return selectTypedefs(data.(yang.HasTypedefs).GetTypedefs())
	case "definitions":
		return selectDefinitionsList(data.(yang.MetaList))
	}
	return nil, nil
}

func selectNotifications(notifications yang.MetaList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	i := &defsIterator{dataList:notifications}
	s.ListIterator = i.ListIterator
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return GroupingsTypedefsDefinitions(s.Position, i.data)
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, i.data, val)
	}
	return
}

func selectMetaList(data *yang.List) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return GroupingsTypedefsDefinitions(s.Position, data)
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaContainer(data *yang.Container) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return GroupingsTypedefsDefinitions(s.Position, data)
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaLeaf(data *yang.Leaf) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "type":
			return selectType(data.DataType)
		}
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaLeafList(data *yang.LeafList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "type":
			return selectType(data.DataType)
		}
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaUses(data *yang.Uses) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaCases(data *yang.Choice) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "definitions":
			return selectDefinitionsList(data)
		}
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}

func selectMetaChoice(data *yang.Choice) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	s.Selector = func(ident string) (*browse.Selection, error) {
		s.Position = yang.FindByIdent2(s.Meta, ident)
		switch ident {
		case "cases":
			return selectMetaCases(data);
		}
		return nil, nil
	}
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, data, val)
	}
	return
}


/**
 when the data is a MetaList, this handles the browse ListIterator-ing
 */
type defsIterator struct {
	data yang.Meta
	iterator yang.MetaIterator
	dataList yang.MetaList
}

func (dl *defsIterator) ListIterator(keys []string, first bool) (bool, error)  {
	dl.data = nil
	if len(keys) > 0 {
		if first {
			dl.data = yang.FindByIdent2(dl.dataList, keys[0])
		}
	} else {
		if first {
			dl.iterator = yang.NewMetaListIterator(dl.dataList, false)
		}
		if dl.iterator.HasNextMeta() {
			dl.data = dl.iterator.NextMeta()
		}
	}
	return dl.data != nil, nil
}

func selectDefinitionsList(dataList yang.MetaList) (s *browse.Selection, err error) {
	s = &browse.Selection{}
	i := &defsIterator{dataList:dataList}
	s.ListIterator = i.ListIterator
	s.Selector = func(ident string) (s2 *browse.Selection, e error) {
		choice := s.Meta.GetFirstMeta().(*yang.Choice)
		if s.Position, e = resolveDefinitionCase(choice, i.data); e != nil {
			return nil, e
		}
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
	s.Reader = func(val *browse.Value) (err error) {
		return browse.ReadField(s.Position, i.data, val)
	}
	return

}

func resolveDefinitionCase(choice *yang.Choice, data yang.Meta) (caseMeta yang.MetaList, err error) {
	caseType := definitionType(data)
	if caseMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*yang.Container); !ok {
		msg := fmt.Sprint("Could not find case meta for ", caseType)
		return nil, &commError{msg}
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
