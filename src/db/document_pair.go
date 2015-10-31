package db
import (
	"schema/browse"
	"schema"
	"errors"
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

type DocumentPair struct {
	oper browse.Data
	config browse.Data
}

func NewDocumentPair(operational browse.Data, config browse.Data) *DocumentPair {
	return &DocumentPair{
		oper:operational,
		config:config,
	}
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
fmt.Printf("document_pair configSel=%p\n", configSel)

	if configSel == nil && operSel == nil {
		return nil, browse.NotFound(path.URL)
	} else if operSel == nil {
		return nil, errors.New("Illegal state")
	}
	var configNode, operNode, comboNode browse.Node
	operNode = operSel.Node()
	if configSel != nil {
		configNode = configSel.Node()
	}
fmt.Printf("document_pair configNode=%p\n", configNode)
	if comboNode, err = self.selectPair(operNode, configNode); err != nil {
		return nil, err
	}
	return operSel.Copy(comboNode), nil
}

func (self *DocumentPair) Schema() schema.MetaList {
	m := self.oper.Schema()
	if m == nil {
		m = self.config.Schema()
	}
	return m
}

func (self *DocumentPair) selectPair(oper browse.Node, config browse.Node) (browse.Node, error) {
	s := &browse.MyNode{}
if config == nil {
	panic("STOP")
}
	IsContainerConfig := config != nil
	s.OnNext = func(state *browse.Selection, meta *schema.List, key []*browse.Value, first bool) (next browse.Node, err error) {
		var operNext, configNext browse.Node
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
	s.OnSelect = func(state *browse.Selection, meta schema.MetaList) (child browse.Node, err error) {
		var configChild, operChild browse.Node
		operChild, err = oper.Select(state, meta)
		if operChild != nil {
fmt.Printf("document_pair IsContainerConfig %v && state.IsConfig() %v\n", IsContainerConfig, state.IsConfig())
			if IsContainerConfig && state.IsConfig() {
				configChild, err = config.Select(state, meta)
			}
			return self.selectPair(operChild, configChild)
		}
		return nil, nil
	}
	s.OnChoose = func(state *browse.Selection, choice *schema.Choice) (choosen schema.Meta, err error) {
		choosen, err = oper.Choose(state, choice)
		return
	}
	return s, nil
}