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

type PersistableBrowser interface {
	ReadSelector(p *browse.Path) (browse.Selection, error)
	WriteSelector(p *browse.Path, strategy browse.Strategy) (browse.Selection, error)
}

type ComboBrowser struct {
	oper browse.Browser
	persist PersistableBrowser
}

func NewComboBrowser(operational browse.Browser, peristable PersistableBrowser) *ComboBrowser {
	return &ComboBrowser{
		oper:operational,
		persist:peristable,
	}
}

func (self *ComboBrowser) ReadSelector(state *browse.WalkState, p *browse.Path) (browse.Selection, error) {
	var err error
	var oper, operRoot, persist browse.Selection
	var operState *browse.WalkState
	if operRoot, operState, err = self.oper.RootSelector(); err != nil {
		return nil, err
	}

	if oper, operState, err = browse.WalkPath(operState, operRoot, p); err != nil {
		return nil, err
	}
	if state.IsConfig() {
		if persist, err = self.persist.ReadSelector(p); err != nil {
			return nil, err
		}
	}

	return self.readMulticast(oper, persist)
}

func (self *ComboBrowser) readMulticast(oper browse.Selection, persist browse.Selection) (browse.Selection, error) {
	s := &browse.MySelection{}

	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (hasMore bool, err error) {
		// TODO: Figure out how to keep operational and config lists in sync when iterating
		// without keys.
		if persist != nil {
			hasMore, err = persist.Next(state, meta, key, first)
		}
		if oper != nil {
			hasMore, err = oper.Next(state, meta, key, first)
		}
		return
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, val *browse.Value) (err error) {
		if persist != nil &&  state.IsConfig() {
			err = persist.Write(state, meta, op, val)
		}
		return err
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		if persist != nil {
			if state.IsConfig() {
				return persist.Read(state, meta)
			}
		}
		if oper != nil {
			return oper.Read(state, meta)
		}
		return nil, nil
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		var persistChild, operChild browse.Selection
		if persist != nil {
			if state.IsConfig() {
				persistChild, err = persist.Select(state, meta)
			}
		}

		if oper != nil {
			operChild, err = oper.Select(state, meta)
		}
		if operChild != nil || persistChild != nil {
			return self.readMulticast(operChild, persistChild)
		}
		return nil, nil
	}
	return s, nil
}