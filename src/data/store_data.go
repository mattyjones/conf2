package data

import (
	"errors"
	"fmt"
	"schema"
	"strings"
)

type StoreData struct {
	Meta  schema.MetaList
	Store Store
}

func NewStoreData(schema schema.MetaList, store Store) *StoreData {
	return &StoreData{
		Meta: schema,
		Store:  store,
	}
}

func (kv *StoreData) Select() *Selection {
	return NewSelection(kv.Meta, kv.Node())
}

func (kv *StoreData) Node() (Node) {
	var err error
	if err = kv.Store.Load(); err != nil {
		return ErrorNode{Err:err}
	}
	return kv.Container("")
}

func (kv *StoreData) OnEvent(sel *Selection, e Event) error {
	switch e {
	case END_EDIT:
		return kv.Store.Save()
	}
	return nil
}

func (kv *StoreData) List(parentPath string) Node {
	s := &MyNode{Label:"StoreData List"}
	var keyList []string
	var i int
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (next Node, err error) {
		if new {
			var childPath string
			if len(sel.path.key) > 0 {
				childPath = kv.listPath(parentPath, sel.path.key)
			} else {
				childPath = parentPath + "=unknown"
			}
			return kv.Container(childPath), nil
		}
		if len(key) > 0 {
			if first {
				path := kv.listPath(parentPath, key)
				if hasMore := kv.Store.HasValues(path); hasMore {
					return kv.Container(path), nil
				}
			} else {
				return nil, nil
			}
		} else {
			if first {
				keyList, err = kv.Store.KeyList(parentPath, meta)
				i = 0
			} else {
				i++
			}
			if hasMore := i < len(keyList); hasMore {
				var key []*Value
				if key, err = CoerseKeys(meta, []string{keyList[i]}); err != nil {
					return nil, err
				}
				sel.path.key = key
				path := kv.listPath(parentPath, key)
				return kv.Container(path), nil
			}
		}
		return
	}
	s.OnEvent = func(sel *Selection, e Event) error {
		switch e {
		case DELETE:
			return kv.Store.RemoveAll(parentPath)
		}
		return kv.OnEvent(sel, e)
	}
	s.OnAction = func(sel *Selection, rpc *schema.Rpc, input *Selection) (output Node, err error) {
		path := kv.listPath(parentPath, sel.path.key)
		var action ActionFunc
		if action, err = kv.Store.Action(path); err != nil {
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
	s := &MyNode{Label:"StoreData Container"}
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
				found := kv.Store.HasValues(candidatePath)
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
		return kv.Store.Value(kv.containerPath(copy, meta), meta.GetDataType()), nil
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
		if kv.Store.HasValues(childPath) {
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
		if err = kv.Store.SetValue(propPath, v); err != nil {
			return err
		}
		if schema.IsKeyLeaf(sel.path.meta, meta) {
			oldPath := copy
			// TODO: Support compound keys
			newKey := []*Value{v}
			newPath := kv.listPathWithNewKey(copy, newKey)
			kv.Store.RenameKey(oldPath, newPath)
		}
		return
	}
	s.OnEvent = func(sel *Selection, e Event) error {
		switch e {
		case DELETE:
			return kv.Store.RemoveAll(copy)
		}
		return kv.OnEvent(sel, e)
    }
	s.OnAction = func(sel *Selection, rpc *schema.Rpc, input *Selection) (output Node, err error) {
		path := kv.containerPath(copy, rpc)
		var action ActionFunc
		if action, err = kv.Store.Action(path); err != nil {
			return
		}
		return action(sel, rpc, input)
	}
	return s
}
