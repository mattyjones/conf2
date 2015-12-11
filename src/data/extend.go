package data
import (
	"schema"
	"fmt"
)

// Used when you want to alter the response from a Node (and the nodes it creates)
// You can alters, reads, writes, event and event the child nodes it creates.

type Extend struct {
	Label    string
	Node     Node
	OnNext   ExtendNextFunc
	OnSelect ExtendSelectFunc
	OnRead   ExtendReadFunc
	OnWrite  ExtendWriteFunc
	OnChoose ExtendChooseFunc
	OnAction ExtendActionFunc
	OnFind   ExtendFindFunc
	OnEvent  ExtendEventFunc
	OnExtend ExtendFunc
}

func (e *Extend) String() string {
	return fmt.Sprintf("(%s) <- %s", e.Node.String(), e.Label)
}

func (e *Extend) Select(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
	var err error
	var child Node
	if e.OnSelect == nil {
		child, err = e.Node.Select(sel, meta, new)
	} else {
		child, err = e.OnSelect(e.Node, sel, meta, new)
	}
	if child == nil || err != nil {
		return child, err
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e.Node, child)
	} else {
		extendedChild := *e
		extendedChild.Node = child
		child = &extendedChild
	}
	return child, err
}

func (e *Extend) Next(sel *Selection, meta *schema.List, new bool, key []*schema.Value, first bool) (Node, error) {
	var err error
	var child Node
	if e.OnNext == nil {
		child, err = e.Node.Next(sel, meta, new, key, first)
	} else {
		child, err = e.OnNext(e.Node, sel, meta, new, key, first)
	}
	if child == nil || err != nil {
		return child, err
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e.Node, child)
	} else {
		extendedChild := *e
		extendedChild.Node = child
		child = &extendedChild
	}
	return child, err
}

func (e *Extend) Read(sel *Selection, meta schema.HasDataType) (*schema.Value, error) {
	if e.OnRead == nil {
		return e.Node.Read(sel, meta)
	} else {
		return e.OnRead(e.Node, sel, meta)
	}
}

func (e *Extend) Write(sel *Selection, meta schema.HasDataType, v *schema.Value) (error) {
	if e.OnWrite == nil {
		return e.Node.Write(sel, meta, v)
	} else {
		return e.OnWrite(e.Node, sel, meta, v)
	}
}

func (e *Extend) Choose(sel *Selection, choice *schema.Choice) (schema.Meta, error) {
	if e.OnWrite == nil {
		return e.Node.Choose(sel, choice)
	} else {
		return e.OnChoose(e.Node, sel, choice)
	}
}

func (e *Extend) Action(sel *Selection, meta *schema.Rpc, input Node) (output Node, err error) {
	if e.OnAction == nil {
		return e.Node.Action(sel, meta, input)
	} else {
		return e.OnAction(e.Node, sel, meta, input)
	}
}

func (e *Extend) Event(sel *Selection, event Event) (err error) {
	if e.OnEvent == nil {
		return e.Node.Event(sel, event)
	} else {
		return e.OnEvent(e.Node, sel, event)
	}
}

func (e *Extend) Find(sel *Selection, p *schema.Path) (err error) {
	if e.OnFind == nil {
		return e.Node.Find(sel, p)
	} else {
		return e.OnFind(e.Node, sel, p)
	}
}

type ExtendNextFunc func(parent Node, sel *Selection, meta *schema.List, new bool, key []*schema.Value, first bool) (next Node, err error)
type ExtendSelectFunc func(parent Node, sel *Selection, meta schema.MetaList, new bool) (child Node, err error)
type ExtendReadFunc func(parent Node, sel *Selection, meta schema.HasDataType) (*schema.Value, error)
type ExtendWriteFunc func(parent Node, sel *Selection, meta schema.HasDataType, val *schema.Value) error
type ExtendChooseFunc func(parent Node, sel *Selection, choice *schema.Choice) (m schema.Meta, err error)
type ExtendActionFunc func(parent Node, sel *Selection, rpc *schema.Rpc, input Node) (output Node, err error)
type ExtendFindFunc func(parent Node, sel *Selection, path *schema.Path) error
type ExtendEventFunc func(parent Node, sel *Selection, e Event) error
type ExtendFunc func(parent Node, child Node) (Node, error)
