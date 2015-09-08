package browse
import (
	"yang"
)

type WalkController interface {
	ListIterator(s *Selection, level int, first bool) (hasMore bool, err error)
	ContainerIterator(s *Selection, level int) yang.MetaIterator
}

type exhaustiveController struct {
	MaxDepth int
}

func (e *exhaustiveController) ListIterator(s *Selection, level int, first bool) (hasMore bool, err error) {
	return s.Iterate([]string{}, first)
}
func (e *exhaustiveController) ContainerIterator(s *Selection, level int) yang.MetaIterator {
	if level >= e.MaxDepth {
		return yang.EmptyInterator(0)
	}
	return yang.NewMetaListIterator(s.Meta, true)
}
