package yang

type MetaIterator interface {
	NextMeta() Meta
	HasNextMeta() bool
}

type MetaListIterator struct {
	position Meta
	currentProxy MetaIterator
	resolveProxies bool
}

type EmptyInterator int

func (e EmptyInterator) HasNextMeta() bool {
	return false;
}
func (e EmptyInterator) NextMeta() Meta {
	return nil;
}

type SingletonIterator struct {
	Meta Meta
}

func (s *SingletonIterator) HasNextMeta() bool {
	return s.Meta != nil;
}
func (s *SingletonIterator) NextMeta() Meta {
	m := s.Meta
	s.Meta = nil
	return m;
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
