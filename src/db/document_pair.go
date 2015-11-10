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
	var createdListItem bool
	var postCreateItem browse.Operation
	s.OnNext = func(state *browse.Selection, meta *schema.List, key []*browse.Value, first bool) (browse.Node, error) {
		var err error
		if createdListItem {
			if err = config.Write(state, meta, browse.POST_CREATE_LIST_ITEM, nil); err != nil {
				return nil, err
			}
			createdListItem = false
		}
		var operNext, configNext browse.Node
		if operNext, err = oper.Next(state, meta, key, first); err != nil {
			return nil, err
		}
		if operNext == nil {
			return nil, nil
		}

		if IsContainerConfig {
			configNext, err = config.Next(state, meta, state.Key(), true)
			if err != nil {
				return nil, err
			}
			if configNext == nil {
				if err = config.Write(state, meta, browse.CREATE_LIST_ITEM, nil); err != nil {
					return nil, err
				}
				createdListItem = true
				configNext, err = config.Next(state, meta, state.Key(), true)
				if err != nil {
					return nil, err
				}
				if configNext == nil {
					return nil, errors.New("Could not create config item")
				}
			}
		}
		return self.selectPair(operNext, configNext), nil
	}
	s.OnWrite = func(state *browse.Selection, meta schema.Meta, op browse.Operation, val *browse.Value) (err error) {
		err = oper.Write(state, meta, op, val)
		if err == nil && state.IsConfig() {

			err = config.Write(state, meta, op, val)
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
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList) (browse.Node, error) {
		var err error
		var configChild, operChild browse.Node
		if operChild, err = oper.Select(state, meta); err != nil {
			return nil, err
		}
		if operChild == nil {
			return nil, nil
		}

		if IsContainerConfig && state.IsConfig() {
			if configChild, err = config.Select(state, meta); err != nil {
				return nil, err
			}
			if configChild == nil {
				if schema.IsList(meta) {
					if err = config.Write(state, meta, browse.CREATE_LIST, nil); err != nil {
						return nil, err
					}
					postCreateItem = browse.POST_CREATE_LIST

				} else {
					if err = config.Write(state, meta, browse.CREATE_CONTAINER, nil); err != nil {
						return nil, err
					}
					postCreateItem = browse.POST_CREATE_CONTAINER
				}

				configChild, err = config.Select(state, meta)
				if err != nil {
					return nil, err
				}
				if configChild == nil {
					return nil, errors.New("Could not create config item")
				}
			}
		}
		return self.selectPair(operChild, configChild), nil
	}
	s.OnUnselect = func(state *browse.Selection, meta schema.MetaList) (err error) {
		if postCreateItem > 0 {
			if err = config.Write(state, meta, postCreateItem, nil); err != nil {
				return err
			}
			postCreateItem = browse.Operation(0)
		}
		if err = oper.Unselect(state, meta); err != nil {
			return err
		}
		if IsContainerConfig {
			if err = config.Unselect(state, meta); err != nil {
				return err
			}
		}
		return
	}
	s.OnChoose = func(state *browse.Selection, choice *schema.Choice) (choosen schema.Meta, err error) {
		choosen, err = oper.Choose(state, choice)
		return
	}
	return s
}
