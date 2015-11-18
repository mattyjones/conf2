package browse
import (
	"schema"
)

type WalkState struct {
	path       schema.MetaPath
	key        []*Value
	insideList bool
}

func (state *WalkState) SelectedMeta() schema.MetaList {
	return state.path.Parent()
}

func (state *WalkState) Position() schema.Meta {
	return state.path.Meta
}

func (state *WalkState) SetPosition(position schema.Meta) {
	state.path.Meta = position
}

func (state *WalkState) Path() *schema.MetaPath {
	return &state.path
}

func (state *WalkState) InsideList() bool {
	return state.insideList
}

func (state *WalkState) Key() []*Value {
	return state.key
}

func (state *WalkState) SetKey(key []*Value) {
	state.key = key
}

func (state *WalkState) Copy() *WalkState {
	copy := *state
	return &copy
}

