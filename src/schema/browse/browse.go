package browse

import (
	"schema"
)

type Browser interface {
	Selector(path *Path, strategy Strategy) (Selection, *WalkState, error)
	Module() (*schema.Module)
}

func WalkPath(state *WalkState, from Selection, path *Path) (Selection, *WalkState, error) {
	finder := NewFindTarget(path)
	err := walk(state, from, finder, 0)
	return finder.target, finder.targetState, err
}

func Walk(state *WalkState, selection Selection, controller WalkController) (err error) {
	return walk(state, selection, controller, 0)
}

func walk(state *WalkState, selection Selection, controller WalkController, level int) (err error) {
	if schema.IsList(state.SelectedMeta()) && !state.InsideList() {
		var hasMore bool
		if hasMore, err = controller.ListIterator(state, selection, level, true); err != nil {
			return
		}
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			state.SetInsideList()

			if err = walk(state, selection, controller, level); err != nil {
				return
			}
			if hasMore, err = controller.ListIterator(state, selection, level, false); err != nil {
				return
			}
		}
	} else {
		var child Selection
		i := controller.ContainerIterator(state, selection, level)
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
				if err = walk(state.Select(), child, controller, level + 1); err != nil {
					return
				}

				err = selection.Unselect(state, metaList)
			}
		}
	}
	return
}

