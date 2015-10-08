package db
import (
	"schema/browse"
	"schema"
	"fmt"
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

type BrowserPair struct {
	oper browse.Browser
	config browse.Browser
}

func NewBrowserPair(operational browse.Browser, config browse.Browser) *BrowserPair {
	return &BrowserPair{
		oper:operational,
		config:config,
	}
}

func (self *BrowserPair) Init() error {
	// Here we initialize the operational browser with the current configuration
	return browse.Update(browse.NewPath(""), self.config, self.oper)
}

func (self *BrowserPair) Selector(path *browse.Path, strategy browse.Strategy) (browse.Selection, *browse.WalkState, error) {
	var err error
	var oper, config, combo browse.Selection
	var operState, configState *browse.WalkState
	if oper, operState, err = self.oper.Selector(path, strategy); err != nil {
		return nil, nil, err
	}

	if config, configState, err = self.config.Selector(path, strategy); err != nil {
		return nil, nil, err
	}

	if config == nil && oper == nil {
		return nil, nil, browse.NotFound(path.URL)
	}
	if combo, err = self.selectPair(oper, config); err != nil {
		return nil, nil, err
	}
	state := operState
	if state == nil {
		state = configState
	}
fmt.Printf("browser_pair - path %s, config nil %v, oper nil %v\n", path.URL, config == nil, oper == nil)
	return combo, state, nil
}

func (self *BrowserPair) Module() *schema.Module {
	m := self.oper.Module()
	if m == nil {
		m = self.config.Module()
	}
	return m
}

func (self *BrowserPair) selectPair(oper browse.Selection, config browse.Selection) (browse.Selection, error) {
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
fmt.Printf("browser_pair - OnWrite\n")
		if oper != nil {
			err = oper.Write(state, meta, op, val)
		}
		if err == nil && config != nil && state.IsConfig() {
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
			return self.selectPair(operChild, configChild)
		}
		return nil, nil
	}
	s.OnChoose = func(state *browse.WalkState, choice *schema.Choice) (choosen schema.Meta, err error) {
		if oper != nil {
			choosen, err = oper.Choose(state, choice)
		}
		// assume that the error is because it's not implemented
		if err != nil && config != nil {
			choosen, err = config.Choose(state, choice)
		}
		return
	}
	return s, nil
}