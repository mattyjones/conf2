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

var operationNames = []string {
	"N/A",
	"CREATE_CHILD",
	"POST_CREATE_CHILD",
	"CREATE_LIST",
	"POST_CREATE_LIST",
	"UPDATE_VALUE",
	"DELETE_CHILD",
	"DELETE_LIST",
	"BEGIN_EDIT",
	"END_EDIT",
	"CREATE_LIST_ITEM",
	"POST_CREATE_LIST_ITEM",
}

func (op Operation) String() string {
	return operationNames[op]
}

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

func Action(impl Selection, rdr Selection, wtr Selection) (err error) {
	meta := impl.WalkState().Position.(*schema.Rpc)
	var input, output Selection
	if input, output, err = impl.Action(meta); err != nil {
		return err
	}
	if err = Insert(rdr, input, NewExhaustiveController()); err != nil {
		return err
	}
	if err = Insert(output, wtr, NewExhaustiveController()); err != nil {
		return err
	}
	return
}

func edit(from Selection, dest Selection, strategy strategy, controller WalkController) (err error) {
	e := editor{}
	var s Selection
	s, err = e.editTarget(from, dest, false, dest.WalkState().InsideList, strategy)
	if err == nil {
		// sync dest and from
		s.WalkState().Meta = from.WalkState().Meta
		dest.WalkState().Meta = s.WalkState().Meta
		if err = dest.Write(dest.WalkState().Meta, BEGIN_EDIT, nil); err == nil {
			if err = WalkExhaustive(s, controller); err == nil {
				err = dest.Write(dest.WalkState().Meta, END_EDIT, nil)
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
		metaList := s.State.Position.(schema.MetaList)
		if sube = to.Write(metaList, CREATE_LIST, nil); sube != nil {
			return
		}
		if toChild, sube = to.Select(metaList); sube != nil {
			return
		}
		createdList = true
		return
	}
	createContainer := func(to Selection) (toChild Selection, sube error) {
		metaList := s.State.Position.(schema.MetaList)
		if sube = to.Write(metaList, CREATE_CHILD, nil); sube != nil {
			return
		}
		if toChild, sube = to.Select(metaList); sube != nil {
			return
		}
		createdContainer = true
		return
	}
	s.OnSelect = func(meta schema.MetaList) (c Selection, err error) {
		fromState := from.WalkState()
		fromState.Meta = s.State.Meta
		fromState.Position = s.State.Position
		var fromChild Selection
		metaList := s.State.Position.(schema.MetaList)
		fromChild, err = from.Select(metaList)
		if err != nil {
			return
		} else if fromChild == nil {
			return
		}

		if fromChild.WalkState().Meta == nil {
			fromChild.WalkState().Meta = s.State.Position.(schema.MetaList)
		}

		var toChild Selection
		toState := to.WalkState()
		toState.Meta = s.State.Meta
		toState.Position = fromState.Position
		toChild, err = to.Select(metaList)
		if err != nil {
			return
		}
		isList := schema.IsList(s.State.Position)

		switch strategy {
		case INSERT:
			if toChild != nil {
				msg := fmt.Sprintf("Found existing container %s", s.ToString())
				return nil, &browseError{Code:DUPLICATE_FOUND, Msg:msg}
			}
			if isList {
				toChild, err = createList(to)
			} else {
				toChild, err = createContainer(to)
			}
		case UPSERT:
			if toChild == nil {
				if isList {
					toChild, err = createList(to)
				} else {
					toChild, err = createContainer(to)
				}
			}
		case UPDATE:
			if toChild == nil {
				msg := fmt.Sprintf("Container not found in list %s", s.ToString())
				return nil, &browseError{Code:NOT_FOUND, Msg:msg}
			}
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}

		if err != nil {
			return nil, err
		}
		if toChild == nil {
			return nil, &browseError{Msg:fmt.Sprint("Unable to create edit destination item for ", s.ToString())}
		}
		return e.editTarget(fromChild, toChild, createdList, false, UPSERT)
	}

	s.OnUnselect = func(meta schema.MetaList) (err error) {
		metaList := s.State.Position.(schema.MetaList)
		if createdContainer {
			if err = to.Write(metaList, POST_CREATE_CHILD, nil); err != nil {
				return
			}
			createdContainer = false
		}
		if createdList {
			if err = to.Write(metaList, POST_CREATE_LIST, nil); err != nil {
				return
			}
			createdList = false
		}
		if err = from.Unselect(metaList); err != nil  {
			return
		}
		if err = to.Unselect(metaList); err != nil {
			return
		}
		return
	}
	s.OnRead = func(meta schema.HasDataType) (v *Value, err error) {
		from.WalkState().Position = meta
		to.WalkState().Position = meta
		if v, err = from.Read(meta); err != nil {
			return
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = to.Write(meta, UPDATE_VALUE, v); err != nil {
				return
			}
		}
		return
	}
	createListItem := func() (sube error) {
		sube = to.Write(s.State.Position, CREATE_LIST_ITEM, nil)
		createdListItem = true
		return
	}

	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(key []*Value, first bool) (hasMore bool, err error) {
		from.WalkState().Meta = s.State.Meta
		to.WalkState().Meta = s.State.Meta
		if createdListItem {
			err = to.Write(s.State.Meta, POST_CREATE_LIST_ITEM, nil)
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

		var toKey []*Value
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

func (e *editor) loadKey(explictKey []*Value, whereToFindKey Selection) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return ReadKeys(whereToFindKey)
}