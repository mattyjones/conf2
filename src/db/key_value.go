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

func (kvs KeyValueStore) Load(p *browse.Path) (map[string]*browse.Value, error) {
	return kvs, nil
}

func (kvs KeyValueStore) Clear(p *browse.Path) error {
	for k, _ := range kvs {
		delete(kvs, k)
	}
	return nil
}

func (kvs KeyValueStore) Upsert(p *browse.Path, vals map[string]*browse.Value) error {
	for k, v := range vals {
		kvs[k] = v
	}
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
	values map[string]*browse.Value
}

func (kv *KeyValues) Module() *schema.Module {
	return kv.module
}

func (kv *KeyValues) Selector(path *browse.Path, strategy browse.Strategy) (s browse.Selection, state *browse.WalkState, err error) {
	selector := &KeyValuesSelector{path: path, store:kv.store, strategy : strategy}
	switch strategy {
	case browse.READ, browse.UPDATE:
		selector.values, err = kv.store.Load(path)
		if err != nil {
			return nil, nil, err
		}
	case browse.INSERT, browse.UPSERT:
		selector.values = make(map[string]*browse.Value, 100)
	}
	s, err = selector.browse(path.URL)
	return browse.WalkPath(browse.NewWalkState(kv.module), s, path)
}

func (kvs *KeyValuesSelector) browse(parentPath string) (browse.Selection, error) {
	s := &browse.MySelection{}
	var list [][]*browse.Value
	var i int
	var created browse.Selection
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
		if first {
			list, err = kvs.buildList(kvs.metaPath(parentPath, state, meta), meta)
			if err != nil {
				return false, err
			}
			i = 0
		}
		if i < len(list) {
			state.SetKey(list[i])
			i++
			return true, nil
		}
		return false, nil
	}
	s.OnRead = func (state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		v, _ := kvs.values[kvs.metaPath(parentPath, state, meta)]
		return v, nil
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		childPath := kvs.metaPath(parentPath, state, meta)
		if kvs.containerExists(childPath) {
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
			switch kvs.strategy {
			case browse.INSERT:
				kvs.store.Clear(kvs.path)
				kvs.store.Upsert(kvs.path, kvs.values)
			case browse.UPSERT, browse.UPDATE:
				// TODO: differentiate here
				kvs.store.Upsert(kvs.path, kvs.values)
			}
		case browse.CREATE_LIST, browse.CREATE_CHILD:
			created, err = kvs.browse(kvs.metaPath(parentPath, state, meta))
		case browse.UPDATE_VALUE:
			kvs.values[kvs.metaPath(parentPath, state, meta)] = v
		}
		return
	}
	return s, nil
}

func (kvs *KeyValuesSelector) containerExists(path string) bool {
	// TODO: performance - most efficient way? sort first?
	for k, _ := range kvs.values {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (kvs *KeyValuesSelector) buildList(path string, meta *schema.List) (keys [][]*browse.Value, err error) {
	// TODO: performance - most efficient way? sort first?
	keysSet := make(map[string]struct{}, 10)
	keyStart := len(path) + 1
	for k, _ := range kvs.values {
		if strings.HasPrefix(k, path) {
			keyEnd := strings.IndexRune(k[keyStart:], '/')
			if keyEnd < 0 {
				continue
			}
			key := k[keyStart:keyStart + keyEnd]
			keysSet[key] = struct{}{}
		}
	}
	keys = make([][]*browse.Value, len(keysSet))
	var i int
	for k, _ := range keysSet {
		if keys[i], err = browse.CoerseKeys(meta, []string{k}); err != nil {
			return nil, err
		}
		i++
	}

	return
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
