package browse
import (
	"schema"
	"errors"
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

func (n *FindTarget) ListIterator(state *WalkState, s Selection, first bool) (next Selection, err error) {
	if !first {
		// when we're finding targets, we never iterate more than one item in a list
		return nil, nil
	}
	level := state.Level()
	if level == len(n.path.Segments) {
		if len(n.path.Segments[level - 1].Keys) == 0 {
			n.setTarget(state, s)
			return nil, nil
		}
	}

	if len(n.path.Segments[level - 1].Keys) == 0 {
		return nil, errors.New("Key required when navigating lists")
	}
	list := state.SelectedMeta().(*schema.List)
	var key []*Value
	key, err = CoerseKeys(list, n.path.Segments[level - 1].Keys)
	if err != nil {
		return nil, err
	}
	state.SetKey(key)
	if next, err = s.Next(state, list, key, true); err != nil {
		return nil, err
	}
	if level == len(n.path.Segments) {
		//state.SetInsideList()
		n.setTarget(state, next)
	}
	return
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

func (n *FindTarget) VisitAction(state *WalkState, s Selection) error {
	level := state.Level()
	if level + 1 != len(n.path.Segments) {
		return errors.New(fmt.Sprint("Target is an action or rpc ", state.String()))
	}
fmt.Printf("find_target_controller - VisitAction state=%s\n", state.String())
	n.setTarget(state, s)
	return nil
}

func (n *FindTarget) ContainerIterator(state *WalkState, s Selection) schema.MetaIterator {
	level := state.Level()
	if level == len(n.path.Segments) {
		n.setTarget(state, s)
		return schema.EmptyInterator(0)
	}
	position := schema.FindByIdentExpandChoices(state.SelectedMeta(), n.path.Segments[level].Ident)
	return &schema.SingletonIterator{Meta:position}
}
