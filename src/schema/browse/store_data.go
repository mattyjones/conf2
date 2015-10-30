package browse

import (
	"schema"
	"fmt"
	"errors"
	"strings"
)

type StoreData struct {
	schema schema.MetaList
	store Store
}

func NewStoreData(schema schema.MetaList, store Store) *StoreData {
	return &StoreData{
		schema : schema,
		store : store,
	}
}

func (kv *StoreData) Schema() schema.MetaList {
	return kv.schema
}

func (kv *StoreData) Selector(path *Path) (selection *Selection, err error) {
	if err = kv.store.Load(); err != nil {
		return nil, err
	}
	root := kv.Container("")
	selection = NewSelection(root, kv.schema)
	if len(path.Segments) > 0 {
		return WalkPath(selection, path)
	}
	return
}

//func fastForwardState(initialState *Selection, path *Path) (state *Selection, err error) {
//	state = initialState
//	for _, seg := range path.Segments {
//		position := schema.FindByIdentExpandChoices(state.SelectedMeta(), seg.Ident)
//		if position == nil {
//			return nil, errors.New(fmt.Sprintf("%s.%s not found in schema", state.SelectedMeta().GetIdent(), seg.Ident))
//		}
//		state.SetPosition(position)
//		state = state.Select()
//		if len(seg.Keys) > 0 {
//			var key []*Value
//			if key, err = CoerseKeys(position.(*schema.List), seg.Keys); err != nil {
//				return nil, err
//			}
//			state = state.SelectListItem(key)
//		}
//	}
//	return state, nil
//}
//
//type storeSelector struct {
//	store Store
//	strategy Strategy
//}

func (kv *StoreData) List(parentPath string) (Node) {
	s := &MyNode{}
	var keyList []string
	var i int
	var created Node
	s.OnNext = func(state *Selection, meta *schema.List, key []*Value, first bool) (next Node, err error) {
		if created != nil {
			return created, nil
		}
		if len(key) > 0 {
			if first {
				path :=	kv.listPath(parentPath, key)
				if hasMore := kv.store.HasValues(path); hasMore {
					return kv.Container(path), nil
				}
			} else {
				return nil, nil
			}
		} else {
			if first {
				keyList, err = kv.store.KeyList(parentPath, meta)
				i = 0
			} else {
				i++
			}
			if hasMore := i < len(keyList); hasMore {
				var key []*Value
				if key, err = CoerseKeys(meta, []string{keyList[i]}); err != nil {
					return nil, err
				}
				state.SetKey(key)
				path :=	kv.listPath(parentPath, key)
				return kv.Container(path), nil
			}
		}
		return
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case END_EDIT:
			kv.store.Save()
		case CREATE_LIST_ITEM:
			childPath := kv.listPath(parentPath, state.Key())
			created = kv.Container(childPath)
		case POST_CREATE_LIST_ITEM:
			created = nil
		}
		return
	}
	s.OnAction = func(state *Selection, rpc *schema.Rpc, input *Selection) (output *Selection, err error) {
		path :=	kv.listPath(parentPath, state.Key())
		var action ActionFunc
		if action, err = kv.store.Action(path); err != nil {
			return
		}
		return action(state, rpc, input)
	}
	return s
}

func (kv *StoreData) containerPath(parentPath string, meta schema.Meta) string {
	if len(parentPath) == 0 {
		return meta.GetIdent()
	}
	return fmt.Sprint(parentPath, "/", meta.GetIdent())
}

func (kv *StoreData) listPath(parentPath string, key []*Value) string {
	// TODO: support compound keys
	return fmt.Sprint(parentPath, "=", key[0].String())
}

func (kv *StoreData) listPathWithNewKey(parentPath string, key []*Value) string {
	eq := strings.LastIndex(parentPath, "=")
	return kv.listPath(parentPath[:eq], key)
}

func (kv *StoreData) Container(parentPath string) (Node) {
	s := &MyNode{}
	//path := storePath{parent:parentPath}
	var created Node
	s.OnChoose = func(state *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		// go thru each case and if there are any properties in the data that are not
		// part of the schema, that disqualifies that case and we move onto next case
		// until one case aligns with data.  If no cases align then input in inconclusive
		// i.e. non-discriminating and we should error out.
		cases := schema.NewMetaListIterator(choice, false)
		for cases.HasNextMeta() {
			kase := cases.NextMeta().(*schema.ChoiceCase)
			aligned := true
			props := schema.NewMetaListIterator(kase, true)
			for props.HasNextMeta() {
				prop := props.NextMeta()
				candidatePath := kv.containerPath(parentPath, prop)
				found := kv.store.HasValues(candidatePath)
				if !found {
					aligned = false
					break;
				} else {
					m = prop
				}
			}
			if aligned {
				return m, nil
			}
		}
		msg := fmt.Sprintf("No discriminating data for choice schema %s ", state.String())
		return nil, errors.New(msg)
	}
	s.OnRead = func (state *Selection, meta schema.HasDataType) (*Value, error) {
		return kv.store.Value(kv.containerPath(parentPath, meta), meta.GetDataType())
	}
	s.OnSelect = func(state *Selection, meta schema.MetaList) (child Node, err error) {
		if (created != nil) {
			child = created
		} else {
			childPath := kv.containerPath(parentPath, meta)
			if kv.store.HasValues(childPath) {
				if schema.IsList(meta) {
					child = kv.List(childPath)
				} else {
					child = kv.Container(childPath)
				}
			}
		}
		return
	}
	s.OnWrite = func(state *Selection, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case END_EDIT:
			kv.store.Save()
		case CREATE_LIST:
			childPath := kv.containerPath(parentPath, meta)
			created = kv.List(childPath)
		case CREATE_CONTAINER:
			childPath := kv.containerPath(parentPath, meta)
			created = kv.Container(childPath)
		case POST_CREATE_LIST, POST_CREATE_CONTAINER:
			created = nil
		case UPDATE_VALUE:
			propPath := kv.containerPath(parentPath, meta)
			if err = kv.store.SetValue(propPath, v); err == nil {
				if schema.IsKeyLeaf(state.SelectedMeta(), meta) {
					oldPath := parentPath
					// TODO: Support compound keys
					newKey := []*Value{v}
					newPath := kv.listPathWithNewKey(parentPath, newKey)
					kv.store.RenameKey(oldPath, newPath)
				}
			}
		}
		return
	}
	s.OnAction = func(state *Selection, rpc *schema.Rpc, input *Selection) (output *Selection, err error) {
		path := kv.containerPath(parentPath, rpc)
		var action ActionFunc
		if action, err = kv.store.Action(path); err != nil {
			return
		}
		return action(state, rpc, input)
	}
	return s
}
