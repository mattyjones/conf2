package browse
import (
	"schema"
	"strconv"
	"strings"
)

type FullWalk struct {
	MaxDepth int
	finalDepth int
	InitialKey []*Value
}

func NewFullWalkFromPath(p *Path) *FullWalk {
	return NewFullWalk(p.query)
}

func NewFullWalk(query string) *FullWalk {
	c := &FullWalk{MaxDepth:32}
	c.parseQuery(query)
	return c
}

func (p *FullWalk) parseQuery(q string) (err error) {
	if len(q) == 0 {
		return nil
	}
	params := strings.Split(q, "&")
	for _, param := range params {
		nameValue := strings.Split(param, "=")
		switch nameValue[0] {
		case "depth":
			if p.MaxDepth, err = strconv.Atoi(nameValue[1]); err != nil {
				return
			}
		}
	}
	return
}

func (p *FullWalk) CloseSelection(s Selection) error {
	return schema.CloseResource(s)
}

func (e *FullWalk) maxedLevel(state *WalkState) bool {
	if e.finalDepth == 0 {
		e.finalDepth = state.Level() + e.MaxDepth
	}
	return state.Level() >= e.finalDepth
}

func (e *FullWalk) ListIterator(state *WalkState, s Selection, first bool) (hasMore bool, err error) {
	if e.maxedLevel(state) {
		return false, nil
	}
	listMeta := state.SelectedMeta().(*schema.List)
	return s.Next(state, listMeta, NO_KEYS, first)
}

func (e *FullWalk) ContainerIterator(state *WalkState, s Selection) schema.MetaIterator {
	if e.maxedLevel(state) {
		return schema.EmptyInterator(0)
	}
	return schema.NewMetaListIterator(state.SelectedMeta(), true)
}