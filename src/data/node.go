package data

import (
	"fmt"
	"net/http"
	"schema"
)

type Node interface {
	fmt.Stringer
	Select(sel *Selection, meta schema.MetaList, new bool) (child Node, err error)
//	Find(sel *Selection, path *Path) error
	Next(sel *Selection, meta *schema.List, new bool, keys []*Value, isFirst bool) (next Node, err error)
	Read(sel *Selection, meta schema.HasDataType) (*Value, error)
	Write(sel *Selection, meta schema.HasDataType, val *Value) error
	Choose(sel *Selection, choice *schema.Choice) (m schema.Meta, err error)
	Event(sel *Selection, e Event) error
	Action(sel *Selection, meta *schema.Rpc, input *Selection) (output Node, err error)
	Peek(sel *Selection, peekId string) interface{}
}

// A way to direct changes to another node to enable CopyOnWrite or other persistable options
type ChangeAwareNode interface {
	DirectChanges(config Node)
	Changes() Node
}

type MyNode struct {
	Label        string
	Peekables    map[string]interface{}
	ChangeAccess Node
	OnNext       NextFunc
	OnSelect     SelectFunc
	OnRead       ReadFunc
	OnWrite      WriteFunc
	OnChoose     ChooseFunc
	OnAction     ActionFunc
	OnEvent      EventFunc
//	OnFind       FindFunc
	OnPeek       PeekFunc
	Resource     schema.Resource
}

func (n *MyNode) DirectChanges(changeNode Node) {
	n.ChangeAccess = changeNode
}

func (n *MyNode) Changes() Node {
	// If there's a change interceptor set, use it otherwise
	// changes go directly back to node
	if n.ChangeAccess != nil {
		return n.ChangeAccess
	}
	return n
}

func (s *MyNode) String() string {
	return s.Label
}

func (s *MyNode) Close() (err error) {
	if s.Resource != nil {
		err = s.Resource.Close()
		s.Resource = nil
	}
	return
}

func (s *MyNode) Select(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
	if s.OnSelect == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Select not implemented on node ", sel.String()),
		}
	}
	return s.OnSelect(sel, meta, new)
}

func (s *MyNode) Next(sel *Selection, meta *schema.List, new bool, keys []*Value, isFirst bool) (Node, error) {
	if s.OnNext == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Next not implemented on node ", sel.String()),
		}
	}
	return s.OnNext(sel, meta, new, keys, isFirst)
}

func (s *MyNode) Read(sel *Selection, meta schema.HasDataType) (*Value, error) {
	if s.OnRead == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Read not implemented on node ", sel.String()),
		}
	}
	return s.OnRead(sel, meta)
}

func (s *MyNode) Write(sel *Selection, meta schema.HasDataType, val *Value) error {
	if s.OnWrite == nil {
		return &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Write not implemented on node ", sel.String()),
		}
	}
	//fmt.Printf("select OnWrite - %s %s\n", op.String(), sel.String())
	return s.OnWrite(sel, meta, val)
}

func (s *MyNode) Choose(sel *Selection, choice *schema.Choice) (m schema.Meta, err error) {
	if s.OnChoose == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Choose not implemented on node ", sel.String()),
		}
	}
	return s.OnChoose(sel, choice)
}

func (s *MyNode) Action(sel *Selection, meta *schema.Rpc, input *Selection) (output Node, err error) {
	if s.OnAction == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg:  fmt.Sprint("Action not implemented on node ", sel.String()),
		}
	}
	return s.OnAction(sel, meta, input)
}

func (s *MyNode) Event(sel *Selection, e Event) (err error) {
	if s.OnEvent != nil {
		return s.OnEvent(sel, e)
	}
	return nil
}

//func (s *MyNode) Find(sel *Selection, p *Path) (err error) {
//	if s.OnFind != nil {
//		return s.OnFind(sel, p)
//	}
//	return nil
//}

func (s *MyNode) Peek(sel *Selection, peekId string) interface{} {
	if s.OnPeek != nil {
		return s.OnPeek(sel, peekId)
	}
	return s.Peekables[peekId]
}

// Useful when you want to return an error from Data.Node().  Any call to get data
// will return same error
//
// func (d *MyData) Node {
//    return ErrorNode(errors.New("bang"))
// }
type ErrorNode struct {
	Err error
}

func (e ErrorNode) Error() string {
	return e.Err.Error()
}

func (e ErrorNode) String() string {
	return e.Error()
}

func (e ErrorNode) Select(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
	return nil, e.Err
}

func (e ErrorNode) Next(*Selection, *schema.List, bool, []*Value, bool) (Node, error) {
	return nil, e.Err
}

func (e ErrorNode) Read(*Selection, schema.HasDataType) (*Value, error) {
	return nil, e.Err
}

func (e ErrorNode) Write(*Selection, schema.HasDataType, *Value) error {
	return e.Err
}

func (e ErrorNode) Choose(*Selection, *schema.Choice) (schema.Meta, error) {
	return nil, e.Err
}

func (e ErrorNode) Event(*Selection, Event) error {
	return e.Err
}

func (e ErrorNode) Action(*Selection, *schema.Rpc, *Selection) (Node, error) {
	return nil, e.Err
}

//func (e ErrorNode) Find(*Selection, *Path) error {
//	return e.Err
//}

func (e ErrorNode) Peek(sel *Selection, peekId string) interface{} {
	return nil
}

type NextFunc func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error)
type SelectFunc func(sel *Selection, meta schema.MetaList, new bool) (child Node, err error)
type ReadFunc func(sel *Selection, meta schema.HasDataType) (*Value, error)
type WriteFunc func(sel *Selection, meta schema.HasDataType, val *Value) error
type ChooseFunc func(sel *Selection, choice *schema.Choice) (m schema.Meta, err error)
type ActionFunc func(sel *Selection, rpc *schema.Rpc, input *Selection) (output Node, err error)
//type FindFunc func(sel *Selection, path *Path) error
type EventFunc func(sel *Selection, e Event) error
type PeekFunc func(sel *Selection, peekId string) interface{}
