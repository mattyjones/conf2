package browse

import (
	"errors"
	"fmt"
	"schema"
)

type FindTarget struct {
	path     *Path
	Target   *Selection
	resource schema.Resource
}

func NewFindTarget(p *Path) *FindTarget {
	return &FindTarget{path: p}
}

func (n *FindTarget) ListIterator(selection *Selection, first bool) (next *Selection, err error) {
	if !first {
		// when we're finding targets, we never iterate more than one item in a list
		return nil, nil
	}
	level := selection.Level()
	if level == len(n.path.Segments) {
		if len(n.path.Segments[level-1].Keys) == 0 {
			n.setTarget(selection)
			return nil, nil
		}
	}

	if len(n.path.Segments[level-1].Keys) == 0 {
		return nil, errors.New("Key required when navigating lists")
	}
	list := selection.SelectedMeta().(*schema.List)
	var key []*Value
	key, err = CoerseKeys(list, n.path.Segments[level-1].Keys)
	if err != nil {
		return nil, err
	}
	selection.SetKey(key)
	var nextNode Node
	if nextNode, err = selection.Node().Next(selection, list, key, true); err != nil {
		return nil, err
	} else if nextNode == nil {
		return nil, err
	}
	next = selection.SelectListItem(nextNode, key)
	if level == len(n.path.Segments) {
		//state.SetInsideList()
		n.setTarget(next)
	}
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
	level := selection.Level()
	if level+1 != len(n.path.Segments) {
		return errors.New(fmt.Sprint("Target is an action or rpc ", selection.String()))
	}
	n.setTarget(selection)
	return nil
}

func (n *FindTarget) ContainerIterator(selection *Selection) schema.MetaIterator {
	level := selection.Level()
	if level == len(n.path.Segments) {
		n.setTarget(selection)
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdentExpandChoices(selection.SelectedMeta(), n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta: position}
}
