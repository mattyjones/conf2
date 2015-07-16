package yang

type MetaIterator interface {
	NextMeta() Meta
	HasNextMeta() bool
}
type MetaProxy interface {
	ResolveProxy() MetaIterator
}
type MetaListIterator struct {
	position Meta
	currentProxy MetaIterator
	resolveProxies bool
}

func NewMetaListIterator(list MetaList, resolveProxies bool) MetaIterator {
	return &MetaListIterator{position:list.GetFirstMeta(), resolveProxies:resolveProxies}
}

func (self *MetaListIterator) HasNextMeta() bool {
	return self.position != nil && self.currentProxy != nil
}

func (self *MetaListIterator) NextMeta() Meta {
	for self.position != nil || self.currentProxy != nil {
		if self.currentProxy != nil {
			next := self.currentProxy.NextMeta()
			if next != nil {
				return next
			} else {
				self.currentProxy = nil
			}
		} else {
			if self.resolveProxies {
				proxy, isProxy := self.position.(MetaProxy)
				if !isProxy {
					next := self.position
					self.position = self.position.GetSibling()
					return next
				} else {
					self.position = self.position.GetSibling()
					self.currentProxy = proxy.ResolveProxy();
				}
			} else {
				next := self.position
				self.position = self.position.GetSibling()
				return next
			}
		}
	}
	return nil
}
