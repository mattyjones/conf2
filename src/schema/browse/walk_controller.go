package browse

import (
	"schema"
)

type WalkController interface {
	ListIterator(state *WalkState, s Selection, first bool) (hasMore bool, err error)
	ContainerIterator(state *WalkState, s Selection) schema.MetaIterator
	CloseSelection(s Selection) error
}

func WalkAll() WalkController {
	return &FullWalk{MaxDepth:32}
}

