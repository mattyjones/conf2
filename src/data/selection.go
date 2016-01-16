package data

import (
	"errors"
	"fmt"
	"schema"
	"regexp"
	"conf2"
)

type Selection struct {
	parent *Selection
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
		State: &WalkState{
			path : NewRootPath(meta, nil),
		},
	}
	return sel
}

func (sel *Selection) Select(node Node) *Selection {
	child := &Selection{
		parent: sel,
		Events: sel.Events,
		Node: node,
		State: &WalkState{
			path : NewContainerPath(sel.State.path, sel.State.Position().(schema.MetaList)),
		},
	}
	return child
}

func (sel *Selection) SelectListItem(node Node, key []*Value) *Selection {
	next := *sel
	// important flag, otherwise we recurse indefinitely
	next.State.insideList = true
	next.Node = node
	if len(key) > 0 {
		next.State.SetKey(key)
		next.State.path = next.State.path.SetKey(key)
	}
	return &next
}

func (sel *Selection) Meta(ident string) schema.Meta {
	return schema.FindByIdent2(sel.State.SelectedMeta(), ident)
}

func (sel *Selection) String() (s string) {
	if sel.Node == nil {
		return ""
	}
	s = sel.Node.String()
	if len(s) > 0 && sel.State.Position() != nil {
		s = s + " " + sel.State.Position().GetIdent()
	}
	return
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
	return sel.Events.Fire(sel.State.path, e)
}

func (sel *Selection) IsConfig() bool {
	if hasDetails, ok := sel.State.Position().(schema.HasDetails); ok {
		return hasDetails.Details().Config(sel.State.path)
	}
	return true
}

func (sel *Selection) On(e Event, listener ListenFunc) *Listener {
	return sel.OnPath(e, sel.State.Path().String(), listener)
}

func (sel *Selection) OnPath(e Event, path string, handler ListenFunc) *Listener {
	listener := &Listener{event: e, path: path, handler: handler}
	sel.Events.AddListener(listener)
	return listener
}

func (sel *Selection) OnChild(e Event, meta schema.MetaList, listener ListenFunc)  *Listener {
	fullPath := sel.State.path.String() + "/" + meta.GetIdent()
	return sel.OnPath(e, fullPath, listener)
}

func (sel *Selection) OnRegex(e Event, regex *regexp.Regexp, handler ListenFunc) *Listener {
	listener := &Listener{event: e, regex: regex, handler: handler}
	sel.Events.AddListener(listener)
	return listener
}

func (sel *Selection) ResolveLeafref(cwd *Selection, lreafrefIdent string, value interface{}) (interface{}, error) {
	m := schema.FindByIdent2(sel.State.SelectedMeta(), lreafrefIdent).(schema.HasDataType)
	sel, prop, err := XPath{}.SelectProperty(cwd, m.GetDataType().Path)
	if err != nil {
		return nil, err
	}
	switch x := sel.State.SelectedMeta().(type) {
	case *schema.List:
		if x.KeyMeta()[0].GetIdent() != prop.GetIdent() {
			return nil, conf2.NewErrC("Cannot resolve leafrefs to non-key values", conf2.NotImplemented)
		}
		key, keyErr := SetValue(x.KeyMeta()[0].GetDataType(), value)
		if keyErr != nil {
			return nil, keyErr
		}
		found, foundErr := sel.Node.Next(sel, x, false, []*Value{key}, true)
		if foundErr != nil {
			return nil, foundErr
		}
		if found != nil {
			sel = sel.SelectListItem(found, []*Value{key})
		}
	}
	if schema.IsLeaf(prop) {
		v, err := sel.Node.Read(sel, prop.(schema.HasDataType))
		if err != nil {
			return nil, err
		}
		return v.Value(), nil
	}
	return sel.Peek(), nil
}

//func (sel *Selection) FindSelectionByPropertyXPath(path string, prop interface{}) (*Selection, error) {
//	panic("TODO")
//	return nil, nil
//}
//
//func (sel *Selection) FindByPropertyXPath(ident string, prop interface{}) (interface{}, error) {
//	m := schema.FindByIdent2(sel.State.SelectedMeta(), ident).(schema.HasDataType)
//	xpath := m.GetDataType().Path
//	hndSel, err := sel.FindSelectionByPropertyXPath(xpath, prop)
//	if err != nil {
//		return nil, err
//	}
//	return hndSel.Peek(), nil
//}
//
//func (sel *Selection) Find(path string) (*Selection, error) {
//	if strings.HasPrefix("../", path) {
//		if sel.parent != nil {
//			return sel.parent.Find(path[3:])
//		} else {
//			return nil, conf2.NewErrC("No parent path to resolve " + path, conf2.NotFound)
//		}
//	}
//	p, err := ParsePath(path, sel.State.SelectedMeta())
//	if err != nil {
//		return nil, err
//	}
//	return WalkPath(sel, p)
//}

func (sel *Selection) Peek() interface{} {
	return sel.Node.Peek(sel)
}
