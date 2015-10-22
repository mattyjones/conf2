package browse

import (
	"schema"
	"fmt"
	"errors"
	"strings"
)

type StoreBrowser struct {
	schema schema.MetaList
	store Store
}

func NewStoreBrowser(schema schema.MetaList, store Store) *StoreBrowser {
	return &StoreBrowser{
		schema : schema,
		store : store,
	}
}

func (kv *StoreBrowser) Schema() schema.MetaList {
	return kv.schema
}

func (kv *StoreBrowser) Selector(path *Path, strategy Strategy) (s Selection, state *WalkState, err error) {
	selector := &storeSelector{store:kv.store, strategy : strategy}
	switch strategy {
	case READ, UPDATE, UPSERT:
		err = kv.store.Load()
		if err != nil {
			return nil, nil, err
		}
	}
	state = NewWalkState(kv.schema)
	if strategy == READ {
		// we walk to the destination and legitimately return nil if nothing is there
		s, err = selector.selectContainer("")
		return WalkPath(state, s, path)
	}

	// here we fast-forward to the destination prepared to insert into the parse hierarchy
	state, err = fastForwardState(state, path)
	if err != nil {
		return nil, nil, err
	}
	if schema.IsList(state.SelectedMeta()) && !state.InsideList() {
		s, _ = selector.selectList(path.URL)
	} else {
		s, _ = selector.selectContainer(path.URL)
	}
	return
}

func fastForwardState(initialState *WalkState, path *Path) (state *WalkState, err error) {
	state = initialState
	for _, seg := range path.Segments {
		position := schema.FindByIdentExpandChoices(state.SelectedMeta(), seg.Ident)
		if position == nil {
			return nil, errors.New(fmt.Sprintf("%s.%s not found in schema", state.SelectedMeta().GetIdent(), seg.Ident))
		}
		state.SetPosition(position)
		state = state.Select()
		if len(seg.Keys) > 0 {
			var key []*Value
			if key, err = CoerseKeys(position.(*schema.List), seg.Keys); err != nil {
				return nil, err
			}
			state = state.SelectListItem(key)
		}
	}
	return state, nil
}

type storeSelector struct {
	store Store
	strategy Strategy
}

func (kvs *storeSelector) selectList(parentPath string) (Selection, error) {
	s := &MySelection{}
	var keyList []string
	var i int
	var created Selection
	s.OnNext = func(state *WalkState, meta *schema.List, key []*Value, first bool) (next Selection, err error) {
		if created != nil {
			return created, nil
		}
		if len(key) > 0 {
			if first {
				path :=	kvs.listPath(parentPath, key)
				if hasMore := kvs.store.HasValues(path); hasMore {
					return kvs.selectContainer(path)
				}
			} else {
				return nil, nil
			}
		} else {
			if first {
				keyList, err = kvs.store.KeyList(parentPath, meta)
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
				path :=	kvs.listPath(parentPath, key)
				return kvs.selectContainer(path)
			}
		}
		return
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case END_EDIT:
			kvs.store.Save()
		case CREATE_LIST_ITEM:
			childPath := kvs.listPath(parentPath, state.Key())
			created, err = kvs.selectContainer(childPath)
		case POST_CREATE_LIST_ITEM:
			created = nil
		}
		return
	}
	s.OnAction, _ = kvs.store.Action(parentPath)
	return s, nil
}

func (kvs *storeSelector) containerPath(parentPath string, meta schema.Meta) string {
	if len(parentPath) == 0 {
		return meta.GetIdent()
	}
	return fmt.Sprint(parentPath, "/", meta.GetIdent())
}

func (kvs *storeSelector) listPath(parentPath string, key []*Value) string {
	// TODO: support compound keys
	return fmt.Sprint(parentPath, "=", key[0].String())
}

func (kvs *storeSelector) listPathWithNewKey(parentPath string, key []*Value) string {
	eq := strings.LastIndex(parentPath, "=")
	return kvs.listPath(parentPath[:eq], key)
}

func (kvs *storeSelector) selectContainer(parentPath string) (Selection, error) {
	s := &MySelection{}
	//path := storePath{parent:parentPath}
	var created Selection
	s.OnChoose = func(state *WalkState, choice *schema.Choice) (m schema.Meta, err error) {
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
				candidatePath := kvs.containerPath(parentPath, prop)
				found := kvs.store.HasValues(candidatePath)
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
	s.OnRead = func (state *WalkState, meta schema.HasDataType) (*Value, error) {
		return kvs.store.Value(kvs.containerPath(parentPath, meta), meta.GetDataType())
	}
	s.OnSelect = func(state *WalkState, meta schema.MetaList) (child Selection, err error) {
		if (created != nil) {
			child = created
		} else {
			childPath := kvs.containerPath(parentPath, meta)
			if kvs.store.HasValues(childPath) {
				if schema.IsList(meta) {
					child, err = kvs.selectList(childPath)
				} else {
					child, err = kvs.selectContainer(childPath)
				}
			}
		}
		return
	}
	s.OnWrite = func(state *WalkState, meta schema.Meta, op Operation, v *Value) (err error) {
		switch op {
		case END_EDIT:
			kvs.store.Save()
		case CREATE_LIST:
			childPath := kvs.containerPath(parentPath, meta)
			created, err = kvs.selectList(childPath)
		case CREATE_CHILD:
			childPath := kvs.containerPath(parentPath, meta)
			created, err = kvs.selectContainer(childPath)
		case POST_CREATE_LIST, POST_CREATE_CHILD:
			created = nil
		case UPDATE_VALUE:
			propPath := kvs.containerPath(parentPath, meta)
			if err = kvs.store.SetValue(propPath, v); err == nil {
				if schema.IsKeyLeaf(state.SelectedMeta(), meta) {
					oldPath := parentPath
					// TODO: Support compound keys
					newKey := []*Value{v}
					newPath := kvs.listPathWithNewKey(parentPath, newKey)
					kvs.store.RenameKey(oldPath, newPath)
				}
			}
		}
		return
	}
	s.OnAction, _ = kvs.store.Action(parentPath)
	return s, nil
}

func (kvs *storeSelector) read() (Selection, error) {
	s := &MySelection{}
	return s, nil
}
