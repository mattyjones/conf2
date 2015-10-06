package db
import (
	"schema/browse"
	"schema"
)

// Details on config nodes v.s. state data
// Section 7.21.1 of RFC6020
// =========================
//   If "config" is not specified, the default is the same as the parent
//   schema node's "config" value.  If the parent node is a "case" node,
//   the value is the same as the "case" node's parent "choice" node.
//
//   If the top node does not specify a "config" statement, the default is
//   "true".
//
//   If a node has "config" set to "false", no node underneath it can have
//   "config" set to "true".


type ComboBrowser struct {
	oper browse.Browser
	config browse.Browser
}

func NewComboBrowser(operational browse.Browser, config browse.Browser) *ComboBrowser {
	return &ComboBrowser{
		oper:operational,
		config:config,
	}
}

func (self *ComboBrowser) Selector(path *browse.Path, strategy browse.Strategy) (browse.Selection, *browse.WalkState, error) {
	var err error
	var oper, config, combo browse.Selection
	var operState, configState *browse.WalkState
	if oper, operState, err = self.oper.Selector(path, strategy); err != nil {
		return nil, nil, err
	} else if oper != nil {
		if oper, operState, err = browse.WalkPath(operState, oper, path); err != nil {
			return nil, nil, err
		}
	}

	if config, configState, err = self.config.Selector(path, strategy); err != nil {
		return nil, nil, err
	} else {
		if config, configState, err = browse.WalkPath(configState, config, path); err != nil {
			return nil, nil, err
		}
	}

	if combo, err = self.readMulticast(oper, config); err != nil {
		return nil, nil, err
	}
	state := operState
	if state == nil {
		state = configState
	}
	return combo, state, nil
}

func (self *ComboBrowser) Module() *schema.Module {
	m := self.oper.Module()
	if m == nil {
		m = self.config.Module()
	}
	return m
}

func (self *ComboBrowser) readMulticast(oper browse.Selection, config browse.Selection) (browse.Selection, error) {
	s := &browse.MySelection{}

	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
		// TODO: Figure out how to keep operational and config lists in sync when iterating
		// without keys.
		if config != nil {
			hasMore, err = config.Next(state, meta, key, first)
		}
		if oper != nil && !hasMore {
			hasMore, err = oper.Next(state, meta, key, first)
		}
		return
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, val *browse.Value) (err error) {
		if config != nil &&  state.IsConfig() {
			err = config.Write(state, meta, op, val)
		}
		return err
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		if config != nil {
			if state.IsConfig() {
				return config.Read(state, meta)
			}
		}
		if oper != nil {
			return oper.Read(state, meta)
		}
		return nil, nil
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		var configChild, operChild browse.Selection
		if config != nil {
			if state.IsConfig() {
				configChild, err = config.Select(state, meta)
			}
		}

		if oper != nil {
			operChild, err = oper.Select(state, meta)
		}
		if operChild != nil || configChild != nil {
			return self.readMulticast(operChild, configChild)
		}
		return nil, nil
	}
	return s, nil
}