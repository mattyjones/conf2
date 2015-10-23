package browse
import (
	"schema"
	"fmt"
	"net/http"
	"errors"
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
		if state == nil {
			return NotFound(path.URL)
		}
	}
	return Edit(state, srcSel, destSel, strategy, LimitedWalk(path.Query))
}


func Action(path *Path, impl Browser, input Browser, output Browser) (err error) {
	var aSel, inputSel, outputSel, rpcOutput Selection
	var aState, outputState *WalkState
	if aSel, aState, err = impl.Selector(path, ACTION); err != nil {
		return err
	}
	if aSel == nil {
		return errors.New(fmt.Sprint("No action found at ", path.URL))
	}
	rpc := aState.SelectedMeta().(*schema.Rpc)
	if inputSel, _, err = input.Selector(NewPath(""), READ); err != nil {
		return err
	}
	if rpcOutput, outputState, err = aSel.Action(aState, rpc, inputSel); err != nil {
		return err
	}
	if rpcOutput != nil {
		outputSel, _, err = output.Selector(NewPath(""), INSERT)
		if err = Edit(outputState, rpcOutput, outputSel, INSERT, LimitedWalk(path.Query)); err != nil {
			return err
		}
	}
	return err
}

func Edit(state *WalkState, from Selection, dest Selection, strategy Strategy, controller WalkController) (err error) {
	e := editor{}
	var s Selection
	if schema.IsList(state.SelectedMeta()) && !state.InsideList() {
		s, err = e.editList(state, from, dest, strategy)
	} else {
		s, err = e.editTarget(state, from, dest, strategy)
	}
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

func (e *editor) editList(state *WalkState, from Selection, to Selection, strategy Strategy) (Selection, error) {
	if to == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get target list selection for ", state.String())}
	}
	if from == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get source list selection for ", state.String())}
	}
	s := &MySelection{}
	var createdListItem bool
	createListItem := func(state *WalkState, meta *schema.List, key []*Value) (next Selection, sube error) {
		sube = to.Write(state, state.SelectedMeta(), CREATE_LIST_ITEM, nil)
		createdListItem = true
		if sube == nil {
			return to.Next(state, meta, key, true)
		}

		return
	}

	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, first bool) (next Selection, err error) {
		if createdListItem {
			err = to.Write(state, meta, POST_CREATE_LIST_ITEM, nil)
			createdListItem = false
			if err != nil {
				return nil, err
			}
		}

		var fromNext Selection
		fromNext, err = from.Next(state, meta, key, first)
		if err != nil || fromNext == nil {
			return nil, err
		}

		//		if !listIterationInitialized {
		//			listIterationInitialized = true
		//			if isNewList {
		//				autoCreateListItems = true
		//			} else {
		//				var isListNotEmpty bool
		//				if isListNotEmpty, err = to.Next(state, meta, NO_KEYS, true); err != nil {
		//					return nil, err
		//				}
		//				if ! isListNotEmpty { // is empty
		//					autoCreateListItems = true
		//				}
		//			}
		//		}
		//		if autoCreateListItems {
		//			fromNext
		//			return true, createListItem(state)
		//		}

		var toKey []*Value
		var toNext Selection
		if toKey, err = e.loadKey(state, key, from); err != nil {
			return nil, err
		}
		if toNext, err = to.Next(state, meta, toKey, true); err != nil {
			return nil, err
		}
		switch strategy {
		case UPDATE:
			if toNext == nil {
				msg := fmt.Sprint("No item found with given key in list ", state.String())
				return nil, &browseError{Code:http.StatusNotFound, Msg:msg}
			}
			return e.editTarget(state, fromNext, toNext, strategy)
		case UPSERT:
			if toNext == nil {
				if toNext, err = createListItem(state, meta, toKey); err != nil {
					return nil, err
				}
			}
			return e.editTarget(state, fromNext, toNext, strategy)
		case INSERT:
			if toNext != nil {
				msg := fmt.Sprint("Duplicate item found with same key in list ", state.String())
				return nil, &browseError{Code:http.StatusConflict, Msg:msg}
			}
			if toNext, err = createListItem(state, meta, toKey); err != nil {
				return nil, err
			}
			return e.editTarget(state, fromNext, toNext, strategy)
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}
		return nil, nil
	}
	return s, nil
}

func (e *editor) editTarget(state *WalkState, from Selection, to Selection, strategy Strategy) (Selection, error) {
	if to == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get target container selection for ", state.String())}
	}
	var createdContainer bool
	var createdList bool
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
				msg := fmt.Sprint("Found existing container ", state.String())
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
				msg := fmt.Sprint("Container not found in list ", state.String())
				return nil, &browseError{Code:http.StatusNotFound, Msg:msg}
			}
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}

		if err != nil {
			return nil, err
		}
		// we always switch to upsert strategy because if there were any conflicts, it would have been
		// discovered in top-most level.
		if isList {
			return e.editList(state, fromChild, toChild, UPSERT)
		}
		return e.editTarget(state, fromChild, toChild, UPSERT)
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

	return s, nil
}

func (e *editor) loadKey(state *WalkState, explictKey []*Value, whereToFindKey Selection) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return ReadKeys(state, whereToFindKey)
}