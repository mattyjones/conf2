package browse

import (
	"schema"
)

type WalkController interface {
	ListIterator(selection *Selection, first bool) (next *Selection, err error)
	ContainerIterator(selection *Selection) schema.MetaIterator
	VisitAction(selection *Selection) error
	CloseSelection(s *Selection) error
}

func WalkAll() WalkController {
	return FullWalk()
}

