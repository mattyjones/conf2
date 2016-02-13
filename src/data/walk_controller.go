package data

import (
	"schema"
)

type ContainerRequest struct {
	New bool
	Target *PathSlice
	Depth int
	First bool
	Meta schema.MetaList
}

type ListRequest struct {
	New bool
	Target *PathSlice
	Depth int
	First bool
	Meta *schema.List
	Key []*Value
}

type WalkController interface {
	ListIterator(selection *Selection, first bool) (next *Selection, err error)
	ContainerIterator(selection *Selection) (schema.MetaIterator, error)
    VisitContainer(sel *Selection, meta schema.MetaList) (child *Selection, err error)
	VisitAction(selection *Selection, rpc *schema.Rpc) (*Selection, error)
	CloseSelection(s *Selection) error
}

func WalkAll() WalkController {
	return FullWalk()
}
