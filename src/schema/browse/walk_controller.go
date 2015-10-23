package browse

import (
	"schema"
)

type WalkController interface {
	ListIterator(state *WalkState, s Selection, first bool) (next Selection, err error)
	ContainerIterator(state *WalkState, s Selection) schema.MetaIterator
	VisitAction(state *WalkState, s Selection) error
	CloseSelection(s Selection) error
}

func WalkAll() WalkController {
	return FullWalk()
}

