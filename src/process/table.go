package process

import (
	"schema"
	"data"
	"strings"
)

type Table struct {
	Corner    *data.Selection
	Row       *data.Selection
	autoCreate  bool
	exhausted bool
	sels      map[string]*data.Selection
	vals      map[string]*schema.Value
}

type Join struct {
	On   *Table
	Into *Table
}

func (j *Join) Iterate() (bool, error) {
	// TODO: link thru key
	var err error
	if err = j.On.Next(false); err != nil {
		return false, err
	}
	if err = j.Into.Next(true); err != nil {
		return false, err
	}
	return j.On.Row != nil && j.Into.Row != nil, nil
}

func (j *Join) Get(key string) (v *schema.Value, err error) {
	if v, err = j.On.Get(key); v == nil && err == nil {
		v, err = j.Into.Get(key)
	}
	return v, err
}

func (t *Table) Next(autoCreate bool) (error) {
	t.autoCreate = autoCreate
	if t.exhausted {
		t.Row = nil
		return nil
	}

	// Container
	if ! schema.IsList(t.Corner.State.SelectedMeta()) {
		t.exhausted = true
		t.Row = t.Corner
		return nil
	}

	// List
	meta := t.Corner.State.SelectedMeta().(*schema.List)
	rowNode, err := t.Corner.Node.Next(t.Corner, meta, t.autoCreate, schema.NO_KEYS, t.Row == nil)
	if err != nil {
		return err
	}
	if rowNode == nil {
		t.exhausted = true
		t.Row = nil
	} else {
		t.Row = t.Corner.SelectListItem(rowNode, t.Corner.State.Key())
	}
	return nil
}

func (t *Table) IsDot(r rune) bool {
	return r == '.'
}

func (t *Table) getSelection(path string) (*data.Selection, error) {
	s, found := t.sels[path]
	if found {
		return s, nil
	}
	sel := t.Row
	ident := path
	dot := strings.LastIndexFunc(path, t.IsDot)
	if dot > 0 {
		ident = path[dot:]
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

func (t *Table) Get(key string) (*schema.Value, error) {
	var err error
	v, found := t.vals[key]
	if found {
		return v, nil
	}
	dot := strings.LastIndexFunc(key, t.IsDot)
	sel := t.Row
	ident := key
	if dot > 0 {
		if sel, err = t.getSelection(key[:dot]); err != nil {
			return nil, err
		}
		ident = key[dot:]
	}
	if v, err = data.GetValue(sel, ident); err != nil {
		return nil, err
	}
	if t.vals == nil {
		t.vals = make(map[string]*schema.Value)
	}
	t.vals[key] = v
	return v, nil
}