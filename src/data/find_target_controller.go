package data

import (
	"errors"
	"schema"
)

type FindTarget struct {
	path     *PathSlice
	position *Path
	Target   *Selection
	resource schema.Resource
}

func NewFindTarget(p *PathSlice) *FindTarget {
	return &FindTarget{
		path: p,
		position: p.NextAfter(p.Head),
	}
}

func (n *FindTarget) ListIterator(selection *Selection, first bool) (next *Selection, err error) {
	if !first {
		// when we're finding targets, we never iterate more than one item in a list
		return nil, nil
	}
	if n.position == n.path.Tail && len(n.position.Key()) == 0 {
		n.setTarget(selection)
		return nil, nil
	}
	if len(n.position.Key()) == 0 {
		return nil, errors.New("Key required when navigating lists")
	}
	list := selection.State.SelectedMeta().(*schema.List)
	var nextNode Node
	if err = selection.Node.Find(selection, n.path.Tail); err != nil {
		return nil, err
	}
	nextNode, err = selection.Node.Next(selection, list, false, n.position.Key(), true)
	if err != nil || nextNode == nil {
		return nil, err
	}
	next = selection.SelectListItem(nextNode, n.position.Key())
	if n.position == n.path.Tail {
		n.setTarget(next)
	}
	n.position = n.path.NextAfter(n.position)
	return
}

func (p *FindTarget) CloseSelection(s *Selection) error {
	if s != p.Target {
		return schema.CloseResource(s)
	}
	return nil
}

func (n *FindTarget) setTarget(selection *Selection) {
	n.Target = selection

	// we take ownership of resource so it's not released until target is used
	//	n.resource = s.Resource
	//	s.Resource = nil
}

func (n *FindTarget) VisitAction(selection *Selection) error {
	n.setTarget(selection)
	return nil
}

func (n *FindTarget) ContainerIterator(selection *Selection) (schema.MetaIterator, error) {
	var err error
	if n.position == nil {
		n.setTarget(selection)
		return schema.EmptyInterator(0), nil
	}
	// should we shorten path to be path[position...tail] ?
	if err = selection.Node.Find(selection, n.path.Tail); err != nil {
		return nil, err
	}
	i := &schema.SingletonIterator{Meta: n.position.Meta()}
	if ! schema.IsList(n.position.Meta()) {
		n.position = n.path.NextAfter(n.position)
	}
	return i, nil
}
