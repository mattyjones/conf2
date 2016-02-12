package data

import (
	"regexp"
	"schema"
	"fmt"
	"net/url"
	"conf2"
	"strings"
)

type Selection struct {
	parent     *Selection
	events     Events
	node       Node
	path       *Path
	insideList bool
}

func (self *Selection) Parent() *Selection {
	return self.parent
}

func (self *Selection) Events() Events {
	return self.events
}

func (self *Selection) Meta() schema.MetaList {
	return self.path.meta
}

func (self *Selection) Node() Node {
	return self.node
}

func (sel *Selection) Fork(node Node) *Selection {
	copy := *sel
	copy.events = &EventsImpl{}
	copy.node = node
	return &copy
}

func (sel *Selection) Key() []*Value {
	return sel.path.key
}

func (sel *Selection) String() string {
	return fmt.Sprint(sel.node.String(), ":", sel.path.String())
}

func NewSelection(meta schema.MetaList, node Node) *Selection {
	return &Selection{
		events: &EventsImpl{},
		path: &Path{meta: meta},
		node:   node,
	}
}

func (sel *Selection) SelectChild(meta schema.MetaList, node Node) *Selection {
	child := &Selection{
		parent: sel,
		events: sel.events,
		path: &Path{parent: sel.path, meta: meta},
		node:   node,
	}
	return child
}

func (sel *Selection) SelectListItem(node Node, key []*Value) *Selection {
	var parentPath *Path
	if sel.parent != nil {
		parentPath = sel.parent.path
	}
	child := &Selection{
		parent:     sel.parent, // NOTE: list item's parent is list's parent, not list!
		events:     sel.events,
		node:       node,
		path:		&Path{parent:parentPath, meta: sel.path.meta, key: key},
		insideList: true,
	}
	return child
}

func (sel *Selection) Path() *Path {
	return sel.path
}

func (sel *Selection) Fire(e Event) (err error) {
	target := sel
	for {
		err = target.node.Event(target, e)
		if err != nil {
			return err
		}
		if e.Type.Bubbles() && ! e.state.propagationStopped {
			if target.parent != nil {
				target = target.parent
				continue
			}
		}
		break
	}
	return sel.events.Fire(sel.path, e)
}

func (sel *Selection) On(e EventType, listener ListenFunc) *Listener {
	return sel.OnPath(e, sel.Path().String(), listener)
}

func (sel *Selection) OnPath(e EventType, path string, handler ListenFunc) *Listener {
	listener := &Listener{event: e, path: path, handler: handler}
	sel.events.AddListener(listener)
	return listener
}

func (sel *Selection) OnChild(e EventType, meta schema.MetaList, listener ListenFunc) *Listener {
	fullPath := sel.path.String() + "/" + meta.GetIdent()
	return sel.OnPath(e, fullPath, listener)
}

func (sel *Selection) OnRegex(e EventType, regex *regexp.Regexp, handler ListenFunc) *Listener {
	listener := &Listener{event: e, regex: regex, handler: handler}
	sel.events.AddListener(listener)
	return listener
}

func (sel *Selection) Peek(peekId string) interface{} {
	return sel.node.Peek(sel, peekId)
}

func isFwdSlash(r rune) bool {
	return r == '/'
}

func (self *Selection) FindLeaf(path string) (*Selection, schema.HasDataType, error) {
	if strings.HasPrefix(path, "../") {
		if self.parent != nil {
			return self.parent.FindLeaf(path[3:])
		} else {
			return nil, nil, conf2.NewErrC("No parent path to resolve " + path, conf2.NotFound)
		}
	}

	slash := strings.LastIndexFunc(path, isFwdSlash)
	sel := self
	ident := path
	if slash > 0 {
		var err error
		if sel, err = sel.Find(path[:slash]); err != nil {
			return nil, nil, err
		}
		ident = path[slash + 1:]
	}
	meta := schema.FindByIdent2(sel.path.meta, ident)
	return sel, meta.(schema.HasDataType), nil
}

// Like Find but panics if path not found or error parsing path
func (self *Selection) Require(path string) (*Selection) {
	sel, err := self.Find(path)
	if err != nil {
		panic(err)
	}
	return sel
}

func (self *Selection) Find(path string) (*Selection, error) {
	if strings.HasPrefix(path, "../") {
		if self.parent != nil {
			return self.parent.Find(path[3:])
		} else {
			return nil, conf2.NewErrC("No parent path to resolve " + path, conf2.NotFound)
		}
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	return self.FindUrl(u)
}

func (self *Selection) FindUrl(url *url.URL) (*Selection, error) {
	if len(url.Path) == 0 {
		return self, nil
	}
	pslice, err := ParsePath(url.Path, self.path.meta)
	if err != nil {
		return nil, err
	}
	pslice.SetParams(url.Query())
	finder := NewFindTarget(pslice)
	err = self.Walk(finder)
	return finder.Target, err
}

func (sel *Selection) Set(ident string, value interface{}) error {
	n := sel.node
	if cw, ok := n.(ChangeAwareNode); ok {
		n = cw.Changes()
	}
	pos := schema.FindByIdent2(sel.path.meta, ident)
	if pos == nil {
		return conf2.NewErrC("property not found " + ident, conf2.NotFound)
	}
	meta := pos.(schema.HasDataType)
	v, e := SetValue(meta.GetDataType(), value)
	if e != nil {
		return e
	}
	return n.Write(sel, meta, v)
}

func (sel *Selection) Get(ident string) (interface{}, error) {
	prop := schema.FindByIdent2(sel.path.meta, ident)
	if prop != nil {
		if schema.IsLeaf(prop) {
			v, err := sel.node.Read(sel, prop.(schema.HasDataType))
			if err != nil {
				return nil, err
			}
			return v.Value(), nil
		}
	}
	return nil, nil
}

func (sel *Selection) GetValue(ident string) (*Value, error) {
	prop := schema.FindByIdent2(sel.path.meta, ident)
	if prop != nil {
		v, err := sel.node.Read(sel, prop.(schema.HasDataType))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, nil
}

func (self *Selection) IsConfig(meta schema.Meta) bool {
	if hasDetails, ok := meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(self.path)
	}
	return true
}

func (sel *Selection) ClearAll() error {
	return sel.node.Event(sel, DELETE.New())
}

func (sel *Selection) FindOrCreate(ident string, autoCreate bool) (*Selection, error) {
	m := schema.FindByIdent2(sel.path.meta, ident)
	var err error
	var child Node
	if m != nil {
		r := ContainerRequest{
			Meta: m.(schema.MetaList),
		}
		child, err = sel.node.Select(sel, r)
		if err != nil {
			return nil, err
		} else if child == nil && autoCreate {
			r.New = true
			child, err = sel.node.Select(sel, r)
			if err != nil {
				return nil, err
			}
		}
		if child != nil {
			return sel.SelectChild(r.Meta, child), nil
		}
	}
	return nil, nil
}

