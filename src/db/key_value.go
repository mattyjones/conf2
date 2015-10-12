package db
import (
	"schema/browse"
	"schema"
	"fmt"
	"strings"
)

type KeyValues struct {
	module *schema.Module
	store Store
}

type KeyValueStore map[string]*browse.Value

func (kvs KeyValueStore) Load() error {
	return nil
}

func (kvs KeyValueStore) Clear() error {
	for k, _ := range kvs {
		delete(kvs, k)
	}
	return nil
}

func (kvs KeyValueStore) HasValues(path string) bool {
	for k, _ := range kvs {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (kvs KeyValueStore) Save() error {
	return nil
}

func (kvs KeyValueStore) KeyList(key string) ([]string, error) {
	// same as mongo store...
}

type KeyList struct {
	set map[string]struct{}

}

func (kl *KeyList)

func (kvs KeyValueStore) Value(key string, dataType *schema.DataType) (*browse.Value, error) {
	if v, found := kvs[key]; found {
		v.Type = dataType
		return v, nil
	}
	return nil, nil
}

func (kvs KeyValueStore) SetValue(key string, v *browse.Value) error {
	kvs[key] = v
	return nil
}

func NewKeyValues(module *schema.Module, store Store) *KeyValues {
	return &KeyValues{
		module : module,
		store : store,
	}
}

type KeyValuesSelector struct {
	store Store
	path *browse.Path
	strategy browse.Strategy
}

func (kv *KeyValues) Module() *schema.Module {
	return kv.module
}

func (kv *KeyValues) Selector(path *browse.Path, strategy browse.Strategy) (s browse.Selection, state *browse.WalkState, err error) {
	selector := &KeyValuesSelector{path: path, store:kv.store, strategy : strategy}
	switch strategy {
	case browse.READ, browse.UPDATE, browse.UPSERT:
		err = kv.store.Load()
		if err != nil {
			return nil, nil, err
		}
	}
	s, err = selector.browse(path.URL)
	return browse.WalkPath(browse.NewWalkState(kv.module), s, path)
}

func (kvs *KeyValuesSelector) browse(parentPath string) (browse.Selection, error) {
	s := &browse.MySelection{}
	path := parentPath
	var keyList []string
	var i int
	var created browse.Selection
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
		if len(key) != 0 {
			if first {
				path = fmt.Sprint(parentPath, "=", key[0].String())
				hasMore = kvs.store.HasValues(path)
			} else {
				return false, nil
			}
		} else {
			if first {
				keyList, err = kvs.store.KeyList(parentPath)
				i = 0
			} else {
				i++
			}
			if hasMore = i < len(keyList); hasMore {
				path = fmt.Sprint(parentPath, "=", keyList[i])
				var key []*browse.Value
				if key, err = browse.CoerseKeys(meta, []string{keyList[i]}); err != nil {
					return false, err
				}
				state.SetKey(key)
			}
		}
		return
	}
	s.OnRead = func (state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		return kvs.store.Value(kvs.metaPath(path, state, meta), meta.GetDataType())
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		childPath := kvs.metaPath(path, state, meta)
		if kvs.store.HasValues(childPath) {
			child, err = kvs.browse(childPath)
		} else if (created != nil) {
			child = created
			created = nil
		}
		return
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, v *browse.Value) (err error) {
		switch op {
		case browse.END_EDIT:
			kvs.store.Save()
		case browse.CREATE_LIST, browse.CREATE_CHILD:
			created, err = kvs.browse(kvs.metaPath(path, state, meta))
		case browse.UPDATE_VALUE:
			kvs.store.SetValue(kvs.metaPath(path, state, meta), v)
		}
		return
	}
	return s, nil
}

func (kvs *KeyValuesSelector) metaPath(parentPath string, state *browse.WalkState, meta schema.Meta) string {
	// TODO: insert keys
	if schema.IsList(state.SelectedMeta()) {
		if len(state.Key()) > 0 {
			return fmt.Sprintf("%s=%s/%s", parentPath, state.Key()[0].String(), meta.GetIdent())
		}
		return parentPath
	}
	return fmt.Sprint(parentPath, "/", meta.GetIdent())
}

func (kvs *KeyValuesSelector) read() (browse.Selection, error) {
	s := &browse.MySelection{}
	return s, nil
}
