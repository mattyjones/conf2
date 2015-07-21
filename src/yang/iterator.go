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
	if self.position != nil {
		return true
	}
	if self.currentProxy != nil {
		return self.currentProxy.HasNextMeta()
	}
	return false
}

func (self *MetaListIterator) NextMeta() Meta {
	for self.HasNextMeta() {
		if self.currentProxy != nil {
			if self.currentProxy.HasNextMeta() {
				return self.currentProxy.NextMeta()
			}
			self.currentProxy = nil
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
