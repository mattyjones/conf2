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
	OnEvent  ExtendEventFunc
	OnExtend ExtendFunc
	OnPeek ExtendPeekFunc
}

func (e *Extend) String() string {
	return fmt.Sprintf("(%s) <- %s", e.Node.String(), e.Label)
}

func (e *Extend) Select(sel *Selection, r ContainerRequest) (Node, error) {
	var err error
	var child Node
	if e.OnSelect == nil {
		child, err = e.Node.Select(sel, r)
	} else {
		child, err = e.OnSelect(e.Node, sel, r)
	}
	if child == nil || err != nil {
		return child, err
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e, sel, r.Meta, child)
	}
	return child, err
}

func (e *Extend) Next(sel *Selection, r ListRequest) (Node, []*Value, error) {
	var err error
	var child Node
	var key []*Value
	if e.OnNext == nil {
		child, key, err = e.Node.Next(sel, r)
	} else {
		child, key, err = e.OnNext(e.Node, sel, r)
	}
	if child == nil || err != nil {
		return child, key, err
	}
	if e.OnExtend != nil {
		child, err = e.OnExtend(e, sel, r.Meta, child)
	}
	return child, key, err
}

func (e *Extend) Extend(n Node) Node {
	extendedChild := *e
	extendedChild.Node = n
	return &extendedChild
}

func (e *Extend) Read(sel *Selection, meta schema.HasDataType) (*Value, error) {
	if e.OnRead == nil {
		return e.Node.Read(sel, meta)
	} else {
		return e.OnRead(e.Node, sel, meta)
	}
}

func (e *Extend) Write(sel *Selection, meta schema.HasDataType, v *Value) (error) {
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

func (e *Extend) Action(sel *Selection, meta *schema.Rpc, input *Selection) (output Node, err error) {
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

func (e *Extend) Peek(sel *Selection, peekId string) interface{} {
	if e.OnPeek == nil {
		if found := e.Node.Peek(sel, peekId); found != nil {
			return found
		}
	}
	return e.OnPeek(e.Node, sel, peekId)
}

type ExtendNextFunc func(parent Node, sel *Selection, r ListRequest) (next Node, key []*Value, err error)
type ExtendSelectFunc func(parent Node, sel *Selection, r ContainerRequest) (child Node, err error)
type ExtendReadFunc func(parent Node, sel *Selection, meta schema.HasDataType) (*Value, error)
type ExtendWriteFunc func(parent Node, sel *Selection, meta schema.HasDataType, val *Value) error
type ExtendChooseFunc func(parent Node, sel *Selection, choice *schema.Choice) (m schema.Meta, err error)
type ExtendActionFunc func(parent Node, sel *Selection, rpc *schema.Rpc, input *Selection) (output Node, err error)
type ExtendEventFunc func(parent Node, sel *Selection, e Event) error
type ExtendFunc func(e *Extend, sel *Selection, meta schema.MetaList, child Node) (Node, error)
type ExtendPeekFunc func(parent Node, sel *Selection, peekId string) interface{}
