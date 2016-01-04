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

func (t *NodeTable) Next() (error) {
	// Container
	if ! schema.IsList(t.Corner.State.SelectedMeta()) {
		if t.Row == nil {
			t.Row = t.Corner
		} else {
			t.Row = nil
		}
		return nil
	}

	// List
	meta := t.Corner.State.SelectedMeta().(*schema.List)
	rowNode, err := t.Corner.Node.Next(t.Corner, meta, t.autoCreate, data.NO_KEYS, t.Row == nil)
	if err != nil {
		return err
	}
	if rowNode == nil {
		t.Row = nil
	} else {
		t.Row = t.Corner.SelectListItem(rowNode, t.Corner.State.Key())
	}
	t.sels = make(map[string]*data.Selection)
	t.vals = make(map[string]*data.Value)
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
	s, err := data.SelectMetaList(sel, ident, t.autoCreate)
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
		sel, err = data.SelectMetaList(parentSel, ident, autoCreate)
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
	if v, err = data.GetValue(sel, ident); err != nil {
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
	return data.ChangeValue(sel, ident, v)
}