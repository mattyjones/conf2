package browse

import (
	"schema"
	"fmt"
	"net/http"
)

type Node interface {
	fmt.Stringer
	Select(state *Selection, meta schema.MetaList) (Node, error)
	Next(state *Selection, meta *schema.List, keys []*Value, isFirst bool) (next Node, err error)
	Read(state *Selection, meta schema.HasDataType) (*Value, error)
	Write(state *Selection, meta schema.Meta, op Operation, val *Value) (error)
	Choose(state *Selection, choice *schema.Choice) (m schema.Meta, err error)
	Unselect(state *Selection, meta schema.MetaList) error
	Action(state *Selection, meta *schema.Rpc, input Node) (output *Selection, err error)
}

type MyNode struct {
	Label string
	OnNext NextFunc
	OnSelect SelectFunc
	OnRead ReadFunc
	OnWrite WriteFunc
	OnUnselect UnselectFunc
	OnChoose ChooseFunc
	OnAction ActionFunc
	Resource schema.Resource
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

func (s *MyNode) Select(state *Selection, meta schema.MetaList) (Node, error) {
	if s.OnSelect == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Select not implemented on node ", state.String()),
		}
	}
	return s.OnSelect(state, meta)
}

func (s *MyNode) Unselect(state *Selection, meta schema.MetaList) error {
	if s.OnUnselect != nil {
		return s.OnUnselect(state, meta)
	}
	return nil
}

func (s *MyNode) Next(state *Selection, meta *schema.List, keys []*Value, isFirst bool) (Node, error) {
	if s.OnNext == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Next not implemented on node ", state.String()),
		}
	}
	return s.OnNext(state, meta, keys, isFirst)
}

func (s *MyNode) Read(state *Selection, meta schema.HasDataType) (*Value, error) {
	if s.OnRead == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Read not implemented on node ", state.String()),
		}
	}
	return s.OnRead(state, meta)
}

func (s *MyNode) Write(state *Selection, meta schema.Meta, op Operation, val *Value) error {
	if s.OnWrite == nil {
		return &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Write not implemented on node ", state.String()),
		}
	}
//fmt.Printf("select OnWrite - %s %s\n", op.String(), state.String())
	return s.OnWrite(state, meta, op, val)
}

func (s *MyNode) Choose(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
	if s.OnChoose == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Choose not implemented on node ", state.String()),
		}
	}
	return s.OnChoose(state, choice)
}

func (s *MyNode) Action(state *Selection, meta *schema.Rpc, input Node) (output *Selection, err error) {
	if s.OnAction == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Action not implemented on node ", state.String()),
		}
	}
	return s.OnAction(state, meta, input)
}

func (my *MyNode) Mixin(delegate Node) {
	my.OnAction = delegate.Action
	my.OnSelect = delegate.Select
	my.OnUnselect = delegate.Unselect
	my.OnNext = delegate.Next
	my.OnRead = delegate.Read
	my.OnWrite = delegate.Write
	my.OnChoose = delegate.Choose
}

type NextFunc func(selection *Selection, meta *schema.List, key []*Value, first bool) (next Node, err error)
type SelectFunc func(selection *Selection, meta schema.MetaList) (child Node, err error)
type ReadFunc func(selection *Selection, meta schema.HasDataType) (*Value, error)
type WriteFunc func(selection *Selection, meta schema.Meta, op Operation, val *Value) (error)
type UnselectFunc func(selection *Selection, meta schema.MetaList) (error)
type ChooseFunc func(selection *Selection, choice *schema.Choice) (m schema.Meta, err error)
type ActionFunc func(selection *Selection, rpc *schema.Rpc, input Node) (output *Selection, err error)
