package db

import (
	"errors"
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
	p := browse.NewPath("")
	var src, dest *browse.Selection
	if src, err = config.Selector(p); err != nil {
		return nil, err
	}
	if dest, err = operational.Selector(p); err != nil {
		return nil, err
	}
	err = browse.Upsert(src, dest)
	return
}

func (self *DocumentPair) Init() (err error) {
	// Here we initialize the operational browser with the current configuration
	var sConfig, sOper *browse.Selection
	if sConfig, err = self.config.Selector(browse.NewPath("")); err != nil {
		return err
	}
	if sOper, err = self.oper.Selector(browse.NewPath("")); err != nil {
		return err
	}
	return browse.Upsert(sConfig, sOper)
}

func (self *DocumentPair) Selector(path *browse.Path) (*browse.Selection, error) {
	var err error
	var operSel, configSel *browse.Selection
	if operSel, err = self.oper.Selector(path); err != nil {
		return nil, err
	}

	if configSel, err = self.config.Selector(path); err != nil {
		return nil, err
	}

	if configSel == nil && operSel == nil {
		return nil, nil
	}
	if configSel != nil && operSel != nil {
		return operSel.Copy(self.selectPair(operSel.Node(), configSel.Node())), nil
	}
	if operSel == nil {
		return nil, errors.New("Operational should exist, illegal state")
	}
	if configSel == nil {
		var configRoot *browse.Selection
		if configRoot, err = self.config.Selector(browse.NewPath("")); err != nil {
			return nil, err
		}
		var operRoot *browse.Selection
		if operRoot, err = self.oper.Selector(browse.NewPath("")); err != nil {
			return nil, err
		}
		combo := self.selectPair(operRoot.Node(), configRoot.Node())
		return browse.WalkPath(browse.NewSelection(combo, operRoot.SelectedMeta()), path)
	}
	combo := self.selectPair(operSel.Node(), configSel.Node())
	return operSel.Copy(combo), nil
}

func (self *DocumentPair) Schema() schema.MetaList {
	m := self.oper.Schema()
	if m == nil {
		m = self.config.Schema()
	}
	return m
}

func (self *DocumentPair) selectPair(oper browse.Node, config browse.Node) browse.Node {
	s := &browse.MyNode{}
	IsContainerConfig := config != nil
	s.OnNext = func(state *browse.Selection, meta *schema.List, new bool, key []*browse.Value, first bool) (browse.Node, error) {
		var err error
		var operNext, configNext browse.Node
		if operNext, err = oper.Next(state, meta, new, key, first); err != nil {
			return nil, err
		}
		if operNext == nil {
			return nil, nil
		}
		if IsContainerConfig {
			configNext, err = config.Next(state, meta, new, state.Key(), true)
			if err != nil {
				return nil, err
			}
		}
		return self.selectPair(operNext, configNext), nil
	}
	s.OnWrite = func(state *browse.Selection, meta schema.HasDataType, val *browse.Value) (err error) {
		err = oper.Write(state, meta, val)
		if err == nil && state.IsConfig() {

			err = config.Write(state, meta, val)
			// TODO: if there's now an error, config and operation are out of sync. To fix
			// this we must "rollback the Write
		}

		return err
	}
	s.OnRead = func(state *browse.Selection, meta schema.HasDataType) (v *browse.Value, err error) {
		if IsContainerConfig && state.IsConfig() {
			v, err = config.Read(state, meta)
		}
		if v == nil && err == nil {
			v, err = oper.Read(state, meta)
		}
		return
	}
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList, createOk bool) (browse.Node, error) {
		var err error
		var configChild, operChild browse.Node
		if operChild, err = oper.Select(state, meta, createOk); err != nil {
			return nil, err
		}
		if operChild == nil {
			return nil, nil
		}

		if IsContainerConfig && state.IsConfig() {
			if configChild, err = config.Select(state, meta, createOk); err != nil {
				return nil, err
			}
		}
		return self.selectPair(operChild, configChild), nil
	}
	s.OnChoose = func(state *browse.Selection, choice *schema.Choice) (choosen schema.Meta, err error) {
		choosen, err = oper.Choose(state, choice)
		return
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) (err error) {
		if err = oper.Event(sel, e); err != nil {
			return err
		}
		if IsContainerConfig && sel.IsConfig() {
			err = config.Event(sel, e)
		}
		return
	}
	return s
}
