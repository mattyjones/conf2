package data

import (
	"errors"
	"fmt"
	"schema"
	"strings"
)

type StoreData struct {
	schema schema.MetaList
	store  Store
}

func NewStoreData(schema schema.MetaList, store Store) *StoreData {
	return &StoreData{
		schema: schema,
		store:  store,
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

func (kv *StoreData) OnEvent(sel *Selection, e Event) error {
	switch e {
	case END_EDIT:
		return kv.store.Save()
	}
	return nil
}

func (kv *StoreData) List(parentPath string) Node {
	s := &MyNode{}
	var keyList []string
	var i int
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		if new {
			childPath := kv.listPath(parentPath, sel.State.Key())
			return kv.Container(childPath), nil
		}
		if len(key) > 0 {
			if first {
				path := kv.listPath(parentPath, key)
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
				sel.State.SetKey(key)
				path := kv.listPath(parentPath, key)
				return kv.Container(path), nil
			}
		}
		return
	}
	s.OnEvent = func(sel *Selection, e Event) error {
		switch e {
		case DELETE:
			return kv.store.RemoveAll(parentPath)
		}
		return kv.OnEvent(sel, e)
	}
	s.OnAction = func(sel *Selection, rpc *schema.Rpc, input *Selection) (output *Selection, err error) {
		path := kv.listPath(parentPath, sel.State.Key())
		var action ActionFunc
		if action, err = kv.store.Action(path); err != nil {
			return
		}
		return action(sel, rpc, input)
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

func (kv *StoreData) Container(copy string) Node {
	s := &MyNode{}
	//path := storePath{parent:parentPath}
	s.OnChoose = func(sel *Selection, choice *schema.Choice) (m schema.Meta, err error) {
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
				candidatePath := kv.containerPath(copy, prop)
				found := kv.store.HasValues(candidatePath)
				if !found {
					aligned = false
					break
				} else {
					m = prop
				}
			}
			if aligned {
				return m, nil
			}
		}
		msg := fmt.Sprintf("No discriminating data for choice schema %s ", sel.String())
		return nil, errors.New(msg)
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (*Value, error) {
		return kv.store.Value(kv.containerPath(copy, meta), meta.GetDataType()), nil
	}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (child Node, err error) {
		if new {
			if schema.IsList(meta) {
				childPath := kv.containerPath(copy, meta)
				return kv.List(childPath), nil
			} else {
				childPath := kv.containerPath(copy, meta)
				return kv.Container(childPath), nil
			}
		}
		childPath := kv.containerPath(copy, meta)
		if kv.store.HasValues(childPath) {
		if schema.IsList(meta) {
				return kv.List(childPath), nil
			} else {
				return kv.Container(childPath), nil
			}
		}
		return
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, v *Value) (err error) {
		propPath := kv.containerPath(copy, meta)
		if err = kv.store.SetValue(propPath, v); err != nil {
			return err
		}
		if schema.IsKeyLeaf(sel.State.SelectedMeta(), meta) {
			oldPath := copy
			// TODO: Support compound keys
			newKey := []*Value{v}
			newPath := kv.listPathWithNewKey(copy, newKey)
			kv.store.RenameKey(oldPath, newPath)
		}
		return
	}
	s.OnEvent = func(sel *Selection, e Event) error {
		switch e {
		case DELETE:
			return kv.store.RemoveAll(copy)
		}
		return kv.OnEvent(sel, e)
    }
	s.OnAction = func(sel *Selection, rpc *schema.Rpc, input *Selection) (output *Selection, err error) {
		path := kv.containerPath(copy, rpc)
		var action ActionFunc
		if action, err = kv.store.Action(path); err != nil {
			return
		}
		return action(sel, rpc, input)
	}
	return s
}
