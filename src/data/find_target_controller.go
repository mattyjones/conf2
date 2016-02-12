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
	list := selection.path.meta.(*schema.List)
	var nextNode Node
	if err = n.fireFindOnFirst(selection); err != nil {
		return nil, err
	}
	nextNode, err = selection.node.Next(selection, list, false, n.position.Key(), true)
	if err != nil {
		return nil, err
	} else if nextNode == nil {
		if n.autocreate {
			nextNode, err = selection.node.Next(selection, list, true, n.position.Key(), true)
			if err != nil {
				return nil, err
			} else if nextNode == nil {
				return nil, conf2.NewErr("Could not autocreate list item for " + selection.path.String())
			}
		} else {
			return nil, nil
		}
	}
	next = selection.SelectListItem(nextNode, n.position.Key())
	if n.position == n.path.Tail {
		n.setTarget(selection)
	}
	n.position = n.path.NextAfter(n.position)
	return
}

func (p *FindTarget) fireFindOnFirst(sel *Selection) error {
	if p.firedFind {
		return nil
	}
	// should we shorten path to be path[position...tail] ?
	p.firedFind = true
	return sel.Fire(FETCH_TREE.NewWithDetails(&FetchDetails{Path: p.path.Tail}))
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

func (n *FindTarget) VisitAction(selection *Selection, rpc *schema.Rpc) (*Selection, error) {
	actionSel := selection.SelectChild(rpc, selection.node)
	n.setTarget(actionSel)
	return actionSel, nil
}

func (n *FindTarget) VisitContainer(sel *Selection, meta schema.MetaList) (*Selection, error) {
	childNode, err := sel.node.Select(sel, meta, false)
	if err != nil {
		return nil, err
	}
	if childNode == nil {
		if !n.autocreate {
			return nil, nil
		}
		childNode, err = sel.node.Select(sel, meta, true)
		if err != nil || childNode == nil {
			return nil, err
		}
	}
	return sel.SelectChild(meta, childNode), nil
}

func (n *FindTarget) ContainerIterator(selection *Selection) (schema.MetaIterator, error) {
	var err error
	if n.position == nil {
		n.setTarget(selection)
		return schema.EmptyInterator(0), nil
	}

	if err = n.fireFindOnFirst(selection); err != nil {
		return nil, err
	}

	i := &schema.SingletonIterator{Meta: n.position.Meta()}
	if !schema.IsList(n.position.Meta()) {
		n.position = n.path.NextAfter(n.position)
	}
	return i, nil
}
