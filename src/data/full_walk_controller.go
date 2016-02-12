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

func (self *ControlledWalk) parseQuery(params map[string][]string) (err error) {
	if depth, found := params["depth"]; found {
		if self.MaxDepth, err = strconv.Atoi(depth[0]); err != nil {
			return
		}
	}
	return
}

func (self *ControlledWalk) CloseSelection(s *Selection) error {
	return schema.CloseResource(s)
}

func (self *ControlledWalk) maxedLevel(selection *Selection) bool {
	level := selection.path.Len()
	if self.finalDepth == 0 {
		self.finalDepth = level + self.MaxDepth
	}
	return level >= self.finalDepth
}

func (self *ControlledWalk) VisitAction(state *Selection, rpc *schema.Rpc) (*Selection, error) {
	// Not sure what a full walk would do when hitting an action, so do nothing
	return nil, nil
}

func (self *ControlledWalk) VisitContainer(sel *Selection, meta schema.MetaList) (*Selection, error) {
	r := ContainerRequest{
		Meta: meta,
	}
	childNode, err := sel.node.Select(sel, r)
	if err == nil && childNode != nil {
		return sel.SelectChild(meta, childNode), nil
	}
	return nil, err
}

func (self *ControlledWalk) ListIterator(sel *Selection, first bool) (next *Selection, err error) {
	if self.maxedLevel(sel) {
		return nil, nil
	}
	r := ListRequest{
		First: first,
		Meta: sel.path.meta.(*schema.List),
	}
	var listNode Node
	listNode, sel.path.key, err = sel.node.Next(sel, r)
	if listNode == nil || err != nil {
		return nil, err
	}
	next = sel.SelectListItem(listNode, sel.path.key)
	return
}

func (self *ControlledWalk) ContainerIterator(selection *Selection) (schema.MetaIterator, error) {
	if self.maxedLevel(selection) {
		return schema.EmptyInterator(0), nil
	}
	return schema.NewMetaListIterator(selection.path.meta, true), nil
}
