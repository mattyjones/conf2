package data

import (
	"schema"
)

func (self *Selection) Walk(controller WalkController) (err error) {
	if schema.IsList(self.path.meta) && !self.insideList {
		var next *Selection
		if next, err = controller.ListIterator(self, true); err != nil {
			return
		}
		for next != nil {
			if err = next.Walk(controller); err != nil {
				return
			}
			if err = next.Fire(LEAVE); err != nil {
				return err
			}
			if next, err = controller.ListIterator(self, false); err != nil {
				return
			}
		}
	} else {
		i, cerr := controller.ContainerIterator(self)
		if cerr != nil {
			return cerr
		}
		for i.HasNextMeta() {
			meta := i.NextMeta()
			if choice, isChoice := meta.(*schema.Choice); isChoice {
				var chosen schema.Meta
				if chosen, err = self.node.Choose(self, choice); err != nil {
					return
				}
				meta = chosen
			}
			if schema.IsLeaf(meta) {
				// only walking here, not interested in value
				if _, err = self.node.Read(self, meta.(schema.HasDataType)); err != nil {
					return err
				}
			} else {
				metaList := meta.(schema.MetaList)
				if schema.IsAction(meta) {
					if _, err = controller.VisitAction(self, metaList.(*schema.Rpc)); err != nil {
						return err
					}
				} else {
					childSel, childErr := controller.VisitContainer(self, metaList)
					if childErr != nil {
						return childErr
					} else if childSel == nil {
						continue
					}

					if err = childSel.Walk(controller); err != nil {
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


