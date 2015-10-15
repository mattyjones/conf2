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
	state = browse.NewWalkState(kv.module)
	if strategy == browse.READ {
		// we walk to the destination and legitimately return nil if nothing is there
		s, err = selector.selectConfig("")
		return browse.WalkPath(state, s, path)
	}

	// here we fast-forward to the destination prepared to insert into the parse hierarchy
	s, _ = selector.selectConfig(path.URL)
	state, err = fastForwardState(state, path)
fmt.Printf("config - s == nil? %v  state %s\n", s == nil, state.String())
	return
}

func fastForwardState(initialState *browse.WalkState, path *browse.Path) (state *browse.WalkState, err error) {
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
	if len(sp.parent) > 0 {
		return fmt.Sprint(sp.String(), "/", metaIdent)
	}
	return metaIdent
}


func (kvs *configSelector) selectConfig(parentPath string) (browse.Selection, error) {
	s := &browse.MySelection{}
	path := configPath{parent:parentPath}
	var keyList []string
	var i int
	var created browse.Selection
	s.OnChoose = func(state *browse.WalkState, choice *schema.Choice) (m schema.Meta, err error) {
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
				candidatePath := path.metaPath(prop.GetIdent())
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
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
fmt.Printf("config - OnNext %s\n", path.metaPath(meta.GetIdent()))
		path.listKey = ""
		if len(key) > 0 {
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
fmt.Printf("config - OnRead %s\n", path.metaPath(meta.GetIdent()))
		return kvs.store.Value(path.metaPath(meta.GetIdent()), meta.GetDataType())
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
fmt.Printf("config - OnSelect %s\n", path.metaPath(meta.GetIdent()))
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
fmt.Printf("config - OnWrite %s\n", path.metaPath(meta.GetIdent()))
		switch op {
		case browse.END_EDIT:
			kvs.store.Save()
		case browse.CREATE_LIST, browse.CREATE_CHILD:
			childPath := path.metaPath(meta.GetIdent())
			created, err = kvs.selectConfig(childPath)
		case browse.CREATE_LIST_ITEM:
			path.listKey = state.Key()[0].String()
		case browse.UPDATE_VALUE:
			if err = kvs.store.SetValue(path.metaPath(meta.GetIdent()), v); err == nil {
				if schema.IsKeyLeaf(state.SelectedMeta(), meta) {
					oldPath := path.String()
					pathCopy := path
					pathCopy.listKey = v.String()
					newPath := pathCopy.String()
					kvs.store.RenameKey(oldPath, newPath)
				}
			}
		}
		return
	}
	return s, nil
}

func (kvs *configSelector) read() (browse.Selection, error) {
	s := &browse.MySelection{}
	return s, nil
}
