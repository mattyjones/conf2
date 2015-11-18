package browse

import (
	"schema"
	"strconv"
	"strings"
)

type ControlledWalk struct {
	MaxDepth   int
	finalDepth int
	InitialKey []*Value
}

//func NewFullWalkFromPath(p *Path) *FullWalk {
//	return LimitedWalk(p.Query)
//}

func LimitedWalk(query string) *ControlledWalk {
	c := FullWalk()
	c.parseQuery(query)
	return c
}

func FullWalk() *ControlledWalk {
	return &ControlledWalk{MaxDepth: 32}
}

func (p *ControlledWalk) parseQuery(q string) (err error) {
	if len(q) == 0 {
		return nil
	}
	params := strings.Split(q, "&")
	for _, param := range params {
		nameValue := strings.Split(param, "=")
		switch nameValue[0] {
		case "depth":
			if p.MaxDepth, err = strconv.Atoi(nameValue[1]); err != nil {
				return
			}
		}
	}
	return
}

func (p *ControlledWalk) CloseSelection(s *Selection) error {
	return schema.CloseResource(s)
}

func (e *ControlledWalk) maxedLevel(selection *Selection) bool {
	if e.finalDepth == 0 {
		e.finalDepth = selection.Level() + e.MaxDepth
	}
	return selection.Level() >= e.finalDepth
}

func (n *ControlledWalk) VisitAction(state *Selection) error {
	return nil
}

func (e *ControlledWalk) ListIterator(selection *Selection, first bool) (next *Selection, err error) {
	if e.maxedLevel(selection) {
		return nil, nil
	}
	listMeta := selection.State.SelectedMeta().(*schema.List)
	var listNode Node
	listNode, err = selection.Node.Next(selection, listMeta, false, NO_KEYS, first)
	if listNode == nil || err != nil {
		return nil, err
	}
	next = selection.SelectListItem(listNode, selection.State.Key())
	return
}

func (e *ControlledWalk) ContainerIterator(selection *Selection) schema.MetaIterator {
	if e.maxedLevel(selection) {
		return schema.EmptyInterator(0)
	}
	return schema.NewMetaListIterator(selection.State.SelectedMeta(), true)
}
