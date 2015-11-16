package browse

import (
	"errors"
	"fmt"
	"schema"
	"regexp"
)

type Selection struct {
	events *Events
	path       schema.MetaPath
	node       Node
	key        []*Value
	insideList bool
}

func (s *Selection) Node() Node {
	return s.node
}

func NewSelection(node Node, meta schema.MetaList) *Selection {
	sel := &Selection{
		node: node,
		events: &Events{},
	}
	sel.path.ParentPath = &schema.MetaPath{Meta: meta}
	return sel
}

func (sel *Selection) Copy(node Node) *Selection {
	copy := *sel
	copy.node = node
	return &copy
}

func (sel *Selection) SelectedMeta() schema.MetaList {
	return sel.path.Parent()
}

func (sel *Selection) Select(node Node) *Selection {
	child := &Selection{
		events: sel.events,
		node: node,
	}
	child.path.ParentPath = &sel.path
	return child
}

func (sel *Selection) SelectListItem(node Node, key []*Value) *Selection {
	next := *sel
	// important flag, otherwise we recurse indefinitely
	next.insideList = true
	next.node = node
	if len(key) > 0 {
		// TODO: Support compound keys
		next.path.Key = key[0].String()
		next.key = key
	}
	return &next
}

func (sel *Selection) Position() schema.Meta {
	return sel.path.Meta
}

func (sel *Selection) SetPosition(position schema.Meta) {
	sel.path.Meta = position
}

func (sel *Selection) Path() *schema.MetaPath {
	return &sel.path
}

func (sel *Selection) String() string {
	if sel.Node() != nil {
		nodeStr := sel.Node().String()
		if len(nodeStr) > 0 {
			return nodeStr + " " + sel.path.Position()
		}
	}
	return sel.path.Position()
}

func (sel *Selection) InsideList() bool {
	return sel.insideList
}

func (sel *Selection) Key() []*Value {
	return sel.key
}

func (sel *Selection) RequireKey() ([]*Value, error) {
	if sel.key == nil {
		return nil, errors.New(fmt.Sprint("Cannot select list without key ", sel.String()))
	}
	return sel.key, nil
}

func (sel *Selection) SetKey(key []*Value) {
	sel.key = key
}

func (sel *Selection) IsConfig() bool {
	if hasDetails, ok := sel.path.Meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(&sel.path)
	}
	return true
}

func (sel *Selection) On(e Event, listener ListenFunc) {
	sel.events.AddByFullPath(e, sel.path.Path(), listener)
}

func (sel *Selection) OnPath(e Event, path string, listener ListenFunc) {
	var fullPath string
	if len(path) == 0 || path[0] != '/' {
		fullPath = sel.path.Path() + "/" + path
	} else {
		fullPath = path
	}
	sel.events.AddByFullPath(e, fullPath, listener)
}

func (sel *Selection) OnRegex(e Event, regex *regexp.Regexp, listener ListenFunc) {
	sel.events.AddByRegex(e, regex, listener)
}

func (sel *Selection) OnChild(e Event, meta schema.MetaList, listener ListenFunc) {
	fullPath := sel.path.Path() + "/" + meta.GetIdent()
	sel.events.AddByFullPath(e, fullPath, listener)
}

func (sel *Selection) Fire(e Event) (err error) {
	err = sel.node.Event(sel, e)
	if err != nil {
		return err
	}
	path := sel.path.Path()
	return sel.events.Fire(path, e)
}

func (sel *Selection) Level() int {
	level := -1
	p := &sel.path
	for p.ParentPath != nil {
		level++
		p = p.ParentPath
	}
	return level
}
