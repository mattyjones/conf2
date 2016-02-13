package data

import (
	"conf2"
	"errors"
	"schema"
)

type FindTarget struct {
	path       *PathSlice
	position   *Path
	Target     *Selection
	resource   schema.Resource
	autocreate bool
	firedFind  bool
}

func NewFindTarget(p *PathSlice) *FindTarget {
	finder := &FindTarget{
		path:     p,
		position: p.NextAfter(p.Head),
	}
	_, finder.autocreate = p.Head.Params()["autocreate"]
	return finder
}

func (self *FindTarget) ListIterator(selection *Selection, first bool) (next *Selection, err error) {
	if !first {
		// when we're finding targets, we never iterate more than one item in a list
		return nil, nil
	}
	if self.position == self.path.Tail && len(self.position.Key()) == 0 {
		self.setTarget(selection)
		return nil, nil
	}
	if len(self.position.Key()) == 0 {
		return nil, errors.New("Key required when navigating lists")
	}
	var nextNode Node
	r := ListRequest{
		Target: self.path,
		Meta: selection.path.meta.(*schema.List),
		First: true,
		Key: self.position.Key(),
	}
	nextNode, selection.path.key, err = selection.node.Next(selection, r)
	if err != nil {
		return nil, err
	} else if nextNode == nil {
		if self.autocreate {
			r.New = true
			nextNode, selection.path.key, err = selection.node.Next(selection, r)
			if err != nil {
				return nil, err
			} else if nextNode == nil {
				return nil, conf2.NewErr("Could not autocreate list item for " + selection.path.String())
			}
		} else {
			return nil, nil
		}
	}
	next = selection.SelectListItem(nextNode, self.position.Key())
	if self.position == self.path.Tail {
		self.setTarget(selection)
	}
	self.position = self.path.NextAfter(self.position)
	return
}

func (self *FindTarget) CloseSelection(s *Selection) error {
	if s != self.Target {
		return schema.CloseResource(s)
	}
	return nil
}

func (self *FindTarget) setTarget(selection *Selection) {
	self.Target = selection

	// we take ownership of resource so it's not released until target is used
	//	n.resource = s.Resource
	//	s.Resource = nil
}

func (self *FindTarget) VisitAction(selection *Selection, rpc *schema.Rpc) (*Selection, error) {
	actionSel := selection.SelectChild(rpc, selection.node)
	self.setTarget(actionSel)
	return actionSel, nil
}

func (self *FindTarget) VisitContainer(sel *Selection, meta schema.MetaList) (*Selection, error) {
	r := ContainerRequest{
		Target: self.path,
		Meta: meta,
	}
	childNode, err := sel.node.Select(sel, r)
	if err != nil {
		return nil, err
	}
	if childNode == nil {
		if !self.autocreate {
			return nil, nil
		}
		r.New = true
		childNode, err = sel.node.Select(sel, r)
		if err != nil || childNode == nil {
			return nil, err
		}
	}
	return sel.SelectChild(meta, childNode), nil
}

func (self *FindTarget) ContainerIterator(selection *Selection) (schema.MetaIterator, error) {
	if self.position == nil {
		self.setTarget(selection)
		return schema.EmptyInterator(0), nil
	}

	i := &schema.SingletonIterator{Meta: self.position.Meta()}
	if !schema.IsList(self.position.Meta()) {
		self.position = self.path.NextAfter(self.position)
	}
	return i, nil
}
