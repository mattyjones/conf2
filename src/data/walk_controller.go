package data

import (
	"schema"
)

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
