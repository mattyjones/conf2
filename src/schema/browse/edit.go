package browse

import (
	"fmt"
	"net/http"
	"schema"
	"conf2"
)

type Strategy int
const (
	UPSERT Strategy = iota + 1
	INSERT
	UPDATE
)

type editor struct{
	fromEvents Events
	toEvents Events
	created bool
}

func SyncData(from Data, to Data, p *Path, s Strategy) (err error) {
	var fromSel, toSel *Selection
	if fromSel, err = from.Selector(p); err != nil {
		return err
	}
	if toSel, err = to.Selector(p); err != nil {
		return err
	}
	return Edit(fromSel, toSel, s, FullWalk())
}

func Insert(from *Selection, to *Selection) (err error) {
	return Edit(from, to, INSERT, FullWalk())
}

func ControlledInsert(from *Selection, to *Selection, cntrl WalkController) (err error) {
	return Edit(from, to, INSERT, cntrl)
}

func Upsert(from *Selection, to *Selection) (err error) {
	return Edit(from, to, UPSERT, FullWalk())
}

func ControlledUpsert(from *Selection, to *Selection, cntrl WalkController) (err error) {
	return Edit(from, to, UPSERT, cntrl)
}

func Update(from *Selection, to *Selection) (err error) {
	return Edit(from, to, UPDATE, FullWalk())
}

func ControlledUpdate(from *Selection, to *Selection, cntrl WalkController) (err error) {
	return Edit(from, to, UPDATE, cntrl)
}

func Delete(sel *Selection) (err error) {
	if err = sel.Fire(BEGIN_EDIT); err != nil {
		return err
	}
	if err = sel.Fire(DELETE); err != nil {
		if suberr := sel.Fire(UNDO_EDIT); suberr != nil {
			conf2.Err.Printf("Could not roll back edit, err=" + suberr.Error())
		}
		return err
	}
	err = sel.Fire(END_EDIT)
	return
}

func Action(impl *Selection, input *Selection) (output *Selection, err error) {
	rpc := impl.State.Position().(*schema.Rpc)
	return impl.Node.Action(impl, rpc, input)
}

func Edit(from *Selection, to *Selection, strategy Strategy, controller WalkController) (err error) {
	e := editor{
		fromEvents: &EventsImpl{Parent:from.Events},
		toEvents: &EventsImpl{Parent:to.Events},
	}
	var n Node
	if schema.IsList(from.State.SelectedMeta()) && !from.State.InsideList() {
		n, err = e.list(from.Node, to.Node, false, strategy)
	} else {
		n, err = e.container(from.Node, to.Node, false, strategy)
	}
	s := &Selection{
		Events: &EventMulticast{
			A: e.fromEvents,
			B: e.toEvents,
		},
		State: from.State,
		Node: n,
	}
	if err == nil {
		if err = to.Fire(BEGIN_EDIT); err == nil {
			if err = Walk(s, controller); err == nil {
				err = to.Fire(END_EDIT)
			} else {
				// TODO: guard against panics not calling undo
				if suberr := to.Fire(UNDO_EDIT); suberr != nil {
					conf2.Err.Printf("Could not roll back edit, err=" + suberr.Error())
				}
			}
		}
	}
	return
}

func (e *editor) list(fromNode Node, toNode Node, new bool, strategy Strategy) (Node, error) {
	if toNode == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get target node")}
	}
	if fromNode == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get source node")}
	}
	to := &Selection{
		Events: e.toEvents,
		Node: toNode,
	}
	from := &Selection{
		Events: e.fromEvents,
		Node: fromNode,
	}
	s := &MyNode{Label: fmt.Sprint("Edit list ", fromNode.String(), "=>", toNode.String())}
	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(sel *Selection, meta *schema.List, _ bool, key []*Value, first bool) (next Node, err error) {
		to.State = sel.State
		from.State = sel.State
		var created bool
		var fromNextNode Node
		fromNextNode, err = fromNode.Next(from, meta, false, key, first)
		if err != nil || fromNextNode == nil {
			return nil, err
		}

		var nextKey []*Value
		var toNextNode Node
		if nextKey, err = e.loadKey(sel, key); err != nil {
			return nil, err
		}
		if len(nextKey) > 0 {
			if toNextNode, err = toNode.Next(to, meta, false, nextKey, true); err != nil {
				return nil, err
			}
		}
		switch strategy {
		case UPDATE:
			if toNextNode == nil {
				msg := fmt.Sprint("No item found with given key in list ", sel.String())
				return nil, &browseError{Code: http.StatusNotFound, Msg: msg}
			}
		case UPSERT:
			if toNextNode == nil {
				if toNextNode, err = toNode.Next(to, meta, true, nextKey, true); err != nil {
					return nil, err
				}
				created = true
			}
		case INSERT:
			if toNextNode != nil {
				msg := fmt.Sprint("Duplicate item found with same key in list ", sel.String())
				return nil, &browseError{Code: http.StatusConflict, Msg: msg}
			}
			if toNextNode, err = toNode.Next(to, meta, true, nextKey, true); err != nil {
				return nil, err
			}
			created = true
		default:
			return nil, &browseError{Msg: "Stratgey not implmented"}
		}
		return e.container(fromNextNode, toNextNode, created, UPSERT)
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		to.State = sel.State
		from.State = sel.State
		return e.handleEvent(sel, from, to, new, event)
	}
	return s, nil
}

func (e *editor) container(fromNode Node, toNode Node, new bool, strategy Strategy) (Node, error) {
	if toNode == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get target container selection")}
	}
	if fromNode == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get source node")}
	}
//conf2.Debug.Printf("container %s, new %v, pathPtr=%p", state.Path().Path(), new, state.Path())
	to := &Selection{
		Events: e.toEvents,
		Node: toNode,
	}
	from := &Selection{
		Events: e.fromEvents,
		Node: fromNode,
	}
	s := &MyNode{Label: fmt.Sprint("Edit container ", fromNode.String(), "=>", toNode.String())}
	s.OnChoose = func(sel *Selection, choice *schema.Choice) (schema.Meta, error) {
		from.State = sel.State
		return from.Node.Choose(from, choice)
	}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, _ bool) (Node, error) {
		to.State = sel.State
		from.State = sel.State
		var created bool
		var err error
		var fromChild Node
		fromChild, err = fromNode.Select(from, meta, false)
		if err != nil || fromChild == nil {
			return nil, err
		}

		var toChild Node
		toChild, err = toNode.Select(to, meta, false)
		if err != nil {
			return nil, err
		}
		isList := schema.IsList(meta)

		switch strategy {
		case INSERT:
			if toChild != nil {
				msg := fmt.Sprint("Found existing container ", sel.String())
				return nil, &browseError{Code: http.StatusConflict, Msg: msg}
			}
			if toChild, err = toNode.Select(to, meta, true); err != nil {
				return nil, err
			}
			created = true
		case UPSERT:
			if toChild == nil {
				if toChild, err = toNode.Select(to, meta, true); err != nil {
					return nil, err
				}
				created = true
			}
		case UPDATE:
			if toChild == nil {
				msg := fmt.Sprint("Container not found in list ", sel.String())
				return nil, &browseError{Code: http.StatusNotFound, Msg: msg}
			}
		default:
			return nil, &browseError{Msg: "Stratgey not implmented"}
		}

		if err != nil {
			return nil, err
		}
		// we always switch to upsert strategy because if there were any conflicts, it would have been
		// discovered in top-most level.
		if isList {
			return e.list(fromChild, toChild, created, UPSERT)
		}
		return e.container(fromChild, toChild, created, UPSERT)
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		to.State = sel.State
		from.State = sel.State
		return e.handleEvent(sel, from, to, new, event)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (v *Value, err error) {
		to.State = sel.State
		from.State = sel.State
		if v, err = fromNode.Read(from, meta); err != nil {
			return
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = toNode.Write(to, meta, v); err != nil {
				return
			}
		}
		return
	}

	return s, nil
}

func (e *editor) handleEvent(sel *Selection, from *Selection, to *Selection, new bool, event Event) (err error) {
	if event == LEAVE && new {
		if err = to.Fire(NEW); err != nil {
			return
		}
	}
	if err = to.Node.Event(sel, event); err != nil {
		return
	}
	if err = from.Node.Event(sel, event); err != nil {
		return
	}
	return
}

func (e *editor) loadKey(selection *Selection, explictKey []*Value) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return selection.State.Key(), nil
}
