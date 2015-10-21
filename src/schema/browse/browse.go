package browse

import (
	"schema"
)

type Browser interface {
	Selector(path *Path, strategy Strategy) (Selection, *WalkState, error)
	Schema() (schema.MetaList)
}

func WalkPath(state *WalkState, from Selection, path *Path) (Selection, *WalkState, error) {
	finder := NewFindTarget(path)
	err := walk(state, from, finder)
	return finder.target, finder.targetState, err
}

func Walk(state *WalkState, selection Selection, controller WalkController) (err error) {
	return walk(state, selection, controller)
}

func walk(state *WalkState, selection Selection, controller WalkController) (err error) {
	if schema.IsList(state.SelectedMeta()) && !state.InsideList() {
		var next Selection
		if next, err = controller.ListIterator(state, selection, true); err != nil {
			return
		}
		for next != nil {
			listItemState := state.SelectListItem(state.Key())

			if err = walk(listItemState, next, controller); err != nil {
				return
			}
			if next, err = controller.ListIterator(state, selection, false); err != nil {
				return
			}
		}
	} else {
		var child Selection
		i := controller.ContainerIterator(state, selection)
		for i.HasNextMeta() {
			state.SetPosition(i.NextMeta())
			if choice, isChoice := state.Position().(*schema.Choice); isChoice {
				var choosen schema.Meta
				if choosen, err = selection.Choose(state, choice); err != nil {
					return
				}
				state.SetPosition(choosen)
			}
			if schema.IsLeaf(state.Position()) {
				// only walking here, not interested in value
				if _, err = selection.Read(state, state.Position().(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := state.Position().(schema.MetaList)
				child, err = selection.Select(state, metaList)
				if err != nil {
					return
				} else if child == nil {
					continue
				}
				defer schema.CloseResource(child)
				if err = walk(state.Select(), child, controller); err != nil {
					return
				}

				err = selection.Unselect(state, metaList)
			}
		}
	}
	return
}

