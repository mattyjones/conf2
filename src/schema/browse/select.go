package browse

import (
	"schema"
	"fmt"
)

type WalkState struct {
	Meta schema.MetaList
	Position schema.Meta
	InsideList bool
}

type Selection interface {
	Select(meta schema.MetaList) (Selection, error)
	Next(keys []*Value, isFirst bool) (hasMore bool, err error)
	Read(meta schema.HasDataType) (*Value, error)
	Write(meta schema.Meta, op Operation, val *Value) (error)
	Choose(choice *schema.Choice) (m schema.Meta, err error)
	Unselect(meta schema.MetaList) error
	WalkState() *WalkState
}

type MySelection struct {
	OnNext NextFunc
	OnSelect SelectFunc
	OnRead ReadFunc
	OnWrite WriteFunc
	OnUnselect UnselectFunc
	OnChoose ChooseFunc
	Resource schema.Resource
	State WalkState
}

func (s *MySelection) Close() (err error) {
	if s.Resource != nil {
		err = s.Resource.Close()
		s.Resource = nil
	}
	return
}

func (s *MySelection) Select(meta schema.MetaList) (Selection, error) {
	if s.OnSelect == nil {
		return nil, &browseError{
			Code: NOT_IMPLEMENTED,
			Msg: fmt.Sprint("Select not implemented on node ", s.ToString()),
		}
	}
	return s.OnSelect(meta)
}

func (s *MySelection) Unselect(meta schema.MetaList) error {
	if s.OnUnselect != nil {
		return s.OnUnselect(meta)
	}
	return nil
}

func (s *MySelection) Next(keys []*Value, isFirst bool) (bool, error) {
	if s.OnNext == nil {
		return false, &browseError{
			Code:NOT_IMPLEMENTED,
			Msg: fmt.Sprint("Next not implemented on node ", s.ToString()),
		}
	}
	return s.OnNext(keys, isFirst)
}

func (s *MySelection) Read(meta schema.HasDataType) (*Value, error) {
	if s.OnRead == nil {
		return nil, &browseError{
			Code: NOT_IMPLEMENTED,
			Msg: fmt.Sprint("Read not implemented on node ", s.ToString()),
		}
	}
	return s.OnRead(meta)
}

func (s *MySelection) Write(meta schema.Meta, op Operation, val *Value) error {
	if s.OnWrite == nil {
		return &browseError{
			Code: NOT_IMPLEMENTED,
			Msg: fmt.Sprint("Write not implemented on node ", s.ToString()),
		}
	}
	return s.OnWrite(meta, op, val)
}

func (s *MySelection) Choose(choice *schema.Choice) (m schema.Meta, err error) {
	if s.OnChoose == nil {
		return nil, &browseError{
			Code:NOT_IMPLEMENTED,
			Msg: fmt.Sprint("Choose not implemented on node ", s.ToString()),
		}
	}
	return s.OnChoose(choice)
}

func (s *MySelection) ToString() string {
	if s.State.Meta != nil {
		if s.State.Position != nil {
			return fmt.Sprintf("%s.%s", s.State.Meta.GetIdent(), s.State.Position.GetIdent())
		}
		return s.State.Meta.GetIdent()
	}
	return "<no meta set>"
}

func (s *MySelection) WalkState() *WalkState {
	return &s.State
}

type NextFunc func(keys []*Value, first bool) (hasMore bool, err error)
type SelectFunc func(meta schema.MetaList) (child Selection, err error)
type ReadFunc func(meta schema.HasDataType) (*Value, error)
type WriteFunc func(meta schema.Meta, op Operation, val *Value) (error)
type UnselectFunc func(meta schema.MetaList) (error)
type ChooseFunc func(choice *schema.Choice) (m schema.Meta, err error)
