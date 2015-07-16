package yang

type DefIterator interface {
	NextDef() Def
	HasNextDef() bool
}
type DefProxy interface {
	ResolveProxy() DefIterator
}
type DefListIterator struct {
	position Def
	currentProxy DefIterator
	resolveProxies bool
}

func NewDefListIterator(list DefList, resolveProxies bool) DefIterator {
	return &DefListIterator{position:list.GetFirstDef(), resolveProxies:resolveProxies}
}

func (self *DefListIterator) HasNextDef() bool {
	return self.position != nil && self.currentProxy != nil
}

func (self *DefListIterator) NextDef() Def {
	for self.position != nil || self.currentProxy != nil {
		if self.currentProxy != nil {
			next := self.currentProxy.NextDef()
			if next != nil {
				return next
			} else {
				self.currentProxy = nil
			}
		} else {
			if self.resolveProxies {
				proxy, isProxy := self.position.(DefProxy)
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
