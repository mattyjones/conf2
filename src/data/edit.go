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
}

func (s *Selection) To(to *Selection) *Editor {
	return &Editor{
		from : s,
		to: to,
	}
}

func (s *Selection) Push(toNode Node) *Editor {
	return &Editor{
		from : s,
		to: s.Fork(toNode),
	}
}

func (s *Selection) Pull(fromNode Node) *Editor {
	return &Editor{
		from : s.Fork(fromNode),
		to: s,
	}
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

func (self *Selection) Delete() (err error) {
	if err = self.Fire(START_TREE_EDIT); err == nil {
		if err = self.Fire(DELETE); err == nil {
			err = self.Fire(END_TREE_EDIT)
		}
	}
	return
}

func (e *Editor) Edit(strategy Strategy, controller WalkController) (err error) {
	var n Node
	if schema.IsList(e.from.path.meta) && !e.from.insideList {
		n, err = e.list(e.from, e.to, false, strategy)
	} else {
		n, err = e.container(e.from, e.to, false, strategy)
	}
	if err != nil {
		return err
	}
	// we could fork "from" or "to", shouldn't matter
	s := e.from.Fork(n)
	if err = e.to.Fire(START_TREE_EDIT); err == nil {
		if err = s.Walk(controller); err == nil {
			if err = e.to.Fire(LEAVE_EDIT); err == nil {
				err = e.to.Fire(END_TREE_EDIT)
			}
		}
	}
	return
}

func (e *Editor) list(from *Selection, to *Selection, new bool, strategy Strategy) (Node, error) {
	s := &MyNode{Label: fmt.Sprint("Edit list ", from.node.String(), "=>", to.node.String())}
	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(sel *Selection, meta *schema.List, _ bool, key []*Value, first bool) (next Node, err error) {
		var created bool
		var fromNextNode Node
		fromNextNode, err = from.node.Next(from, meta, false, key, first)
		if err != nil || fromNextNode == nil {
			return nil, err
		}

		sel.path.key = from.path.key
		var toNextNode Node
		if len(sel.path.key) > 0 {
			if toNextNode, err = to.node.Next(to, meta, false, sel.path.key, true); err != nil {
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
				if toNextNode, err = to.node.Next(to, meta, true, sel.path.key, true); err != nil {
					return nil, err
				}
				created = true
			}
		case INSERT:
			if toNextNode != nil {
				msg := fmt.Sprint("Duplicate item found with same key in list ", sel.String())
				return nil, conf2.NewErrC(msg, conf2.Conflict)
			}
			if toNextNode, err = to.node.Next(to, meta, true, sel.path.key, true); err != nil {
				return nil, err
			}
			created = true
		default:
			return nil, conf2.NewErrC("Stratgey not implmented", conf2.NotImplemented)
		}
		if err != nil {
			return nil, err
		} else  if toNextNode == nil {
			return nil, conf2.NewErr("Could not create destination list node " + to.String())
		}
		fromChild := from.SelectListItem(fromNextNode, sel.path.key)
		toChild := to.SelectListItem(toNextNode, sel.path.key)
		return e.container(fromChild, toChild, created, UPSERT)
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		return e.handleEvent(sel, from, to, new, event)
	}
	return s, nil
}

func (e *Editor) container(from *Selection, to *Selection, new bool, strategy Strategy) (Node, error) {
	s := &MyNode{Label: fmt.Sprint("Edit container ", from.node.String(), "=>", to.node.String())}
	s.OnChoose = func(sel *Selection, choice *schema.Choice) (schema.Meta, error) {
		return from.node.Choose(from, choice)
	}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, _ bool) (Node, error) {
		var created bool
		var err error
		var fromChildNode Node
		fromChildNode, err = from.node.Select(from, meta, false)
		if err != nil || fromChildNode == nil {
			return nil, err
		}

		var toChildNode Node
		toChildNode, err = to.node.Select(to, meta, false)
		if err != nil {
			return nil, err
		}
		isList := schema.IsList(meta)

		switch strategy {
		case INSERT:
			if toChildNode != nil {
				return nil, conf2.NewErrC("Found existing container " + sel.String(), conf2.Conflict)
			}
			if toChildNode, err = to.node.Select(to, meta, true); err != nil {
				return nil, err
			}
			created = true
		case UPSERT:
			if toChildNode == nil {
				if toChildNode, err = to.node.Select(to, meta, true); err != nil {
					return nil, err
				}
				created = true
			}
		case UPDATE:
			if toChildNode == nil {
				return nil, conf2.NewErrC("Container not found in list " + sel.String(), conf2.NotFound)
			}
		default:
			return nil, conf2.NewErrC("Stratgey not implemented", conf2.NotImplemented)
		}

		if err != nil {
			return nil, err
		} else if toChildNode == nil {
			return nil, conf2.NewErr("Could not create destination container node " + to.String())
		}
		// we always switch to upsert strategy because if there were any conflicts, it would have been
		// discovered in top-most level.
		fromChild := from.SelectChild(meta, fromChildNode)
		toChild := to.SelectChild(meta, toChildNode)
		if isList {
			return e.list(fromChild, toChild, created, UPSERT)
		}
		return e.container(fromChild, toChild, created, UPSERT)
	}
	s.OnEvent = func(sel *Selection, event Event) (err error) {
		return e.handleEvent(sel, from, to, new, event)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (v *Value, err error) {
		if v, err = from.node.Read(from, meta); err != nil {
			return
		}
		if v == nil && strategy != UPDATE {
			if meta.GetDataType().HasDefault() {
				v = &Value{Type:meta.GetDataType()}
				v.CoerseStrValue(meta.GetDataType().Default())
			}
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = to.node.Write(to, meta, v); err != nil {
				return
			}
		}
		return
	}

	return s, nil
}

func (e *Editor) handleEvent(sel *Selection, from *Selection, to *Selection, new bool, event Event) (err error) {
	if event == LEAVE {
		if new {
			if err = to.Fire(NEW); err != nil {
				return
			}
		}
		if err = to.Fire(LEAVE_EDIT); err != nil {
			return
		}
	}
	if err = to.node.Event(sel, event); err != nil {
		return
	}
	if err = from.node.Event(sel, event); err != nil {
		return
	}
	return
}

func (e *Editor) loadKey(selection *Selection, explictKey []*Value) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return selection.path.key, nil
}

