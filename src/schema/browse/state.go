package browse
import (
	"schema"
)

type WalkState struct {
	path schema.MetaPath
	insideList bool
	key []*Value
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

func (state *WalkState) SetInsideList() {
	state.insideList = true
}

func (state *WalkState) SetKey(key []*Value) {
	//state.key = key
}

func (state *WalkState) Key() []*Value {
	return state.key
}

func (state *WalkState) IsConfig() bool {
	if hasDetails, ok := state.path.Meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(&state.path)
	}
	return true
}
