package browse
import (
	"schema"
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

func Insert(from Selection, to Selection, controller WalkController) error {
	return edit(from, to, INSERT, controller)
}

func Upsert(from Selection, to Selection, controller WalkController) error {
	return edit(from, to, UPSERT, controller)
}

func Delete(from Selection, to Selection, p *Path, controller WalkController) error {
	return edit(from, to, DELETE, controller)
}

func Update(from Selection, to Selection, controller WalkController) error {
	return edit(from, to, UPDATE, controller)
}

func edit(from Selection, dest Selection, strategy strategy, controller WalkController) (err error) {
	e := editor{}
	var s Selection
	s, err = e.editTarget(from, dest, strategy)
	if err == nil {
		s.WalkState().Meta = from.WalkState().Meta
		dest.WalkState().Meta = s.WalkState().Meta
		if err = dest.Write(BEGIN_EDIT, nil); err == nil {
			if err = WalkExhaustive(s, controller); err == nil {
				err = dest.Write(END_EDIT, nil)
			}
		}
	}
	return
}

func (e *editor) editTarget(from Selection, to Selection, strategy strategy) (Selection, error) {
	var createdChild bool
	var createdList bool
	s := &MySelection{}
	s.State.insideList = from.WalkState().insideList
	s.OnChoose = func(choice *schema.Choice) (schema.Meta, error) {
		return from.Chooze(choice)
	}
	s.OnSelect = func() (c Selection, err error) {
		fromState := from.WalkState()
		fromState.Meta = s.State.Meta
		fromState.Position = s.State.Position
		var fromChild Selection
		fromChild, err = from.Select()
		if err != nil {
			return
		}
		if !fromState.Found {
			s.State.Found = false
			return
		}
		if fromChild.WalkState().Meta == nil {
			fromChild.WalkState().Meta = s.State.Position.(schema.MetaList)
		}

		var toChild Selection
		toState := to.WalkState()
		toState.Meta = s.State.Meta
		toState.Position = fromState.Position
		toChild, err = to.Select()
		if err != nil {
			return
		}

		s.State.Found = fromState.Found
		if fromChild == nil || (!fromState.Found && !toState.Found) {
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
//				if schema.IsList(s.Position) {
//					err = to.DeleteList()
//				} else {
//					err = to.DeleteChild()
//				}
//				s.Found = false
//			}
//		} else if from.Found && !to.Found {
			switch strategy {
			case UPSERT, INSERT, CLEAR:
				if schema.IsList(s.State.Position) {
					err = to.Write(CREATE_LIST, nil)
					createdList = true
				} else {
					err = to.Write(CREATE_CHILD, nil)
					createdChild = true
				}

				if err == nil {
					toChild, err = to.Select()
					toChild.WalkState().Meta = fromChild.WalkState().Meta
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
//				if schema.IsList(s.Position) {
//					err = to.DeleteList()
//				} else {
//					err = to.DeleteChild()
//				}
//				s.Found = false
//			}
//		}

		if err == nil && s.State.Found {
			return e.editTarget(fromChild, toChild, nextStrategy)
		}

		return
	}
	s.OnUnselect = func() (err error) {
		if createdChild {
			if err = to.Write(POST_CREATE_CHILD, nil); err != nil {
				return
			}
			createdChild = false
		}
		if createdList {
			if err = to.Write(POST_CREATE_LIST, nil); err != nil {
				return
			}
			createdList = false
		}
		if err = from.Unselect(); err != nil  {
			return
		}
		if err = to.Unselect(); err != nil {
			return
		}
		return
	}
	s.OnRead = func(v *Value) (err error) {
		from.WalkState().Position = s.State.Position
		to.WalkState().Position = s.State.Position
		if err = from.Read(v); err != nil {
			return
		}
		if err = to.Write(UPDATE_VALUE, v); err != nil {
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
	s.OnNext = func(fromKeys []interface{}, first bool) (hasMore bool, err error) {
		from.WalkState().Meta = s.State.Meta
		to.WalkState().Meta = s.State.Meta
		hasMore, err = from.Next(fromKeys, first)

		if err != nil {
			return
		}

		if hasMore {
			_, err = to.Next(fromKeys, first)
			if err != nil {
				return
			}
		}

		// TODO: Consider to.hasMore results on LIST_ITEM calls
		if first && hasMore {
			err = to.Write(CREATE_LIST_ITEM, nil)
		} else if !first && hasMore {
			err = to.Write(POST_CREATE_LIST_ITEM, nil)
			if err == nil {
				err = to.Write(CREATE_LIST_ITEM, nil)
			}
		} else if !first && !hasMore {
			err = to.Write(POST_CREATE_LIST_ITEM, nil)
		}

		return

// TODO assumption for now
//			toKeys := fromKeys
//			if len(toKeys) == 0 {
//				keyIdents := s.Meta.(*schema.List).Keys
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
