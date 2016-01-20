package data

import (
	"schema"
	"strconv"
)

type ControlledWalk struct {
	MaxDepth   int
	finalDepth int
	InitialKey []*Value
}

func LimitedWalk(params map[string][]string) *ControlledWalk {
	c := FullWalk()
	c.parseQuery(params)
	return c
}

func FullWalk() *ControlledWalk {
	return &ControlledWalk{MaxDepth: 32}
}

func (p *ControlledWalk) parseQuery(params map[string][]string) (err error) {
	if depth, found := params["depth"]; found {
		if p.MaxDepth, err = strconv.Atoi(depth[0]); err != nil {
			return
		}
	}
	return
}

func (p *ControlledWalk) CloseSelection(s *Selection) error {
	return schema.CloseResource(s)
}

func (e *ControlledWalk) maxedLevel(selection *Selection) bool {
	level := selection.path.Len()
	if e.finalDepth == 0 {
		e.finalDepth = level + e.MaxDepth
	}
	return level >= e.finalDepth
}

func (n *ControlledWalk) VisitAction(state *Selection, rpc *schema.Rpc) (*Selection, error) {
	// Not sure what a full walk would do when hitting an action, so do nothing
	return nil, nil
}

func (n *ControlledWalk) VisitContainer(sel *Selection, meta schema.MetaList) (*Selection, error) {
	childNode, err := sel.node.Select(sel, meta, false)
	if err == nil && childNode != nil {
		return sel.SelectChild(meta, childNode), nil
	}
	return nil, err
}

func (e *ControlledWalk) ListIterator(selection *Selection, first bool) (next *Selection, err error) {
	if e.maxedLevel(selection) {
		return nil, nil
	}
	listMeta := selection.path.meta.(*schema.List)
	var listNode Node
	listNode, err = selection.node.Next(selection, listMeta, false, NO_KEYS, first)
	if listNode == nil || err != nil {
		return nil, err
	}
	next = selection.SelectListItem(listNode, selection.path.key)
	return
}

func (e *ControlledWalk) ContainerIterator(selection *Selection) (schema.MetaIterator, error) {
	if e.maxedLevel(selection) {
		return schema.EmptyInterator(0), nil
	}
	return schema.NewMetaListIterator(selection.path.meta, true), nil
}
