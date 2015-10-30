package browse
import (
	"schema"
	"fmt"
	"errors"
)

type Selection struct {
	path schema.MetaPath
	node Node
	key []*Value
	insideList bool
}

//func (s *Selection) Write(op Operation) error {
//	return s.node.Write(s, s.SelectedMeta(), op, nil)
//}
//
//func (s *Selection) WriteValue(op Operation, v *Value) error {
//	return s.node.Write(s, s.Position(), UPDATE_VALUE, v)
//}

func (s *Selection) Node() Node {
	return s.node
}

//func (s *Selection) Next(op Operation, meta schema.MetaList, v *Value) (*Selection, error) {
//	n, err := s.node.Next(s, meta, op, v)
//	if err == nil && n != nil{
//		return s.SelectListItem(n, s.Key())
//	}
//	return nil, err
//}

//func (s *Selection) Choose(choice *schema.Choice) (schema.Meta, error) {
//	return s.node.Choose(s, choice)
//}

func NewSelection(node Node, meta schema.MetaList) *Selection {
	state := &Selection{node:node}
	state.path.ParentPath = &schema.MetaPath{Meta:meta}
	return state
}

func (state *Selection) Copy(node Node) *Selection {
	copy := *state
	copy.node = node
	return &copy
}

func (state *Selection) SelectedMeta() schema.MetaList {
	return state.path.Parent()
}

func (state *Selection) Select(node Node) *Selection {
	child := &Selection{node : node}
	child.path.ParentPath = &state.path
	return child
}

func (state *Selection) SelectListItem(node Node, key []*Value) *Selection {
	next := *state
	// important flag, otherwise we recurse indefinitely
	next.insideList = true
	next.node = node
	if len(key) > 0 {
		// TODO: Support compound keys
		next.path.Key = key[0].String()
		next.key = key
	}
	return &next
}

func (state *Selection) Position() schema.Meta {
	return state.path.Meta
}

func (state *Selection) SetPosition(position schema.Meta)  {
	state.path.Meta = position
}

func (state *Selection) Path() *schema.MetaPath {
	return &state.path
}

func (state *Selection) String() string {
	return state.path.String()
}

func (state *Selection) InsideList() bool {
	return state.insideList
}

func (state *Selection) Key() []*Value {
	return state.key
}

func (state *Selection) RequireKey() ([]*Value, error) {
	if state.key == nil {
		return nil, errors.New(fmt.Sprint("Cannot select list without key ", state.String()))
	}
	return state.key, nil
}

func (state *Selection) SetKey(key []*Value) {
	state.key = key
}

func (state *Selection) IsConfig() bool {
	if hasDetails, ok := state.path.Meta.(schema.HasDetails); ok {
		return hasDetails.Details().Config(&state.path)
	}
	return true
}

func (state *Selection) Level() int {
	level := -1
	p := &state.path
	for p.ParentPath != nil {
		level++
		p = p.ParentPath
	}
	return level
}
