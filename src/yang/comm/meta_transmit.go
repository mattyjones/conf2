package comm
import (
	"yang"
	"reflect"
)

/**
 * This is used to encode YANG models. In order to navigate the YANG model it needs a model
 * which is the YANG YANG model.  It can be confusing which is the data and which is the
 * meta.
 */
type MetaTransmitter struct {
	metaRoot yang.Meta // read: meta-meta
	data yang.Meta // read: meta
	out Receiver
	level int
}

func (self *MetaTransmitter) Transmit() (err error) {
	self.out.StartTransaction()
	if err = self.transmitMeta(self.metaRoot, self.data); err == nil {
		self.out.EndTransaction()
	}
	return
}

func (self *MetaTransmitter) getValue(meta yang.HasType, obj interface{}) string {
	fieldName := yang.MetaNameToFieldName(meta.GetIdent())
	objType := reflect.ValueOf(obj).Elem()
	value := objType.FieldByName(fieldName)
	switch meta.Type() {
		case "bool":
			if value.Bool() {
				return "true"
			}
			return "false"
		default:
			return value.String()
	}
}

func (self *MetaTransmitter) transmitList(meta *yang.List, data yang.MetaList) {
	self.out.NewList(meta)
	// Do not resolve proxies in data otherwise we will recurse infinitely!
	dataItems := yang.NewMetaListIterator(data, false)
	for dataItems.HasNextMeta() {
		next := dataItems.NextMeta()
		self.transmitMeta(meta, next)
	}
	self.out.ExitList(meta)
}

func (self *MetaTransmitter) definitionType(data yang.Meta) string {
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

func (self *MetaTransmitter) transmitDefinitions(meta *yang.List, data yang.MetaList) (err error) {
	self.out.NewList(meta)
	choice := meta.GetFirstMeta().(*yang.Choice)
	dataItems := yang.NewMetaListIterator(data, false)
	for dataItems.HasNextMeta() {
		next := dataItems.NextMeta()
		caseType := self.definitionType(next)
		if itemMeta, ok := choice.GetCase(caseType).GetFirstMeta().(*yang.Container); ok {
			self.transmitObject(itemMeta, next)
		} else {
			return &commError{"Expected container meta for definition"}
		}
	}
	self.out.ExitList(meta)
	return nil
}

func (self *MetaTransmitter) transmitMeta(meta yang.Meta, data yang.Meta) (err error) {
	if data == nil {
		return
	}
	switch field := meta.(type) {
	case *yang.List:
		switch field.GetIdent() {
		case "groupings":
			self.transmitList(field, data.(yang.HasGroupings).GetGroupings())
		case "definitions":
			self.transmitDefinitions(field, data.(yang.MetaList))
		case "enumerations":
			self.transmitList(field, data.(*yang.Typedef).GetEnumerations())
		case "typedefs":
			self.transmitList(field, data.(yang.HasTypedefs).GetTypedefs())
		case "rpcs":
			self.transmitList(field, data.(*yang.Module).GetRpcs())
		case "notifications":
			self.transmitList(field, data.(*yang.Module).GetNotifications())
		}
    case *yang.Module:
		self.transmitObject(field, data)
	case *yang.Container:
		self.transmitObject(field, data)
	case *yang.Leaf:
		switch field.GetIdent() {
		case "rev-date":
			self.out.PutStringLeaf(field, data.(*yang.Module).Revision.GetIdent())
		default:
			dataValue := self.getValue(field, data)
			if dataValue != "" {
				self.out.PutStringLeaf(field, dataValue)
			}
		}
		// TODO: Support PutIntLeaf by looking at type
	case *yang.LeafList:
		// TODO: Support PutIntLeafList by looking at type
	}

	return nil
}

func (self *MetaTransmitter) transmitObject(meta yang.MetaList, data yang.Meta) (err error) {
	self.out.NewObject(meta)
	i := yang.NewMetaListIterator(meta, true)
	for i.HasNextMeta() {
		next := i.NextMeta()
		self.transmitMeta(next, data)
	}
	self.out.ExitObject(meta)
	return
}