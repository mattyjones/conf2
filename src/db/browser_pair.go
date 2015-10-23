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
	return browse.Upsert(browse.NewPath(""), self.config, self.oper)
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
	return combo, state, nil
}

func (self *BrowserPair) Schema() schema.MetaList {
	m := self.oper.Schema()
	if m == nil {
		m = self.config.Schema()
	}
	return m
}

func (self *BrowserPair) selectPair(oper browse.Selection, config browse.Selection) (browse.Selection, error) {
	s := &browse.MySelection{}
	IsContainerConfig := config != nil
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (next browse.Selection, err error) {
		var operNext, configNext browse.Selection
		operNext, err = oper.Next(state, meta, key, first)
		if err == nil && operNext != nil {
			if IsContainerConfig {
				configNext, err = config.Next(state, meta, state.Key(), true)
				if err != nil {
					return nil, err
				}
			}
			return self.selectPair(operNext, configNext)
		}

		return
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, val *browse.Value) (err error) {
		err = oper.Write(state, meta, op, val)
		if err == nil && state.IsConfig() {
			err = config.Write(state, meta, op, val)
			// TODO: if there's now an error, config and operation are out of sync. To fix
			// this we must "rollback the Write
		}

		return err
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (v *browse.Value, err error) {
		if IsContainerConfig && state.IsConfig() {
			v, err = config.Read(state, meta)
		}
		if v == nil && err == nil {
			v, err = oper.Read(state, meta)
		}
		return
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		var configChild, operChild browse.Selection
		operChild, err = oper.Select(state, meta)
		if operChild != nil {
			if IsContainerConfig && state.IsConfig() {
				configChild, err = config.Select(state, meta)
			}
			return self.selectPair(operChild, configChild)
		}
		return nil, nil
	}
	s.OnChoose = func(state *browse.WalkState, choice *schema.Choice) (choosen schema.Meta, err error) {
		choosen, err = oper.Choose(state, choice)
		return
	}
	return s, nil
}