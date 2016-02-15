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
	p := path
	selection := self.Selection
	if strings.HasPrefix(path, "../") {
		for strings.HasPrefix(p, "../") {
			if selection.parent != nil {
				selection = selection.parent
				p = p[3:]
			} else {
				self.LastErr = conf2.NewErrC("No parent path to resolve " + p, conf2.NotFound)
				return self
			}
		}
	}
	var u *url.URL
	u, self.LastErr = url.Parse(p)
	if self.LastErr != nil {
		return self
	}
	return self.FindUrl(u)
}

func (self Selector) FindUrl(url *url.URL) Selector {
	if self.LastErr != nil {
		return self
	}

	self.Target, self.LastErr = ParseUrlPath(url, self.Selection.Meta())
	if self.LastErr != nil {
		return self
	}
	if url != nil {
		self.Target.SetParams(url.Query())
	}
	self.EditControl = LimitedWalk(self.Target.Params())
	findController := NewFindTarget(self.Target)
	if self.LastErr = self.Selection.Walk(findController); self.LastErr == nil {
		self.Selection = findController.Target
	}
	return self
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
