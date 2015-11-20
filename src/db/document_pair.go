package db

import (
	"schema"
	"schema/browse"
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

type DocumentPair struct {
	oper   browse.Data
	config browse.Data
}

func NewDocumentPair(operational browse.Data, config browse.Data) (pair *DocumentPair, err error) {
	pair = &DocumentPair{
		oper:   operational,
		config: config,
	}
	// Here we initialize the operational browser with the current configuration
	err = browse.SyncData(config, operational, browse.NewPath(""), browse.UPSERT)
	return
}

func NewSelectionPair(oper *browse.Selection, config *browse.Selection) (s *browse.Selection, err error) {
	pair := &selectionPair{
		operEvents : &browse.EventsImpl{Parent:oper.Events},
	}
	var configNode browse.Node
	if config != nil {
		pair.configEvents = &browse.EventsImpl{Parent:config.Events}
		configNode = config.Node
	}
	combo := pair.selectPair(oper.State, oper.Node, configNode)
	s = &browse.Selection{
		Events: &browse.EventMulticast{
			A: pair.configEvents,
			B: pair.operEvents,
		},
		State: oper.State,
		Node: combo,
	}
	return
}

type selectionPair struct {
	configEvents browse.Events
	operEvents browse.Events
}

func (self *DocumentPair) Selector(path *browse.Path) (s *browse.Selection, err error) {
	var configSel, operSel *browse.Selection
	if operSel, err = self.oper.Selector(path); err != nil {
		return
	}
	if operSel == nil {
		return nil, nil
	}
	if configSel, err = self.config.Selector(path); err != nil {
		return
	}
	return NewSelectionPair(operSel, configSel)
}

func (self *DocumentPair) Schema() schema.MetaList {
	return self.oper.Schema()
}

func (self *selectionPair) selectPair(state *browse.WalkState, operNode browse.Node, configNode browse.Node) browse.Node {
	oper := &browse.Selection{
		Events: self.operEvents,
		Node: operNode,
	}
	var config *browse.Selection
	IsContainerConfig := configNode != nil
	if IsContainerConfig {
		config = &browse.Selection{
			Events: self.configEvents,
			Node: configNode,
		}
	}
	s := &browse.MyNode{}
	s.OnNext = func(sel *browse.Selection, meta *schema.List, new bool, key []*browse.Value, first bool) (browse.Node, error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		var err error
		var operNext, configNext browse.Node
		if operNext, err = operNode.Next(oper, meta, new, key, first); err != nil {
			return nil, err
		}
		if operNext == nil {
			return nil, nil
		}
		if IsContainerConfig {
			configNext, err = configNode.Next(config, meta, new, sel.State.Key(), true)
			if err != nil {
				return nil, err
			}
		}
		return self.selectPair(sel.State, operNext, configNext), nil
	}
	s.OnWrite = func(sel *browse.Selection, meta schema.HasDataType, val *browse.Value) (err error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		err = operNode.Write(oper, meta, val)
		if err == nil && sel.IsConfig() {
			err = configNode.Write(config, meta, val)
			// TODO: if there's now an error, config and operation are out of sync. To fix
			// this we must "rollback the Write
		}
		return err
	}
	s.OnRead = func(sel *browse.Selection, meta schema.HasDataType) (v *browse.Value, err error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		if IsContainerConfig && sel.IsConfig() {
			v, err = configNode.Read(config, meta)
		}
		if v == nil && err == nil {
			v, err = operNode.Read(oper, meta)
		}
		return
	}
	s.OnSelect = func(sel *browse.Selection, meta schema.MetaList, createOk bool) (browse.Node, error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		var err error
		var configChild, operChild browse.Node
		if operChild, err = operNode.Select(oper, meta, createOk); err != nil {
			return nil, err
		}
		if operChild == nil {
			return nil, nil
		}

		if IsContainerConfig && sel.IsConfig() {
			if configChild, err = configNode.Select(config, meta, createOk); err != nil {
				return nil, err
			}
		}
		return self.selectPair(sel.State, operChild, configChild), nil
	}
	s.OnChoose = func(sel *browse.Selection, choice *schema.Choice) (m schema.Meta, err error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		return operNode.Choose(sel, choice)
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) (err error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		if err = operNode.Event(oper, e); err != nil {
			return err
		}
		if IsContainerConfig && sel.IsConfig() {
			err = configNode.Event(config, e)
		}
		return
	}
	return s
}
