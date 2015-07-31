package comm
import (
	"yang"
	"yang/browse"
	"fmt"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type MetaTransmitter struct {
	data *yang.Module // read: meta
	meta *yang.Module
}

func (self *MetaTransmitter) GetSelector() (s browse.Selection) {
	moduleMeta := self.meta.GetFirstMeta().(yang.MetaList)
	return selectModule(moduleMeta, self.data)
}

func selectModule(containerMeta yang.MetaList, data *yang.Module) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "revision":
			selection = selectRevision(data.Revision)
		case "rpcs":
			selection = selectList(selectRpc, data.GetRpcs())
		case "notifications":
			/* TODO: replace w/notifications */
			selection = selectList(selectContainerContainer, data.GetNotifications())
		case "typedefs":
			selection = selectList(selectTypedef, data.GetTypedefs())
		case "groupings":
			selection = selectList(selectGrouping, data.GetTypedefs())
		case "definitions":
			selection = selectDefinitionsList(meta.(yang.MetaList), data.DataDefs())
		default:
			return defaultHandler(s, op, containerMeta, meta, data, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}

		return
	}
	return s
}

func selectType(typeMeta yang.Meta, dataType *yang.DataType) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		return defaultHandler(s, op, typeMeta, meta, dataType, v)
	}
	return s
}

func selectTypedef(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	tdef := data.(*yang.Typedef)
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "type":
			selection = selectType(meta, tdef.GetDataType())
		default:
			return defaultHandler(s, op, meta.GetParent(), meta, tdef, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectGrouping(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "typedefs":
			selection = selectList(selectTypedef, data.(yang.HasTypedefs).GetTypedefs())
		case "groupings":
			selection = selectList(selectGrouping, data.(yang.HasGroupings).GetGroupings())
		case "definitions":
			selection = selectDefinitionsList(meta.(yang.MetaList), data.(yang.MetaList))
		default:
			return defaultHandler(s, op, containerMeta, meta, data, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectRpc(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
			case "input":
				selection = selectDefinitionsList(meta.(yang.MetaList), data.(*yang.Rpc).Input)
			case "output":
				selection = selectDefinitionsList(meta.(yang.MetaList), data.(*yang.Rpc).Output)
			default:
				return defaultHandler(s, op, containerMeta, meta, data, v)
		}
		switch op {
		case browse.READ_VALUE:
			return selection(op, meta, v)
		case browse.SELECT_CHILD:
			v.Selection = selection
		}

		return
	}
	return s
}

func selectRevision(rev *yang.Revision) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		switch meta.GetIdent() {
		case "rev-date":
			if op == browse.READ_VALUE {
				v.Val.Str = rev.Ident
				err = v.Send(meta)
			} else {
				err = browse.NotImplemented(meta)
			}
		case "revision":
			if err = v.EnterContainer(meta); err != nil {
				i := yang.NewMetaListIterator(meta.(yang.MetaList), true)
				for (i.HasNextMeta()) {
					if err = s(browse.READ_VALUE, i.NextMeta(), v); err != nil {
						return err
					}
				}
				err = v.ExitContainer(meta)
			}
		default:
			err = browse.UseReflection(op, meta, rev, v)
		}
		return
	}
	return s
}

func selectListContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "typedefs":
			selection = selectList(selectTypedef, data.(yang.HasTypedefs).GetTypedefs())
		case "groupings":
			selection = selectList(selectGrouping, data.(yang.HasGroupings).GetGroupings())
		case "definitions":
			selection = selectDefinitionsList(meta.(yang.MetaList), data.(yang.MetaList))
		default:
			return defaultHandler(s, op, containerMeta, meta, data, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectUses(usesMeta yang.Meta, uses *yang.Uses) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		return defaultHandler(s, op, usesMeta, meta, uses, v)

	}
	return s
}

func selectContainerContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "typedefs":
			selection = selectList(selectTypedef, data.(yang.HasTypedefs).GetTypedefs())
		case "groupings":
			selection = selectList(selectGrouping, data.(yang.HasGroupings).GetGroupings())
		case "definitions":
			selection = selectDefinitionsList(meta.(yang.MetaList), data.(yang.MetaList))
		default:
			return defaultHandler(s, op, containerMeta, meta, data, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectLeafContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	leaf := data.(*yang.Leaf)
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "type":
			selection = selectType(meta, leaf.GetDataType())
		default:
			return defaultHandler(s, op, containerMeta, meta, leaf, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectLeafListContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	leafList := data.(*yang.LeafList)
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		var selection browse.Selection
		switch meta.GetIdent() {
		case "type":
			selection = selectType(meta, leafList.GetDataType())
		default:
			return defaultHandler(s, op, containerMeta, meta, leafList, v)
		}
		switch op {
		case browse.SELECT_CHILD:
			v.Selection = selection
		case browse.READ_VALUE:
			return selection(browse.READ_VALUE, meta, v)
		}
		return
	}
	return s
}

func selectUsesContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		return defaultHandler(s, op, containerMeta, meta, data, v)
	}
	return s
}

func selectChoiceContainer(containerMeta yang.MetaList, data yang.Meta) (s browse.Selection) {
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		return defaultHandler(s, op, containerMeta, meta, data, v)
	}
	return s
}

func definitionSelection(choice *yang.Choice, data yang.Meta) (selection browse.Selection, caseMeta yang.MetaList, err error) {
	if caseMeta, err = resolveDefinitionCase(choice, data); err == nil {
		switch caseMeta.GetIdent() {
		case "list":
			selection = selectListContainer(caseMeta, data)
		case "container":
			selection = selectContainerContainer(caseMeta, data)
		case "leaf":
			selection = selectLeafContainer(caseMeta, data)
		case "leaf-list":
			selection = selectLeafListContainer(caseMeta, data)
		case "uses":
			selection = selectUsesContainer(caseMeta, data)
		case "choice":
			selection = selectChoiceContainer(caseMeta, data)
		}
	}
	return
}

func selectDefinitionsList(containerMeta yang.MetaList, data yang.MetaList) (s browse.Selection) {
	choice := containerMeta.GetFirstMeta().(*yang.Choice)
	s = func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		switch meta.GetIdent() {
		case "body-stmt":

		}
		switch op {
		case browse.SELECT_CHILD:
			all := yang.NewMetaListIterator(data, false)
			found := yang.FindByIdent(all, v.Val.Keys[0])
			if found == nil {
				return v.NotFound(v.Val.Keys[0])
			}
			v.Selection, v.Position, err = definitionSelection(choice, found)
		case browse.READ_VALUE:
			v.EnterList(meta)
			all := yang.NewMetaListIterator(data, false)
			for all.HasNextMeta() {
				v.EnterListItem(meta)
				def := all.NextMeta()
				selection, caseMeta, err := definitionSelection(choice, def)
				if err != nil {
					return err
				}
				if err = selection(browse.READ_VALUE, caseMeta, v); err != nil {
					return err
				}
				v.ExitListItem(meta)
			}
			v.ExitList(meta)
		default:
			return browse.NotImplemented(meta)
		}

		return
	}

	return s
}

func selectList(delegate SelectionDelegate, data yang.MetaList) browse.Selection {
	return func(op browse.Operation, meta yang.Meta, v *browse.Visitor) (err error) {
		switch op {
		case browse.SELECT_CHILD:
			all := yang.NewMetaListIterator(data, false)
			found := yang.FindByIdent(all, v.Val.Keys[0])
			if found == nil {
				return v.NotFound(v.Val.Keys[0])
			}
			v.Selection = delegate(meta.(yang.MetaList), found)
		case browse.READ_VALUE:
			v.EnterList(meta)
			all := yang.NewMetaListIterator(data, false)
			for all.HasNextMeta() {
				v.EnterListItem(meta)
				def := all.NextMeta()
				selection := delegate(meta.(yang.MetaList), def)
				if err = selection(browse.READ_VALUE, meta, v); err != nil {
					return err
				}
				v.ExitListItem(meta)
			}
			v.ExitList(meta)
		default:
			return browse.NotImplemented(meta)
		}
		return
	}
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
		panic("unknown type")
	}
}

type SelectionDelegate func(yang.MetaList, yang.Meta) browse.Selection

func readAll(selection browse.Selection, meta yang.Meta, v *browse.Visitor) (err error) {
	if err = v.EnterContainer(meta); err != nil {
		i := yang.NewMetaListIterator(meta.(yang.MetaList), true)
		for (i.HasNextMeta()) {
			if err = selection(browse.READ_VALUE, i.NextMeta(), v); err != nil {
				return err
			}
		}
		err = v.ExitContainer(meta)
	}
	return
}

func defaultHandler(s browse.Selection, op browse.Operation, containerMeta yang.Meta, meta yang.Meta, data interface{}, v *browse.Visitor) (err error) {
	if op == browse.READ_VALUE && meta == containerMeta {
		return readAll(s, meta, v)
	} else {
		return browse.UseReflection(op, meta, data, v)
	}
}

