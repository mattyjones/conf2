package browse
import (
	"schema"
	"strconv"
	"strings"
)

type FullWalk struct {
	MaxDepth int
	InitialKey []*Value
}

func NewFullWalkFromPath(p *Path) *FullWalk {
	c := &FullWalk{MaxDepth:32}
	c.parseQuery(p.query)
	return c
}

func NewFullWalk(query string) *FullWalk {
	c := &FullWalk{MaxDepth:32}
	c.parseQuery(query)
	return c
}

func (p *FullWalk) parseQuery(q string) (err error) {
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

func (e *FullWalk) ListIterator(state *WalkState, s Selection, level int, first bool) (hasMore bool, err error) {
	if level >= e.MaxDepth {
		return false, nil
	}
	listMeta := state.SelectedMeta().(*schema.List)
	return s.Next(state, listMeta, state.key, first)
}

func (e *FullWalk) ContainerIterator(state *WalkState, s Selection, level int) schema.MetaIterator {
	if level >= e.MaxDepth {
		return schema.EmptyInterator(0)
	}
	return schema.NewMetaListIterator(state.SelectedMeta(), true)
}