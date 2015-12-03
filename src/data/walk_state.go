package data
import (
	"schema"
)

type WalkState struct {
	path       *schema.Path
	position   schema.Meta
	key			[]*schema.Value
	insideList bool
}

func (state *WalkState) SelectedMeta() schema.MetaList {
	return state.path.Meta()
}

func (state *WalkState) Position() schema.Meta {
	return state.position
}

func (state *WalkState) SetPosition(position schema.Meta) {
	state.position = position
}

func (state *WalkState) Path() *schema.Path {
	return state.path
}

func (state *WalkState) String() string {
	if state.position == nil {
		return state.path.String()
	}
	return state.path.String() + "." + state.position.GetIdent()
}

func (state *WalkState) InsideList() bool {
	return state.insideList
}

func (state *WalkState) Key() []*schema.Value {
	return state.key
}

func (state *WalkState) SetKey(key []*schema.Value) {
	state.key = key
}

func (state *WalkState) Copy() *WalkState {
	copy := *state
	return &copy
}

