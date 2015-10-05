package browse
import (
	"schema"
	"fmt"
	"net/http"
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

type Strategy int
const (
	UPSERT Strategy = iota + 1
	INSERT
	UPDATE
	DELETE
	CLEAR
	READ
	ACTION
)

type editor struct {}

func Insert(path *Path, src Browser, dest Browser) (err error) {
	return modifyingOperationWithInput(path, src, dest, INSERT)
}

func Upsert(path *Path, src Browser, dest Browser) (err error) {
	return modifyingOperationWithInput(path, src, dest, UPSERT)
}

func Update(path *Path, src Browser, dest Browser) (err error) {
	return modifyingOperationWithInput(path, src, dest, UPDATE)
}

func Delete(path *Path, b Browser) error {
	sel, state, err := b.Selector(path, DELETE)
	if err != nil {
		return err
	}
	return sel.Write(state, state.SelectedMeta(), DELETE_CHILD, nil)
}

func modifyingOperationWithInput(path *Path, src Browser, dest Browser, strategy Strategy) (err error) {
	var destSel, srcSel Selection
	var destState, srcState *WalkState
	if destSel, destState, err = dest.Selector(path, strategy); err != nil {
		return err
	}
	if destSel == nil {
		return NotFound(path.URL)
	}
	if srcSel, srcState, err = src.Selector(path, READ); err != nil {
		return err
	}

	state := destState
	if state == nil {
		state = srcState
	}
fmt.Printf("edit - state=%s, srcSel=%v\n", state.String(), srcSel)
	return edit(state, srcSel, destSel, strategy, NewFullWalk(path.query))
}

func Action(path *Path, impl Browser, input Browser, output Browser) (err error) {
	var aSel, inputSel, outputSel, rpcInput, rpcOutput Selection
	var aState *WalkState
	if aSel, aState, err = impl.Selector(path, ACTION); err != nil {
		return err
	}
	rpc := aState.Position().(*schema.Rpc)
	if rpcInput, rpcOutput, err = aSel.Action(aState, rpc); err != nil {
		return err
	}

	inputSel, _, err = input.Selector(path, READ)
	if err = edit(aState, inputSel, rpcInput, INSERT, NewFullWalk(path.query)); err != nil {
		return err
	}
	outputSel, _, err = input.Selector(path, INSERT)
	if err = edit(aState, rpcOutput, outputSel, INSERT, NewFullWalk(path.query)); err != nil {
		return err
	}
	return nil
}

//func SelectionInsert(state *WalkState, from Selection, to Selection, controller WalkController) error {
//	return edit(state, from, to, INSERT, controller)
//}
//
//func Upsert(state *WalkState, from Selection, to Selection, controller WalkController) error {
//	return edit(state, from, to, UPSERT, controller)
//}
//
//func Delete(state *WalkState, from Selection, to Selection, p *Path, controller WalkController) error {
//	return edit(state, from, to, DELETE, controller)
//}
//
//func Update(state *WalkState, from Selection, to Selection, controller WalkController) error {
//	return edit(state, from, to, UPDATE, controller)
//}

func edit(state *WalkState, from Selection, dest Selection, strategy Strategy, controller WalkController) (err error) {
	e := editor{}
	var s Selection
	s, err = e.editTarget(from, dest, false, state.InsideList(), strategy)
	if err == nil {
		// sync dest and from
		if err = dest.Write(state, state.SelectedMeta(), BEGIN_EDIT, nil); err == nil {
			if err = Walk(state, s, controller); err == nil {
				err = dest.Write(state, state.SelectedMeta(), END_EDIT, nil)
			}
		}
	}
	return
}

func (e *editor) editTarget(from Selection, to Selection, isNewList bool, isInsideList bool, strategy Strategy) (Selection, error) {
	var createdContainer bool
	var createdList bool
	var createdListItem bool
	var autoCreateListItems bool
	var listIterationInitialized bool
	s := &MySelection{}
	s.OnChoose = func(state *WalkState, choice *schema.Choice) (schema.Meta, error) {
		return from.Choose(state, choice)
	}
	createList := func(state *WalkState, to Selection) (toChild Selection, sube error) {
		if sube = to.Write(state, state.Position(), CREATE_LIST, nil); sube != nil {
			return
		}
		metaList := state.Position().(schema.MetaList)
		if toChild, sube = to.Select(state, metaList); sube != nil {
			return
		}
		createdList = true
		return
	}
	createContainer := func(state *WalkState, to Selection) (toChild Selection, sube error) {
		if sube = to.Write(state, state.Position(), CREATE_CHILD, nil); sube != nil {
			return
		}
		metaList := state.Position().(schema.MetaList)
		if toChild, sube = to.Select(state, metaList); sube != nil {
			return
		}
		createdContainer = true
		return
	}
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (c Selection, err error) {
fmt.Printf("edit - OnSelect\n")
		var fromChild Selection
		fromChild, err = from.Select(state, meta)
		if err != nil {
			return
		} else if fromChild == nil {
			return
		}

		var toChild Selection
		toChild, err = to.Select(state, meta)
		if err != nil {
			return
		}
		isList := schema.IsList(meta)

		switch strategy {
		case INSERT:
			if toChild != nil {
				msg := fmt.Sprintf("Found existing container %s", state.String())
				return nil, &browseError{Code:http.StatusConflict, Msg:msg}
			}
			if isList {
				toChild, err = createList(state, to)
			} else {
				toChild, err = createContainer(state, to)
			}
		case UPSERT:
			if toChild == nil {
				if isList {
					toChild, err = createList(state, to)
				} else {
					toChild, err = createContainer(state, to)
				}
			}
		case UPDATE:
			if toChild == nil {
				msg := fmt.Sprintf("Container not found in list %s", state.String())
				return nil, &browseError{Code:http.StatusNotFound, Msg:msg}
			}
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}

		if err != nil {
			return nil, err
		}
		if toChild == nil {
			return nil, &browseError{Msg:fmt.Sprint("Unable to create edit destination item for ", state.String())}
		}
		return e.editTarget(fromChild, toChild, createdList, false, UPSERT)
	}

	s.OnUnselect = func(state *WalkState, meta schema.MetaList) (err error) {
		if createdContainer {
			if err = to.Write(state, meta, POST_CREATE_CHILD, nil); err != nil {
				return
			}
			createdContainer = false
		}
		if createdList {
			if err = to.Write(state, meta, POST_CREATE_LIST, nil); err != nil {
				return
			}
			createdList = false
		}
		if err = from.Unselect(state, meta); err != nil  {
			return
		}
		if err = to.Unselect(state, meta); err != nil {
			return
		}
		return
	}
	s.OnRead = func(state *WalkState, meta schema.HasDataType) (v *Value, err error) {
		if v, err = from.Read(state, meta); err != nil {
			return
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = to.Write(state, meta, UPDATE_VALUE, v); err != nil {
				return
			}
		}
		return
	}
	createListItem := func(state *WalkState) (sube error) {
		sube = to.Write(state, state.SelectedMeta(), CREATE_LIST_ITEM, nil)
		createdListItem = true
		return
	}

	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, first bool) (hasMore bool, err error) {
fmt.Printf("edit - OnNext\n")
		if createdListItem {
			err = to.Write(state, meta, POST_CREATE_LIST_ITEM, nil)
			createdListItem = false
			if err != nil {
				return false, err
			}
		}

		hasMore, err = from.Next(state, meta, key, first)
		if err != nil || ! hasMore {
			return hasMore, err
		}

		if !listIterationInitialized {
			listIterationInitialized = true
			if isNewList {
				autoCreateListItems = true
			} else {
				var isListNotEmpty bool
				if isListNotEmpty, err = to.Next(state, meta, NO_KEYS, true); err != nil {
					return false, err
				}
				if ! isListNotEmpty { // is empty
					autoCreateListItems = true
				}
			}
		}
		if autoCreateListItems {
			return true, createListItem(state)
		}

		var toKey []*Value
		var foundItem bool
		if toKey, err = e.loadKey(state, key, from); err != nil {
			return false, err
		}
		if foundItem, err = to.Next(state, meta, toKey, true); err != nil {
			return false, err
		}
		switch strategy {
		case UPDATE:
			if !foundItem {
				msg := fmt.Sprintf("No item found with given key in list %s", state.String())
				return false, &browseError{Code:http.StatusNotFound, Msg:msg}
			}
		case UPSERT:
			if !foundItem {
				return true, createListItem(state)
			}
		case INSERT:
			if foundItem {
				msg := fmt.Sprintf("Duplicate item found with same key in list %s", state.String())
				return false, &browseError{Code:http.StatusConflict, Msg:msg}
			}
			return true, createListItem(state)
		default:
			return false, &browseError{Msg:"Stratgey not implmented"}
		}
		return true, nil
	}

	return s, nil
}

func (e *editor) loadKey(state *WalkState, explictKey []*Value, whereToFindKey Selection) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	if len(state.Key()) > 0 {
		return state.Key(), nil
	}
	return ReadKeys(state, whereToFindKey)
}