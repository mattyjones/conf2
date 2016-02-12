package process

import (
	"schema"
	"data"
	"strings"
)

type Table interface {
	Set(identPath string, val interface{}) error
	Get(identPath string) (interface{}, error)
	Select(identPath string, autocreate bool) (t Table, err error)
	HasNext() bool
	Next() error
}

type NodeTable struct {
	Corner     *data.Selection
	Row        *data.Selection
	autoCreate bool
	sels       map[string]*data.Selection
	vals       map[string]*data.Value
}

func (t *NodeTable) HasNext() (bool) {
	return t.Row != nil
}

func (self *NodeTable) Next() (error) {
	// Container
	if ! schema.IsList(self.Corner.Meta()) {
		if self.Row == nil {
			self.Row = self.Corner
		} else {
			self.Row = nil
		}
		return nil
	}

	// List
	r := data.ListRequest{
		Meta: self.Corner.Meta().(*schema.List),
		First: self.Row == nil,
		New : self.autoCreate,
	}
	r.New = self.autoCreate
	rowNode, key, err := self.Corner.Node().Next(self.Corner, r)
	if err != nil {
		return err
	}
	if rowNode == nil {
		self.Row = nil
	} else {
		self.Row = self.Corner.SelectListItem(rowNode, key)
	}
	self.sels = make(map[string]*data.Selection)
	self.vals = make(map[string]*data.Value)
	return nil
}

func (t *NodeTable) IsDot(r rune) bool {
	return r == '.'
}

func (t *NodeTable) getSelection(path string) (*data.Selection, error) {
	s, found := t.sels[path]
	if found {
		return s, nil
	}
	sel := t.Row
	ident := path
	dot := strings.LastIndexFunc(path, t.IsDot)
	if dot > 0 {
		ident = path[dot + 1:]
		var selErr error
		if sel, selErr = t.getSelection(path[:dot]); selErr != nil {
			return nil, selErr
		}
	}

	// if s is a list, what do we do?  select 0th?
	s, err := sel.FindOrCreate(ident, t.autoCreate)
	if err != nil {
		return nil, err
	}
	if t.sels == nil {
		t.sels = make(map[string]*data.Selection)
	}
	t.sels[path] = s
	return s, nil
}

func (t *NodeTable) Select(identPath string, autoCreate bool) (Table, error) {
	parentSel, ident, err := t.resolveIdentPath(identPath)
	if parentSel != nil && err == nil {
		var sel *data.Selection
		sel, err = parentSel.FindOrCreate(ident, autoCreate)
		if err != nil {
			return nil, err
		}
		if sel != nil {
			return &NodeTable{
				Corner: sel,
				autoCreate: autoCreate,
			}, nil
		}
	}
	return nil, err
}

func (t *NodeTable) resolveIdentPath(identPath string) (sel *data.Selection, ident string, err error) {
	dot := strings.LastIndexFunc(identPath, t.IsDot)
	sel = t.Row
	ident = identPath
	if dot > 0 {
		if sel, err = t.getSelection(identPath[:dot]); err != nil {
			return
		}
		ident = identPath[dot + 1:]
	}
	return
}

func (t *NodeTable) Get(identPath string) (interface{}, error) {
	var err error
	v, found := t.vals[identPath]
	if found {
		return v, nil
	}
	sel, ident, err := t.resolveIdentPath(identPath)
	if err != nil {
		return nil, err
	}
	if sel == nil {
		return nil, err
	}
	if v, err = sel.GetValue(ident); err != nil {
		return nil, err
	}
	if t.vals == nil {
		t.vals = make(map[string]*data.Value)
	}
	t.vals[identPath] = v
	if v == nil {
		return nil, nil
	}
	return v.Value(), nil
}

func (t *NodeTable) Set(identPath string, v interface{}) error {
	if v == nil {
		return nil
	}
	sel, ident, err := t.resolveIdentPath(identPath)
	if err != nil {
		return err
	}
	return sel.Set(ident, v)
}