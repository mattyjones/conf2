package db
import (
	"schema"
	"schema/browse"
	"fmt"
	"errors"
)

type Config struct {
	module *schema.Module
	store Store
}

func NewConfig(module *schema.Module, store Store) *Config {
	return &Config{
		module : module,
		store : store,
	}
}

func (kv *Config) Module() *schema.Module {
	return kv.module
}

func (kv *Config) Selector(path *browse.Path, strategy browse.Strategy) (s browse.Selection, state *browse.WalkState, err error) {
	selector := &configSelector{store:kv.store, strategy : strategy}
	switch strategy {
	case browse.READ, browse.UPDATE, browse.UPSERT:
		err = kv.store.Load()
		if err != nil {
			return nil, nil, err
		}
	}
	if strategy == browse.READ {
		// we walk to the destination and legitimately return nil if nothing is there
		s, err = selector.selectConfig("")
		return browse.WalkPath(browse.NewWalkState(kv.module), s, path)
	}

	// here we fast-forward to the destination prepared to insert into the parse hierarchy
	s, _ = selector.selectConfig(path.URL)
	state, err = FastForward(browse.NewWalkState(kv.module), path)
	return
}

func FastForward(initialState *browse.WalkState, path *browse.Path) (state *browse.WalkState, err error) {
	state = initialState
	for _, seg := range path.Segments {
		position := schema.FindByIdentExpandChoices(state.SelectedMeta(), seg.Ident)
		if position == nil {
			return nil, errors.New(fmt.Sprintf("%s.%s not found in schema", state.SelectedMeta().GetIdent(), seg.Ident))
		}
		state.SetPosition(position)
		state = state.Select()
	}
	return state, nil
}

type configSelector struct {
	store Store
	strategy browse.Strategy
}


type configPath struct {
	parent string
	listKey string
}

func (sp configPath) String() string {
	p := sp.parent
	if len(sp.listKey) > 0 {
		p = fmt.Sprint(p, "=", sp.listKey)
	}
	return p
}

func (sp configPath) metaPath(metaIdent string) string {
	return fmt.Sprint(sp.String(), "/", metaIdent)
}


func (kvs *configSelector) selectConfig(parentPath string) (browse.Selection, error) {
	s := &browse.MySelection{}
	path := configPath{parent:parentPath}
	var keyList []string
	var i int
	var created browse.Selection
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
		path.listKey = ""
		if len(key) != 0 {
			if first {
				path.listKey = key[0].String()
				hasMore = kvs.store.HasValues(path.String())
			} else {
				return false, nil
			}
		} else {
			if first {
				keyList, err = kvs.store.KeyList(parentPath, meta)
				i = 0
			} else {
				i++
			}
			if hasMore = i < len(keyList); hasMore {
				path.listKey = keyList[i]
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
		return kvs.store.Value(path.metaPath(meta.GetIdent()), meta.GetDataType())
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		childPath := path.metaPath(meta.GetIdent())
		if kvs.store.HasValues(childPath) {
			child, err = kvs.selectConfig(childPath)
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
			childPath := path.metaPath(meta.GetIdent())
			created, err = kvs.selectConfig(childPath)
		case browse.CREATE_LIST_ITEM:
			path.listKey = state.Key()[0].String()
		case browse.UPDATE_VALUE:
			kvs.store.SetValue(path.metaPath(meta.GetIdent()), v)
		}
		return
	}
	return s, nil
}

func (kvs *configSelector) read() (browse.Selection, error) {
	s := &browse.MySelection{}
	return s, nil
}
