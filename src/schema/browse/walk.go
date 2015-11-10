package browse

import (
	"schema"
)

func WalkPath(selection *Selection, path *Path) (*Selection, error) {
	finder := NewFindTarget(path)
	err := walk(selection, finder)
	return finder.Target, err
}

func Walk(selection *Selection, controller WalkController) (err error) {
	return walk(selection, controller)
}

func walk(selection *Selection, controller WalkController) (err error) {
	if schema.IsList(selection.SelectedMeta()) && !selection.InsideList() {
		var next *Selection
		if next, err = controller.ListIterator(selection, true); err != nil {
			return
		}
		for next != nil {
			if err = walk(next, controller); err != nil {
				return
			}
			if next, err = controller.ListIterator(selection, false); err != nil {
				return
			}
		}
	} else {
		var child Node
		i := controller.ContainerIterator(selection)
		for i.HasNextMeta() {
			selection.SetPosition(i.NextMeta())
			if choice, isChoice := selection.Position().(*schema.Choice); isChoice {
				var choosen schema.Meta
				if choosen, err = selection.Node().Choose(selection, choice); err != nil {
					return
				}
				selection.SetPosition(choosen)
			}
			if schema.IsLeaf(selection.Position()) {
				// only walking here, not interested in value
				if _, err = selection.Node().Read(selection, selection.Position().(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := selection.Position().(schema.MetaList)
				if schema.IsAction(selection.Position()) {
					if err = controller.VisitAction(selection); err != nil {
						return err
					}
				} else {
					child, err = selection.Node().Select(selection, metaList)
					if err != nil {
						return
					} else if child == nil {
						continue
					}
					defer schema.CloseResource(child)
					if err = walk(selection.Select(child), controller); err != nil {
						return
					}

					err = selection.Node().Unselect(selection, metaList)
				}
			}
		}
	}
	return
}
