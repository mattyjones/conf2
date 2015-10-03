package browse
import (
	"schema"
	"fmt"
)

type FindTarget struct {
	path *Path
	target Selection
	targetState *WalkState
	resource schema.Resource
}

func NewFindTarget(p *Path) *FindTarget {
	return &FindTarget{path:p}
}

func (n *FindTarget) ListIterator(state *WalkState, s Selection, level int, first bool) (hasMore bool, err error) {
	if level == len(n.path.Segments) {
		if len(n.path.Segments[level - 1].Keys) == 0 {
			n.setTarget(state, s)
			return false, nil
		}
		if !first {
			n.setTarget(state, s)
			return false, nil
		}
	}
	if first && level > 0 && level <= len(n.path.Segments) {
		segment := n.path.Segments[level - 1]
		keysAsStrings := segment.Keys
		list, isList := state.SelectedMeta().(*schema.List)
		if !isList {
			return false, &browseError{Msg:fmt.Sprintf("Key \"%s\" specified when not a list", keysAsStrings)}
		}
		var key []*Value
		key, err = CoerseKeys(list, keysAsStrings)
		state.SetKey(key)
		if err != nil {
			return false, err
		}
		return s.Next(state, list, key, first)
	} else {
		return false, nil
	}
}

func (p *FindTarget) CloseSelection(s Selection) error {
	if s != p.target {
		return schema.CloseResource(s)
	}
	return nil
}

func (n *FindTarget) setTarget(state *WalkState, s Selection) {
	n.target = s
	n.targetState = state
	// we take ownership of resource so it's not released until target is used
	//	n.resource = s.Resource
	//	s.Resource = nil
}

func (n *FindTarget) ContainerIterator(state *WalkState, s Selection, level int) schema.MetaIterator {
	//n.path.Key = nil
	if level >= len(n.path.Segments) {
		n.setTarget(state, s)
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdentExpandChoices(state.SelectedMeta(), n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta:position}
}
