package browse
import (
	"schema"
	"fmt"
	"errors"
)

type WalkState struct {
	path schema.MetaPath
	key []*Value
	insideList bool
}

func NewWalkState(meta schema.MetaList) *WalkState {
	state := &WalkState{}
	state.path.ParentPath = &schema.MetaPath{Meta:meta}
	return state
}

func (state *WalkState) SelectedMeta() schema.MetaList {
	return state.path.Parent()
}

func (state *WalkState) Select() *WalkState {
	child := &WalkState{}
	child.path.ParentPath = &state.path
	return child
}

func (state *WalkState) SelectListItem(key []*Value) *WalkState {
	next := *state
	// important flag, otherwise we recurse indefinitely
	next.insideList = true
	if len(key) > 0 {
		// TODO: Support compound keys
		next.path.Key = key[0].String()
		next.key = key
	}
	return &next
}

func (state *WalkState) Position() schema.Meta {
	return state.path.Meta
}

func (state *WalkState) SetPosition(position schema.Meta)  {
	state.path.Meta = position
}

func (state *WalkState) Path() *schema.MetaPath {
	return &state.path
}

func (state *WalkState) String() string {
	return state.path.String()
}

func (state *WalkState) InsideList() bool {
	return state.insideList
}

func (state *WalkState) Key() []*Value {
	return state.key
}

func (state *WalkState) RequireKey() ([]*Value, error) {
	if state.key == nil {
		return nil, errors.New(fmt.Sprint("Cannot select list without key ", state.String()))
	}
	return state.key, nil
}

func (state *WalkState) SetKey(key []*Value) {
	state.key = key
}

func (state *WalkState) IsConfig() bool {
	if hasDetails, ok := state.path.Meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(&state.path)
	}
	return true
}

func (state *WalkState) Level() int {
	level := -1
	p := &state.path
	for p.ParentPath != nil {
		level++
		p = p.ParentPath
	}
	return level
}
