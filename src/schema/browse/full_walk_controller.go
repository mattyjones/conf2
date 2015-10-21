package browse
import (
	"schema"
	"strconv"
	"strings"
)

type ControlledWalk struct {
	MaxDepth int
	finalDepth int
	InitialKey []*Value
}

//func NewFullWalkFromPath(p *Path) *FullWalk {
//	return LimitedWalk(p.Query)
//}

func LimitedWalk(query string) *ControlledWalk {
	c := FullWalk()
	c.parseQuery(query)
	return c
}

func FullWalk() *ControlledWalk {
	return &ControlledWalk{MaxDepth:32}
}

func (p *ControlledWalk) parseQuery(q string) (err error) {
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

func (p *ControlledWalk) CloseSelection(s Selection) error {
	return schema.CloseResource(s)
}

func (e *ControlledWalk) maxedLevel(state *WalkState) bool {
	if e.finalDepth == 0 {
		e.finalDepth = state.Level() + e.MaxDepth
	}
	return state.Level() >= e.finalDepth
}

func (e *ControlledWalk) ListIterator(state *WalkState, s Selection, first bool) (next Selection, err error) {
	if e.maxedLevel(state) {
		return nil, nil
	}
	listMeta := state.SelectedMeta().(*schema.List)
	return s.Next(state, listMeta, NO_KEYS, first)
}

func (e *ControlledWalk) ContainerIterator(state *WalkState, s Selection) schema.MetaIterator {
	if e.maxedLevel(state) {
		return schema.EmptyInterator(0)
	}
	return schema.NewMetaListIterator(state.SelectedMeta(), true)
}