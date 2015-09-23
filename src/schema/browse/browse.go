package browse

import (
	"schema"
)

type Browser interface {
	RootSelector() (Selection, error)
	Module() (*schema.Module)
}

type WalkController interface {
	ListIterator(s Selection, level int, first bool) (hasMore bool, err error)
	ContainerIterator(s Selection, level int) schema.MetaIterator
	CloseSelection(s Selection) error
}

func WalkPath(from Selection, path *Path) (s Selection, err error) {
	nest := path.FindTargetController()
	err = walk(from, nest, 0)
	return nest.target, err
}

func WalkExhaustive(selection Selection, controller WalkController) (err error) {
	return walk(selection, controller, 0)
}

func walk(selection Selection, controller WalkController, level int) (err error) {
	state := selection.WalkState()
	if schema.IsList(state.Meta) && !state.InsideList {
		var hasMore bool
		if hasMore, err = controller.ListIterator(selection, level, true); err != nil {
			return
		}
		for i := 0; hasMore; i++ {

			// important flag, otherwise we recurse indefinitely
			state.InsideList = true

			if err = walk(selection, controller, level); err != nil {
				return
			}
			if hasMore, err = controller.ListIterator(selection, level, false); err != nil {
				return
			}
		}
	} else {
		var child Selection
		i := controller.ContainerIterator(selection, level)
		for i.HasNextMeta() {
			state.Position = i.NextMeta()
			if choice, isChoice := state.Position.(*schema.Choice); isChoice {
				if state.Position, err = selection.Choose(choice); err != nil {
					return
				}
			}
			if schema.IsLeaf(state.Position) {
				// only walking here, not interested in value
				if _, err = selection.Read(state.Position.(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := state.Position.(schema.MetaList)
				child, err = selection.Select(metaList)
				if err != nil {
					return
				} else if child == nil {
					continue
				}
				child.WalkState().Meta = metaList
				defer schema.CloseResource(child)

				if err = walk(child, controller, level + 1); err != nil {
					return
				}

				err = selection.Unselect(metaList)
			}
		}
	}
	return
}

