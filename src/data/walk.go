package data

import (
	"schema"
)

func WalkPath(selection *Selection, path *schema.PathSlice) (*Selection, error) {
	if path.Empty() {
		return selection, nil
	}
	finder := NewFindTarget(path)
	err := Walk(selection, finder)
	return finder.Target, err
}

func WalkDataPath(data Data, path string) (*Selection, error) {
	p, err := schema.ParsePath(path, data.Schema())
	if err != nil {
		return nil, err
	}
	return WalkPath(NewSelection(data.Node(), data.Schema()), p)
}

func Walk(selection *Selection, controller WalkController) (err error) {
	state := selection.State
	if schema.IsList(state.SelectedMeta()) && !state.InsideList() {
		var next *Selection
		if next, err = controller.ListIterator(selection, true); err != nil {
			return
		}
		for next != nil {
			if err = Walk(next, controller); err != nil {
				return
			}
			if err = next.Fire(LEAVE); err != nil {
				return err
			}
			if next, err = controller.ListIterator(selection, false); err != nil {
				return
			}
		}
	} else {
		var child Node
		i, cerr := controller.ContainerIterator(selection)
		if cerr != nil {
			return cerr
		}
		for i.HasNextMeta() {
			state.SetPosition(i.NextMeta())
			if choice, isChoice := state.Position().(*schema.Choice); isChoice {
				var chosen schema.Meta
				if chosen, err = selection.Node.Choose(selection, choice); err != nil {
					return
				}
				state.SetPosition(chosen)
			}
			if schema.IsLeaf(state.Position()) {
				// only walking here, not interested in value
				if _, err = selection.Node.Read(selection, state.Position().(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := state.Position().(schema.MetaList)
				if schema.IsAction(state.Position()) {
					if err = controller.VisitAction(selection); err != nil {
						return err
					}
				} else {
					if child, err = selection.Node.Select(selection, metaList, false); err != nil {
						return err
					}
					if child == nil {
						continue
					}
					defer schema.CloseResource(child)
					childSel := selection.Select(child)

					if err = Walk(childSel, controller); err != nil {
						return
					}
					if err = childSel.Fire(LEAVE); err != nil {
						return err
					}
				}
			}
		}
	}
	return
}
