package data

import (
	"errors"
	"fmt"
	"schema"
	"regexp"
)

type Selection struct {
	Events Events
	Node   Node
	State  *WalkState
}

func NewSelectionFromState(node Node, state *WalkState) *Selection {
	return &Selection{
		Node: node,
		Events: &EventsImpl{},
		State: state.Copy(),
	}
}

func NewSelection(node Node, meta schema.MetaList) *Selection {
	sel := &Selection{
		Node: node,
		Events: &EventsImpl{},
		State: &WalkState{},
	}
	sel.State.path.ParentPath = &schema.MetaPath{Meta: meta}
	return sel
}

func (sel *Selection) Select(node Node) *Selection {
	child := &Selection{
		Events: sel.Events,
		Node: node,
		State: &WalkState{},
	}
	child.State.path.ParentPath = &sel.State.path
	return child
}

func (sel *Selection) SelectListItem(node Node, key []*Value) *Selection {
	next := *sel
	// important flag, otherwise we recurse indefinitely
	next.State.insideList = true
	next.Node = node
	if len(key) > 0 {
		// TODO: Support compound keys
		next.State.path.Key = key[0].String()
		next.State.key = key
	}
	return &next
}

func (sel *Selection) String() string {
	if sel.Node != nil {
		nodeStr := sel.Node.String()
		if len(nodeStr) > 0 {
			return nodeStr + " " + sel.State.path.Position()
		}
	}
	return sel.State.path.Position()
}

func (sel *Selection) RequireKey(key []*Value, err error) {
	key = sel.State.Key()
	if key == nil {
		err = errors.New(fmt.Sprint("Cannot select list without key ", sel.String()))
	}
	return
}

func (sel *Selection) Fire(e Event) (err error) {
	err = sel.Node.Event(sel, e)
	if err != nil {
		return err
	}
	path := sel.State.Path().Path()
	return sel.Events.Fire(path, e)
}

func (sel *Selection) IsConfig() bool {
	if hasDetails, ok := sel.State.Path().Meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(sel.State.Path())
	}
	return true
}

func (sel *Selection) On(e Event, listener ListenFunc) *Listener {
	return sel.OnPath(e, sel.State.Path().Path(), listener)
}

func (sel *Selection) OnPath(e Event, path string, handler ListenFunc) *Listener {
	listener := &Listener{event: e, path: path, handler: handler}
	sel.Events.AddListener(listener)
	return listener
}

func (sel *Selection) OnChild(e Event, meta schema.MetaList, listener ListenFunc)  *Listener {
	fullPath := sel.State.Path().Path() + "/" + meta.GetIdent()
	return sel.OnPath(e, fullPath, listener)
}

func (sel *Selection) OnRegex(e Event, regex *regexp.Regexp, handler ListenFunc) *Listener {
	listener := &Listener{event: e, regex: regex, handler: handler}
	sel.Events.AddListener(listener)
	return listener
}

func (sel *Selection) Level() int {
	level := -1
	p := sel.State.Path()
	for p.ParentPath != nil {
		level++
		p = p.ParentPath
	}
	return level
}