package browse
import (
	"schema"
	"fmt"
	"net/http"
	"errors"
)

type Operation int
const (
	CREATE_CONTAINER Operation = 1 + iota // 1
	POST_CREATE_CONTAINER             // 2
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


func Insert(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, INSERT, FullWalk())
}

func InsertByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, INSERT)
}

func UpdateByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, UPDATE)
}

func UpsertByNode(selection *Selection, src Node, dest Node) (err error) {
	return EditByNode(selection, src, dest, UPSERT)
}

func EditByNode(selection *Selection, src Node, dest Node, strategy Strategy) (err error) {
	e := editor{}
	var n Node
	if schema.IsList(selection.SelectedMeta()) && !selection.InsideList() {
		n, err = e.list(src, dest, strategy)
	} else {
		n, err = e.container(src, dest, strategy)
	}
	if err == nil {
		s := selection.Copy(n)
		if err = dest.Write(s, s.SelectedMeta(), BEGIN_EDIT, nil); err == nil {
			if err = Walk(s, FullWalk()); err == nil {
				err = dest.Write(s, s.SelectedMeta(), END_EDIT, nil)
			}
		}
	}
	return
}

func ControlledInsert(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, INSERT, cntrl)
}

func Upsert(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, UPSERT, FullWalk())
}

func ControlledUpsert(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, UPSERT, cntrl)
}

func Update(src *Selection, dest *Selection) (err error) {
	return Edit(src, dest, UPDATE, FullWalk())
}

func ControlledUpdate(src *Selection, dest *Selection, cntrl WalkController) (err error) {
	return Edit(src, dest, UPDATE, cntrl)
}

func Delete(sel *Selection) error {
	return sel.Node().Write(sel, sel.SelectedMeta(), DELETE_CHILD, nil)
}

func Action(impl *Selection, input Node) (output *Selection, err error) {
	rpc := impl.Position().(*schema.Rpc)
	return impl.Node().Action(impl, rpc, input)
}

func Edit(from *Selection, dest *Selection, strategy Strategy, controller WalkController) (err error) {
	e := editor{}
	var n Node
	if schema.IsList(from.SelectedMeta()) && !from.InsideList() {
		n, err = e.list(from.Node(), dest.Node(), strategy)
	} else {
		n, err = e.container(from.Node(), dest.Node(), strategy)
	}
	if err == nil {
		// sync dest and from
		if err = dest.Node().Write(dest, dest.SelectedMeta(), BEGIN_EDIT, nil); err == nil {
			s := from.Copy(n)
			if err = Walk(s, controller); err == nil {
				err = dest.Node().Write(dest, dest.SelectedMeta(), END_EDIT, nil)
			}
		}
	}
	return
}

func (e *editor) list(from Node, to Node, strategy Strategy) (Node, error) {
	if to == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get target node")}
	}
	if from == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get source node")}
	}
	s := &MyNode{Label:fmt.Sprint("Edit list ", from.String(), "=>", to.String())}
	var createdListItem bool
	createListItem := func(selection *Selection, meta *schema.List, key []*Value) (next Node, err error) {
		err = to.Write(selection, meta, CREATE_LIST_ITEM, nil)
		createdListItem = true
		if err == nil {
			next, err = to.Next(selection, meta, key, true)
			if next == nil {
				return nil, errors.New("Could not create list item")
			}
			return
		}

		return
	}
	// List Edit - See "List Edit State Machine" diagram for additional documentation
	s.OnNext = func(selection *Selection, meta *schema.List, key []*Value, first bool) (next Node, err error) {
		if createdListItem {
			err = to.Write(selection, meta, POST_CREATE_LIST_ITEM, nil)
			createdListItem = false
			if err != nil {
				return nil, err
			}
		}

		var fromNextNode Node
		fromNextNode, err = from.Next(selection, meta, key, first)
		if err != nil || fromNextNode == nil {
			return nil, err
		}

		var nextKey []*Value
		var toNextNode Node
		if nextKey, err = e.loadKey(selection, key); err != nil {
			return nil, err
		}
		if len(nextKey) > 0 {
			if toNextNode, err = to.Next(selection, meta, nextKey, true); err != nil {
				return nil, err
			}
		}
		switch strategy {
		case UPDATE:
			if toNextNode == nil {
				msg := fmt.Sprint("No item found with given key in list ", selection.String())
				return nil, &browseError{Code:http.StatusNotFound, Msg:msg}
			}
		case UPSERT:
			if toNextNode == nil {
				if toNextNode, err = createListItem(selection, meta, nextKey); err != nil {
					return nil, err
				}
			}
		case INSERT:
			if toNextNode != nil {
				msg := fmt.Sprint("Duplicate item found with same key in list ", selection.String())
				return nil, &browseError{Code:http.StatusConflict, Msg:msg}
			}
			if toNextNode, err = createListItem(selection, meta, nextKey); err != nil {
				return nil, err
			}
		default:
			return nil, &browseError{Msg:"Stratgey not implmented"}
		}
		return e.container(fromNextNode, toNextNode, UPSERT)
	}
	return s, nil
}

func (e *editor) container(from Node, to Node, strategy Strategy) (Node, error) {
	if to == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get target container selection")}
	}
	if from == nil {
		return nil, &browseError{Msg:fmt.Sprint("Unable to get source node")}
	}
	var createdContainer bool
	var createdList bool
	s := &MyNode{Label:fmt.Sprint("Edit container ", from.String(), "=>", to.String())}
	s.OnChoose = func(state *Selection, choice *schema.Choice) (schema.Meta, error) {
		return from.Choose(state, choice)
	}
	createList := func(selection *Selection) (toChild Node, err error) {
		if err = to.Write(selection, selection.Position(), CREATE_LIST, nil); err != nil {
			return nil, err
		}
		metaList := selection.Position().(schema.MetaList)
		var listNode Node
		if listNode, err = to.Select(selection, metaList); err != nil {
			return nil, err
		}
		if listNode == nil {
			msg := fmt.Sprint("Failure selecting newly created list ", selection.String())
			return nil, &browseError{Code:http.StatusNotFound, Msg:msg}

		}
		createdList = true
		return listNode, nil
	}
	createContainer := func(selection *Selection) (toChild Node, err error) {
		if err = to.Write(selection, selection.Position(), CREATE_CONTAINER, nil); err != nil {
			return
		}
		metaList := selection.Position().(schema.MetaList)
		var containerNode Node
		if containerNode, err = to.Select(selection, metaList); err != nil {
			return
		}
		if containerNode == nil {
			msg := fmt.Sprint("Failure selection newly created container ", selection.String())
			return nil, &browseError{Code:http.StatusNotFound, Msg:msg}

		}
		createdContainer = true
		return containerNode, nil
	}
	s.OnSelect = func(selection *Selection, meta schema.MetaList) (c Node, err error) {
		var fromChild Node
		fromChild, err = from.Select(selection, meta)
		if err != nil {
			return
		} else if fromChild == nil {
			return
		}

		var toChild Node
		toChild, err = to.Select(selection, meta)
		if err != nil {
			return
		}
		isList := schema.IsList(meta)

		switch strategy {
		case INSERT:
			if toChild != nil {
				msg := fmt.Sprint("Found existing container ", selection.String())
				return nil, &browseError{Code:http.StatusConflict, Msg:msg}
			}
			if isList {
				toChild, err = createList(selection)
			} else {
				toChild, err = createContainer(selection)
			}
		case UPSERT:
			if toChild == nil {
				if isList {
					toChild, err = createList(selection)
				} else {
					toChild, err = createContainer(selection)
				}
			}
		case UPDATE:
			if toChild == nil {
				msg := fmt.Sprint("Container not found in list ", selection.String())
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
			return e.list(fromChild, toChild, UPSERT)
		}
		return e.container(fromChild, toChild, UPSERT)
	}
	s.OnUnselect = func(selection *Selection, meta schema.MetaList) (err error) {
		if createdContainer {
			if err = to.Write(selection, meta, POST_CREATE_CONTAINER, nil); err != nil {
				return
			}
			createdContainer = false
		}
		if createdList {
			if err = to.Write(selection, meta, POST_CREATE_LIST, nil); err != nil {
				return
			}
			createdList = false
		}
		if err = from.Unselect(selection, meta); err != nil  {
			return
		}
		if err = to.Unselect(selection, meta); err != nil {
			return
		}
		return
	}
	s.OnRead = func(selection *Selection, meta schema.HasDataType) (v *Value, err error) {
		if v, err = from.Read(selection, meta); err != nil {
			return
		}
		if v != nil {
			v.Type = meta.GetDataType()
			if err = to.Write(selection, meta, UPDATE_VALUE, v); err != nil {
				return
			}
		}
		return
	}

	return s, nil
}

func (e *editor) loadKey(selection *Selection, explictKey []*Value) ([]*Value, error) {
	if len(explictKey) > 0 {
		return explictKey, nil
	}
	return selection.Key(), nil
}