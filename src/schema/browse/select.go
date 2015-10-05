package browse

import (
	"schema"
	"fmt"
	"net/http"
)


type Selection interface {
	Select(state *WalkState, meta schema.MetaList) (Selection, error)
	Next(state *WalkState, meta *schema.List, keys []*Value, isFirst bool) (hasMore bool, err error)
	Read(state *WalkState, meta schema.HasDataType) (*Value, error)
	Write(state *WalkState, meta schema.Meta, op Operation, val *Value) (error)
	Choose(state *WalkState, choice *schema.Choice) (m schema.Meta, err error)
	Unselect(state *WalkState, meta schema.MetaList) error
	Action(state *WalkState, meta *schema.Rpc) (input Selection, output Selection, err error)
}

type MySelection struct {
	OnNext NextFunc
	OnSelect SelectFunc
	OnRead ReadFunc
	OnWrite WriteFunc
	OnUnselect UnselectFunc
	OnChoose ChooseFunc
	OnAction ActionFunc
	Resource schema.Resource
}

func (s *MySelection) Close() (err error) {
	if s.Resource != nil {
		err = s.Resource.Close()
		s.Resource = nil
	}
	return
}

func (s *MySelection) Select(state *WalkState, meta schema.MetaList) (Selection, error) {
	if s.OnSelect == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Select not implemented on node ", state.String()),
		}
	}
	return s.OnSelect(state, meta)
}

func (s *MySelection) Unselect(state *WalkState, meta schema.MetaList) error {
	if s.OnUnselect != nil {
		return s.OnUnselect(state, meta)
	}
	return nil
}

func (s *MySelection) Next(state *WalkState, meta *schema.List, keys []*Value, isFirst bool) (bool, error) {
	if s.OnNext == nil {
		return false, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Next not implemented on node ", state.String()),
		}
	}
	return s.OnNext(state, meta, keys, isFirst)
}

func (s *MySelection) Read(state *WalkState, meta schema.HasDataType) (*Value, error) {
	if s.OnRead == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Read not implemented on node ", state.String()),
		}
	}
	return s.OnRead(state, meta)
}

func (s *MySelection) Write(state *WalkState, meta schema.Meta, op Operation, val *Value) error {
	if s.OnWrite == nil {
		return &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Write not implemented on node ", state.String()),
		}
	}
	return s.OnWrite(state, meta, op, val)
}

func (s *MySelection) Choose(state *WalkState, choice *schema.Choice) (m schema.Meta, err error) {
	if s.OnChoose == nil {
		return nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Choose not implemented on node ", state.String()),
		}
	}
	return s.OnChoose(state, choice)
}

func (s *MySelection) Action(state *WalkState, rpc *schema.Rpc) (input Selection, output Selection, err error) {
	if s.OnAction == nil {
		return nil, nil, &browseError{
			Code: http.StatusNotImplemented,
			Msg: fmt.Sprint("Action not implemented on node ", state.String()),
		}
	}
	return s.OnAction(state, rpc)
}

type NextFunc func(state *WalkState, meta *schema.List, keys []*Value, first bool) (hasMore bool, err error)
type SelectFunc func(state *WalkState, meta schema.MetaList) (child Selection, err error)
type ReadFunc func(state *WalkState, meta schema.HasDataType) (*Value, error)
type WriteFunc func(state *WalkState, meta schema.Meta, op Operation, val *Value) (error)
type UnselectFunc func(state *WalkState, meta schema.MetaList) (error)
type ChooseFunc func(state *WalkState, choice *schema.Choice) (m schema.Meta, err error)
type ActionFunc func(state *WalkState, rpc *schema.Rpc) (input Selection, output Selection, err error)
