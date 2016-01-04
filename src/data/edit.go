package data

import (
	"fmt"
	"schema"
	"conf2"
)

type Strategy int
const (
	UPSERT Strategy = iota + 1
	INSERT
	UPDATE
)

type Editor struct{
	from *Selection
	to *Selection
	created bool
}

func NodeToNode(fromNode Node, toNode Node, schema schema.MetaList) *Editor {
	return &Editor{
		from : NewSelection(fromNode, schema),
		to: NewSelection(toNode, schema),
	}
}

func SelectionToNode(from *Selection, toNode Node) *Editor {
	return &Editor{
		from : from,
		to: NewSelectionFromState(toNode, from.State),
	}
}

func NodeToSelection(fromNode Node, to *Selection) *Editor {
	return &Editor{
		from : NewSelectionFromState(fromNode, to.State),
		to: to,
	}
}

func SelectionToSelection(from *Selection, to *Selection) *Editor {
	return &Editor{
		from : from,
		to: to,
	}
}

func NodeToPath(fromNode Node, data Data, path string) (*Editor, error) {
	var err error
	var p *PathSlice
	if p, err = ParsePath(path, data.Schema()); err != nil {
		return nil, err
	}
	var to *Selection
	if to, err = WalkPath(NewSelection(data.Node(), data.Schema()), p); err != nil {
		return nil, err
	}
	if to == nil {
		return nil, PathNotFound(path)
	}
	return NodeToSelection(fromNode, to), nil
}

func PathToNode(data Data, path string, toNode Node) (*Editor, error) {
	var err error
	var p *PathSlice
	if p, err = ParsePath(path, data.Schema()); err != nil {
		return nil, err
	}
	var from *Selection
	if from, err = WalkPath(NewSelection(data.Node(), data.Schema()), p); err != nil {
		return nil, err
	}
	if from == nil {
		return nil, PathNotFound(path)
	}
	return SelectionToNode(from, toNode), nil
}

func PathToPath(fromData Data, toData Data, path string) (*Editor, error) {
	var err error
	var p *PathSlice
	if p, err = ParsePath(path, fromData.Schema()); err != nil {
		return nil, err
	}
	from := NewSelection(fromData.Node(), fromData.Schema())
	if from, err = WalkPath(from, p); err != nil {
		return nil, err
	}
	to := NewSelection(toData.Node(), toData.Schema())
	if to, err = WalkPath(to, p); err != nil {
		return nil, err
	}
	return SelectionToSelection(from, to), nil
}

func (e *Editor) Insert() (err error) {
	return e.Edit(INSERT, FullWalk())
}

func (e *Editor) ControlledInsert(cntrl WalkController) (err error) {
	return e.Edit(INSERT, cntrl)
}

func (e *Editor) Upsert() (err error) {
	return e.Edit(UPSERT, FullWalk())
}

func (e *Editor) ControlledUpsert(cntrl WalkController) (err error) {
	return e.Edit(UPSERT, cntrl)
}

func (e *Editor) Update() (err error) {
	return e.Edit(UPDATE, FullWalk())
}

func (e *Editor) ControlledUpdate(cntrl WalkController) (err error) {
	return e.Edit(UPDATE, cntrl)
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

func (e *Editor) Edit(strategy Strategy, controller WalkController) (err error) {
	var n Node
	if schema.IsList(e.from.State.SelectedMeta()) && !e.from.State.InsideList() {
		n, err = e.list(e.from.Node, e.to.Node, false, strategy, "")
	} else {
		n, err = e.container(e.from.Node, e.to.Node, false, strategy, "")
	}
	s := &Selection{
		Events: &EventMulticast{
			A: e.from.Events,
			B: e.to.Events,
		},
		State: e.from.State,
		Node: n,
	}
	if err == nil {
		if err = e.to.Fire(BEGIN_EDIT); err == nil {
			if err = Walk(s, controller); err == nil {
				err = e.to.Fire(END_EDIT)
			} else {
				// TODO: guard against panics not calling undo
				if suberr := e.to.Fire(UNDO_EDIT); suberr != nil {
					conf2.Err.Printf("Could not roll back edit, err=" + suberr.Error())
				}
			}
		}
	}
	return
}

func (e *Editor) list(fromNode Node, toNode Node, new bool, strategy Strategy, path string) (Node, error) {
	if toNode == nil {
		return nil, conf2.NewErrC("Unable to get target node " + path, conf2.NotFound)
	}
	if fromNode == nil {
		return nil, conf2.NewErrC("Unable to get source node" + path, conf2.NotFound)
	}
	to := &Selection{
		Events: e.to.Events,
		Node: toNode,
	}
	from := &Selection{
		Events: e.from.Events,
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
				return nil, conf2.NewErrC(msg, conf2.NotFound)
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
				return nil, conf2.NewErrC(msg, conf2.Conflict)
			}
			if toNextNode, err = toNode.Next(to, meta, true, nextKey, true); err != nil {
				return nil, err
			}
			created = true
		default:
			return nil, conf2.NewErrC("Stratgey not implmented", conf2.NotImplemented)
		}
		return e.container(fromNextNode, toNextNode, created, UPSERT, sel.State.String())
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		to.State = sel.State
		from.State = sel.State
		return e.handleEvent(sel, from, to, new, event)
	}
	return s, nil
}

func (e *Editor) container(fromNode Node, toNode Node, new bool, strategy Strategy, path string) (Node, error) {
	if toNode == nil {
		return nil, conf2.NewErrC("Unable to get target container selection " + path, conf2.NotFound)
	}
	if fromNode == nil {
		return nil, conf2.NewErrC("Unable to get source node" + path, conf2.NotFound)
	}
	to := &Selection{
		Events: e.to.Events,
		Node: toNode,
	}
	from := &Selection{
		Events: e.from.Events,
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
				return nil, conf2.NewErrC("Found existing container " + sel.String(), conf2.Conflict)
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
				return nil, conf2.NewErrC("Container not found in list " + sel.String(), conf2.NotFound)
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
			return e.list(fromChild, toChild, created, UPSERT, sel.State.String())
		}
		return e.container(fromChild, toChild, created, UPSERT, sel.State.String())
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

func (e *Editor) handleEvent(sel *Selection, from *Selection, to *Selection, new bool, event Event) (err error) {
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

func (e *Editor) loadKey(selection *Selection, explictKey []*Value) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return selection.State.Key(), nil
}
