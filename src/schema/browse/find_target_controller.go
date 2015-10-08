package browse
import (
	"schema"
	"errors"
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

func (n *FindTarget) ListIterator(state *WalkState, s Selection, first bool) (hasMore bool, err error) {
	if !first {
		// when we're finding targets, we never iterate more than one item in a list
		return false, nil
	}
	level := state.Level()
	if level == len(n.path.Segments) {
		n.setTarget(state, s)
		if len(n.path.Segments[level - 1].Keys) == 0 {
			return false, nil
		}
	}

	if len(n.path.Segments[level - 1].Keys) == 0 {
		return false, errors.New("Key required when navigating lists")
	}
	list := state.SelectedMeta().(*schema.List)
	var key []*Value
	key, err = CoerseKeys(list, n.path.Segments[level - 1].Keys)
	if err != nil {
		return false, err
	}
	state.SetKey(key)
	return s.Next(state, list, key, true)
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

func (n *FindTarget) ContainerIterator(state *WalkState, s Selection) schema.MetaIterator {
	//n.path.Key = nil
	level := state.Level()
	if level == len(n.path.Segments) {
		n.setTarget(state, s)
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdentExpandChoices(state.SelectedMeta(), n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta:position}
}
