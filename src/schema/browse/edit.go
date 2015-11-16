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
	created bool
}

func Insert(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, INSERT, FullWalk())
}

func InsertByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, INSERT)
}

func UpdateByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, UPDATE)
}

func UpsertByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, UPSERT)
}

func EditByNode(selection *Selection, src Node, dest Node, strategy Strategy) (err error) {
	e := editor{}
	var n Node
	if schema.IsList(selection.SelectedMeta()) && !selection.InsideList() {
		n, err = e.list(src, dest, false, strategy)
	} else {
		n, err = e.container(src, dest, false, strategy)
	}
	if err == nil {
		s := selection.Copy(n)
		if err = s.Fire(BEGIN_EDIT); err == nil {
			if err = Walk(s, FullWalk()); err == nil {
				err = s.Fire(END_EDIT)
			} else {
				// TODO: guard against panics not calling undo
				if suberr := s.Fire(UNDO_EDIT); suberr != nil {
					conf2.Err.Printf("Could not roll back edit, err=" + suberr.Error())
				}
			}
		}
	}
	return
}

func ControlledInsert(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, INSERT, cntrl)
}

func Upsert(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, UPSERT, FullWalk())
}

func ControlledUpsert(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, UPSERT, cntrl)
}

func Update(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, UPDATE, FullWalk())
}

func ControlledUpdate(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, UPDATE, cntrl)
}

func Delete(sel *Selection, node Node) (err error) {
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

func Action(impl *Selection, input Node) (output *Selection, err error) {
	rpc := impl.Position().(*schema.Rpc)
	return impl.Node().Action(impl, rpc, input)
}

func Edit(from *Selection, dest *Selection, strategy Strategy, controller WalkController) (err error) {
	e := editor{}
	var n Node
	if schema.IsList(from.SelectedMeta()) && !from.InsideList() {
		n, err = e.list(from.Node(), dest.Node(), false, strategy)
	} else {
		n, err = e.container(from.Node(), dest.Node(), false, strategy)
	}
	if err == nil {
		// sync dest and from
		s := from.Copy(n)
		if err = s.Fire(BEGIN_EDIT); err == nil {
			if err = Walk(s, controller); err == nil {
				err = s.Fire(END_EDIT)
			} else {
				if suberr := s.Fire(UNDO_EDIT); suberr != nil {
					conf2.Err.Printf("Could not roll back edit, err=" + suberr.Error())
				}
			}
		}
	}
	return
}

func (e *editor) list(from Node, to Node, new bool, strategy Strategy) (Node, error) {
	if to == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get target node")}
	}
	if from == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get source node")}
	}
	s := &MyNode{Label: fmt.Sprint("Edit list ", from.String(), "=>", to.String())}
	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(selection *Selection, meta *schema.List, _ bool, key []*Value, first bool) (next Node, err error) {
		var created bool
		var fromNextNode Node
		fromNextNode, err = from.Next(selection, meta, false, key, first)
		if err != nil || fromNextNode == nil {
			return nil, err
		}

		var nextKey []*Value
		var toNextNode Node
		if nextKey, err = e.loadKey(selection, key); err != nil {
			return nil, err
		}
		if len(nextKey) > 0 {
			if toNextNode, err = to.Next(selection, meta, false, nextKey, true); err != nil {
				return nil, err
			}
		}
		switch strategy {
		case UPDATE:
			if toNextNode == nil {
				msg := fmt.Sprint("No item found with given key in list ", selection.String())
				return nil, &browseError{Code: http.StatusNotFound, Msg: msg}
			}
		case UPSERT:
			if toNextNode == nil {
				if toNextNode, err = to.Next(selection, meta, true, nextKey, true); err != nil {
					return nil, err
				}
				created = true
			}
		case INSERT:
			if toNextNode != nil {
				msg := fmt.Sprint("Duplicate item found with same key in list ", selection.String())
				return nil, &browseError{Code: http.StatusConflict, Msg: msg}
			}
			if toNextNode, err = to.Next(selection, meta, true, nextKey, true); err != nil {
				return nil, err
			}
			created = true
		default:
			return nil, &browseError{Msg: "Stratgey not implmented"}
		}
		return e.container(fromNextNode, toNextNode, created, UPSERT)
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		return e.handleEvent(sel, from, to, new, event)
	}
	return s, nil
}

func (e *editor) handleEvent(sel *Selection, from Node, to Node, new bool, event Event) (err error) {
	if event == LEAVE && new {
		if err = sel.Fire(NEW); err != nil {
			return
		}
	}
	if err = to.Event(sel, event); err != nil {
		return
	}
	switch event {
	// don't send write-only events to reader
	case BEGIN_EDIT, END_EDIT, NEW, UNDO_EDIT, DELETE:
		return
	}
	if err = from.Event(sel, event); err != nil {
		return
	}
	return
}

func (e *editor) container(from Node, to Node, new bool, strategy Strategy) (Node, error) {
	if to == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get target container selection")}
	}
	if from == nil {
		return nil, &browseError{Msg: fmt.Sprint("Unable to get source node")}
	}
	s := &MyNode{Label: fmt.Sprint("Edit container ", from.String(), "=>", to.String())}
	s.OnChoose = func(state *Selection, choice *schema.Choice) (schema.Meta, error) {
		return from.Choose(state, choice)
	}
	s.OnSelect = func(selection *Selection, meta schema.MetaList, _ bool) (Node, error) {
		var created bool
		var err error
		var fromChild Node
		fromChild, err = from.Select(selection, meta, false)
		if err != nil || fromChild == nil {
			return nil, err
		}

		var toChild Node
		toChild, err = to.Select(selection, meta, false)
		if err != nil {
			return nil, err
		}
		isList := schema.IsList(meta)

		switch strategy {
		case INSERT:
			if toChild != nil {
				msg := fmt.Sprint("Found existing container ", selection.String())
				return nil, &browseError{Code: http.StatusConflict, Msg: msg}
			}
			if toChild, err = to.Select(selection, meta, true); err != nil {
				return nil, err
			}
			created = true
		case UPSERT:
			if toChild == nil {
				if toChild, err = to.Select(selection, meta, true); err != nil {
					return nil, err
				}
				created = true
			}
		case UPDATE:
			if toChild == nil {
				msg := fmt.Sprint("Container not found in list ", selection.String())
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
		return e.handleEvent(sel, from, to, new, event)
	}
	s.OnRead = func(selection *Selection, meta schema.HasDataType) (v *Value, err error) {
		if v, err = from.Read(selection, meta); err != nil {
			return
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = to.Write(selection, meta, v); err != nil {
				return
			}
		}
		return
	}

	return s, nil
}

func (e *editor) loadKey(selection *Selection, explictKey []*Value) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return selection.Key(), nil
}
