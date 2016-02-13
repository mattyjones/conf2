package data

import (
	"net/url"
	"strings"
	"conf2"
)

type Selector struct {
	Selection   *Selection
	Target      *PathSlice
	EditControl *ControlledWalk
	Edit        *Editor
	LastErr     error
}

func (self Selector) Find(path string) Selector {
	return (&self).find(path, nil)
}

func (self Selector) FindUrl(url *url.URL) Selector {
	return (&self).find("", url)
}

func (self *Selector) find(path string, url *url.URL) Selector {
	if self.LastErr != nil {
		return *self
	}
	sel := self.Selection
	p := path
	if url != nil {
		p = url.Path
	}
	s := Selector{}
	for strings.HasPrefix(p, "../") {
		if sel.parent != nil {
			sel = sel.parent
			p = p[3:]
		} else {
			s.LastErr = conf2.NewErrC("No parent path to resolve " + p, conf2.NotFound)
			return s
		}
	}

	s.Target, s.LastErr = ParsePath(p, sel.Meta())
	if s.LastErr != nil {
		return s
	}
	if url != nil {
		s.Target.SetParams(url.Query())
	}
	s.EditControl = LimitedWalk(s.Target.Params())
	findController := NewFindTarget(s.Target)
	if s.LastErr = sel.Walk(findController); s.LastErr == nil {
		s.Selection = findController.Target
	}
	return s
}

func (self Selector) Push(toNode Node) Selector {
	self.Edit = &Editor{
		from: self.Selection,
		to:   self.Selection.Fork(toNode),
	}
	return self
}

func (self Selector) Pull(fromNode Node) Selector {
	self.Edit = &Editor{
		from: self.Selection.Fork(fromNode),
		to:   self.Selection,
	}
	return self
}

func (self Selector) Insert() Selector {
	(&self).edit(INSERT)
	return self
}

func (self Selector) Upsert() Selector {
	(&self).edit(UPSERT)
	return self
}

func (self Selector) Update() Selector {
	(&self).edit(UPDATE)
	return self
}

func (self *Selector) edit(op Strategy) {
	if self.LastErr != nil {
		return
	}
	if self.EditControl != nil {
		self.LastErr = self.Edit.Edit(op, self.EditControl)
	} else {
		self.LastErr = self.Edit.Edit(op, FullWalk())
	}
}
