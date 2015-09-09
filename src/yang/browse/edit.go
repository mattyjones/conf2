package browse
import (
	"yang"
	"fmt"
)

type Operation int
const (
	CREATE_CHILD Operation = 1 + iota // 1
	POST_CREATE_CHILD                 // 2
	CREATE_LIST                       // 3
	POST_CREATE_LIST                  // 4
	UPDATE_VALUE                      // 5
	DELETE_CHILD                      // 6
	DELETE_LIST                       // 7
	BEGIN_EDIT                        // 8
	END_EDIT                          // 9
	CREATE_LIST_ITEM                  // 10
	POST_CREATE_LIST_ITEM             // 11
)

type strategy int
const (
	UPSERT strategy = iota + 1
	INSERT
	UPDATE
	DELETE
	CLEAR
)

type editor struct {
}

func Insert(from *Selection, to *Selection, controller WalkController) error {
	return edit(from, to, INSERT, controller)
}

func Upsert(from *Selection, to *Selection, controller WalkController) error {
	return edit(from, to, UPSERT, controller)
}

func Delete(from *Selection, to *Selection, p *Path, controller WalkController) error {
	return edit(from, to, DELETE, controller)
}

func Update(from *Selection, to *Selection, controller WalkController) error {
	return edit(from, to, UPDATE, controller)
}

func edit(from *Selection, dest *Selection, strategy strategy, controller WalkController) (err error) {
	e := editor{}
	var s *Selection
	s, err = e.editTarget(from, dest, strategy)
	if err == nil {
		s.Meta = from.Meta
		dest.Meta = s.Meta
		if err = dest.Edit(BEGIN_EDIT, nil); err == nil {
			if err = WalkExhaustive(s, controller); err == nil {
				err = dest.Edit(END_EDIT, nil)
			}
		}
	}
	return
}

func (e *editor) editTarget(from *Selection, to *Selection, strategy strategy) (*Selection, error) {
	var createdChild bool
	var createdList bool
	s := &Selection{}
	s.insideList = from.insideList
	s.Choose = func(choice *yang.Choice) (yang.Meta, error) {
		return from.Choose(choice)
	}
	s.Enter = func() (c *Selection, err error) {
		from.Meta = s.Meta
		from.Position = s.Position
		var fromChild *Selection
		fromChild, err = from.Enter()
		if err != nil {
			return
		}
		if !from.Found {
			s.Found = false
			return
		}
		if fromChild.Meta == nil {
			fromChild.Meta = s.Position.(yang.MetaList)
		}

		var toChild *Selection
		to.Meta = s.Meta
		to.Position = from.Position
		toChild, err = to.Enter()
		if err != nil {
			return
		}

		s.Found = from.Found
		if fromChild == nil || (!from.Found && !to.Found) {
			return
		}

		nextStrategy := strategy
//		if from.Found && to.Found {
//			switch strategy {
//			case INSERT:
//				err = &browseError{Msg:"Duplicate object found"}
//			case UPDATE:
//				strategy = CLEAR
//			case DELETE:
//				if yang.IsList(s.Position) {
//					err = to.DeleteList()
//				} else {
//					err = to.DeleteChild()
//				}
//				s.Found = false
//			}
//		} else if from.Found && !to.Found {
			switch strategy {
			case UPSERT, INSERT, CLEAR:
				if yang.IsList(s.Position) {
					err = to.CreateList()
					createdList = true
				} else {
					err = to.CreateChild()
					createdChild = true
				}

				if err == nil {
					toChild, err = to.Enter()
					if err == nil && toChild == nil {
						err = &browseError{Msg:"Could not select object that was just created"}
					}
				}
			case UPDATE, DELETE:
				err = &browseError{Msg:"No such object"}
			}
//		} else if !from.Found && to.Found {
//			switch strategy {
//			case DELETE, CLEAR:
//				if yang.IsList(s.Position) {
//					err = to.DeleteList()
//				} else {
//					err = to.DeleteChild()
//				}
//				s.Found = false
//			}
//		}

		if err == nil && s.Found {
			return e.editTarget(fromChild, toChild, nextStrategy)
		}

		return
	}
	s.Exit = func() (err error) {
		if createdChild {
			if err = to.Edit(POST_CREATE_CHILD, nil); err != nil {
				return
			}
			createdChild = false
		}
		if createdList {
			if err = to.Edit(POST_CREATE_LIST, nil); err != nil {
				return
			}
			createdList = false
		}
		if from.Exit != nil {
			if err = from.Exit(); err != nil {
				return
			}
		}
		if to.Exit != nil {
			if err = to.Exit(); err != nil {
				return
			}
		}
		return
	}
	s.ReadValue = func(v *Value) (err error) {
		from.Position = s.Position
		to.Position = s.Position
		if from.ReadValue == nil {
			msg := fmt.Sprint("Read not implemented on ", from.Meta.GetIdent())
			return &browseError{Msg:msg}
		}
		if err = from.ReadValue(v); err != nil {
			return
		}
		if err = to.SetValue(v); err != nil {
			return
		}

// TODO: support strategies
//		var copy bool
//		var clear bool
//		if from.Found && to.Found {
//			switch strategy {
//			case UPSERT, UPDATE, CLEAR:
//				copy = true
//			}
//		} else if from.Found && !to.Found {
//			switch strategy {
//			case INSERT, UPSERT, CLEAR:
//				copy = true
//			}
//		} else if !from.Found && to.Found {
//			switch strategy {
//			case UPDATE, CLEAR:
//				clear = true
//			}
//		}
//		if copy || clear {
//			if err = to.SetValue(v); err != nil {
//				return
//			}
//		}
		return
	}
	s.Iterate = func(fromKeys []string, first bool) (hasMore bool, err error) {
		from.Meta = s.Meta
		to.Meta = s.Meta
		if from.Iterate == nil {
			msg := fmt.Sprint("Missing destination iterator on ", from.Meta.GetIdent())
			return false, &browseError{Msg:msg}
		}
		hasMore, err = from.Iterate(fromKeys, first)

		if err != nil {
			return
		}

		if hasMore {
			_, err = to.Iterate(fromKeys, first)
			if err != nil {
				return
			}
		}

		// TODO: Consider to.hasMore results on LIST_ITEM calls
		if first && hasMore {
			err = to.Edit(CREATE_LIST_ITEM, nil)
		} else if !first && hasMore {
			err = to.Edit(POST_CREATE_LIST_ITEM, nil)
			if err == nil {
				err = to.Edit(CREATE_LIST_ITEM, nil)
			}
		} else if !first && !hasMore {
			err = to.Edit(POST_CREATE_LIST_ITEM, nil)
		}

		return

// TODO assumption for now
//			toKeys := fromKeys
//			if len(toKeys) == 0 {
//				keyIdents := s.Meta.(*yang.List).Keys
//				toKeys = make([]string, len(keyIdents))
//				for i, keyIdent := range keyIdents {
//					v := &Value{}
//					if _, err = from.Select(keyIdent); err != nil {
//						return
//					}
//					if err = from.Read(v); err != nil {
//						return
//					}
//					// TODO: don't assume key is a string
//					toKeys[i] = v.Str
//				}
//			}
//
//			// ignore if exists or not, next Select will detect existance for lists and container
//			// selections.
//			_, err = to.Iterate(toKeys, true)
//			if err != nil {
//				return
//			}
	}

	return s, nil
}
