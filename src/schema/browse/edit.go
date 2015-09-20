package browse
import (
	"schema"
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
	s, err = e.editTarget(from, dest, false, dest.WalkState().InsideList, strategy)
	if err == nil {
		// sync dest and from
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

func (e *editor) editTarget(from Selection, to Selection, isNewList bool, isInsideList bool, strategy strategy) (Selection, error) {
	var createdContainer bool
	var createdList bool
	var createdListItem bool
	var autoCreateListItems bool
	var listIterationInitialized bool
	s := &MySelection{}
	s.State.InsideList = from.WalkState().InsideList
	s.OnChoose = func(choice *schema.Choice) (schema.Meta, error) {
		return from.Choose(choice)
	}
	createList := func(to Selection) (toChild Selection, sube error) {
		if sube = to.Write(CREATE_LIST, nil); sube != nil {
			return
		}
		if toChild, sube = to.Select(); sube != nil {
			return
		}
		createdList = true
		return
	}
	createContainer := func(to Selection) (toChild Selection, sube error) {
		if sube = to.Write(CREATE_CHILD, nil); sube != nil {
			return
		}
		if toChild, sube = to.Select(); sube != nil {
			return
		}
		createdContainer = true
		return
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
		isList := schema.IsList(s.State.Position)
		found := to.WalkState().Found

		switch strategy {
		case INSERT:
			if found {
				msg := fmt.Sprintf("Found existing container %s", s.ToString())
				return nil, &browseError{Code:DUPLICATE_FOUND, Msg:msg}
			}
			if isList {
				toChild, err = createList(to)
			} else {
				toChild, err = createContainer(to)
			}
		case UPSERT:
			if !found {
				if isList {
					toChild, err = createList(to)
				} else {
					toChild, err = createContainer(to)
				}
			}
		case UPDATE:
			if !found {
				msg := fmt.Sprintf("Container not found in list %s", s.ToString())
				return nil, &browseError{Code:NOT_FOUND, Msg:msg}
			}
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}

		if err != nil {
			return nil, err
		}
		s.State.Found = true
		return e.editTarget(fromChild, toChild, createdList, false, UPSERT)
	}

	s.OnUnselect = func() (err error) {
		if createdContainer {
			if err = to.Write(POST_CREATE_CHILD, nil); err != nil {
				return
			}
			createdContainer = false
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
		if from.WalkState().Found {
			if err = to.Write(UPDATE_VALUE, v); err != nil {
				return
			}
		}
		return
	}
	createListItem := func() (sube error) {
		sube = to.Write(CREATE_LIST_ITEM, nil)
		createdListItem = true
		return
	}

	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(key []Value, first bool) (hasMore bool, err error) {
		from.WalkState().Meta = s.State.Meta
		to.WalkState().Meta = s.State.Meta
		if createdListItem {
			err = to.Write(POST_CREATE_LIST_ITEM, nil)
			createdListItem = false
			if err != nil {
				return false, err
			}
		}

		hasMore, err = from.Next(key, first)
		if err != nil || ! hasMore {
			return hasMore, err
		}

		if !listIterationInitialized {
			listIterationInitialized = true
			if isNewList {
				autoCreateListItems = true
			} else {
				var isListNotEmpty bool
				if isListNotEmpty, err = to.Next(NO_KEYS, true); err != nil {
					return false, err
				}
				if ! isListNotEmpty { // is empty
					autoCreateListItems = true
				}
			}
		}
		if autoCreateListItems {
			return true, createListItem()
		}

		var toKey []Value
		var foundItem bool
		if toKey, err = e.loadKey(key, from); err != nil {
			return false, err
		}
		if foundItem, err = to.Next(toKey, true); err != nil {
			return false, err
		}
		switch strategy {
		case UPDATE:
			if !foundItem {
				msg := fmt.Sprintf("No item found with given key in list %s", s.ToString())
				return false, &browseError{Code:NOT_FOUND, Msg:msg}
			}
		case UPSERT:
			if !foundItem {
				return true, createListItem()
			}
		case INSERT:
			if foundItem {
				msg := fmt.Sprintf("Duplicate item found with same key in list %s", s.ToString())
				return false, &browseError{Code:DUPLICATE_FOUND, Msg:msg}
			}
			return true, createListItem()
		default:
			return false, &browseError{Msg:"Stratgey not implmented"}
		}
		return true, nil
	}

	return s, nil
}

func (e *editor) loadKey(explictKey []Value, whereToFindKey Selection) ([]Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return ReadKeys(whereToFindKey)
}