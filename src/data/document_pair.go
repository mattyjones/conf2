package data

import (
	"schema"
)

// Details on config nodes v.s. operational nodes in Section 7.21.1 of RFC6020
// ============================================================================
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
	var edit *Editor
	if edit, err = PathToPath(config, operational, ""); err == nil {
		err = edit.Upsert()
	}

	return
}

type selectionPair struct {
	configEvents Events
	operEvents Events
}

func (self *DocumentPair) Node() (Node) {
	pair := &selectionPair{
		configEvents : &EventsImpl{},
		operEvents : &EventsImpl{},
	}
	return pair.selectPair(self.oper.Node(), self.config.Node())
}

func (self *DocumentPair) Schema() schema.MetaList {
	return self.oper.Schema()
}

func (self *selectionPair) selectPair(operNode Node, configNode Node) Node {
	var oper *Selection
	var config *Selection
	IsContainerConfig := configNode != nil
	s := &MyNode{}
	onInit := func(sel *Selection) (err error) {
		if IsContainerConfig && config == nil {
			config = &Selection{
				Events: self.configEvents,
				Node: configNode,
			}
		}
		if oper == nil {
			oper = &Selection{
				Events: self.operEvents,
				Node: operNode,
			}
		}
		if config != nil {
			config.State = sel.State
		}
		oper.State = sel.State
		return
	}
	s.OnEvent = func(sel *Selection, e Event) error {
		var err error
		if err = onInit(sel); err != nil {
			return err
		}
		if err = operNode.Event(oper, e); err != nil {
			return err
		}
		self.operEvents.Fire(sel.State.path, e)
		if IsContainerConfig && sel.IsConfig() {
			if err = configNode.Event(config, e); err != nil {
				return err
			}
			if err = self.configEvents.Fire(sel.State.path, e); err != nil {
				return err
			}
		}
		return err
	}
	s.OnNext = func(sel *Selection, meta *schema.List, new bool, key []*Value, first bool) (Node, error) {
		var err error
		if err = onInit(sel); err != nil {
			return nil, err
		}
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
			if configNext == nil && ! new {
				configNext, err = configNode.Next(config, meta, true, sel.State.Key(), true)
				if err != nil {
					return nil, err
				}
			}
		}
		return self.selectPair(operNext, configNext), nil
	}
	s.OnWrite = func(sel *Selection, meta schema.HasDataType, val *Value) (err error) {
		if err = onInit(sel); err != nil {
			return err
		}
		err = operNode.Write(oper, meta, val)
		if err == nil && sel.IsConfig() {
			err = configNode.Write(config, meta, val)
			// TODO: if there's now an error, config and operation are out of sync. To fix
			// this we must "rollback the Write
		}
		return err
	}
	s.OnRead = func(sel *Selection, meta schema.HasDataType) (v *Value, err error) {
		if err = onInit(sel); err != nil {
			return nil, err
		}
		if IsContainerConfig && sel.IsConfig() {
			v, err = configNode.Read(config, meta)
		}
		if v == nil && err == nil {
			v, err = operNode.Read(oper, meta)
		}
		return
	}
	s.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (Node, error) {
		var err error
		if err = onInit(sel); err != nil {
			return nil, err
		}
		var configChild, operChild Node
		if operChild, err = operNode.Select(oper, meta, new); err != nil {
			return nil, err
		}
		if operChild == nil {
			return nil, nil
		}

		if IsContainerConfig && sel.IsConfig() {
			if configChild, err = configNode.Select(config, meta, new); err != nil {
				return nil, err
			}
			if configChild == nil && ! new {
				if configChild, err = configNode.Select(config, meta, true); err != nil {
					return nil, err
				}
			}
		}
		return self.selectPair(operChild, configChild), nil
	}
	s.OnChoose = func(sel *Selection, choice *schema.Choice) (m schema.Meta, err error) {
		if err = onInit(sel); err != nil {
			return nil, err
		}
		return operNode.Choose(sel, choice)
	}
	// at this time, config has no generic ability to handle actions, so forward to operational.
	s.OnAction = operNode.Action
	return s
}
