package data

import (
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

type DocumentPair struct {
	oper   Data
	config Data
}

func NewDocumentPair(operational Data, config Data) (pair *DocumentPair, err error) {
	pair = &DocumentPair{
		oper:   operational,
		config: config,
	}
	// Here we initialize the operational browser with the current configuration
	err = SyncData(config, operational, NewPath(""), UPSERT)
	return
}

func NewSelectionPair(oper *Selection, config *Selection) (s *Selection, err error) {
	pair := &selectionPair{
		operEvents : &EventsImpl{Parent:oper.Events},
	}
	var configNode Node
	if config != nil {
		pair.configEvents = &EventsImpl{Parent:config.Events}
		configNode = config.Node
	}
	combo := pair.selectPair(oper.State, oper.Node, configNode)
	s = &Selection{
		Events: &EventMulticast{
			A: pair.configEvents,
			B: pair.operEvents,
		},
		State: oper.State,
		Node: combo,
	}
	return
}

type selectionPair struct {
	configEvents Events
	operEvents Events
}

func (self *DocumentPair) Selector(path *Path) (s *Selection, err error) {
	var configSel, operSel *Selection
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

func (self *selectionPair) selectPair(state *WalkState, operNode Node, configNode Node) Node {
	oper := &Selection{
		Events: self.operEvents,
		Node: operNode,
	}
	var config *Selection
	IsContainerConfig := configNode != nil
	if IsContainerConfig {
		config = &Selection{
			Events: self.configEvents,
			Node: configNode,
		}
	}
	s := &MyNode{}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (Node, error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		var err error
		var operNext, configNext Node
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
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *Value) (err error) {
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
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (v *Value, err error) {
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
	s.OnSelect = func(sel *Selection, meta schema.MetaList, createOk bool) (Node, error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		var err error
		var configChild, operChild Node
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
	s.OnChoose = func(sel *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		return operNode.Choose(sel, choice)
	}
	s.OnEvent = func(sel *Selection, e Event) (err error) {
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
	// at this time, config has no generic ability to handle actions, so forward to operational.
	s.OnAction = operNode.Action
	return s
}
