package data
import (
	"schema"
	"errors"
)

func ChangeValue(sel *Selection, ident string, value interface{}) error {
	n := sel.Node
	if cw, ok := n.(ChangeAwareNode); ok {
		n = cw.Changes()
	}
	pos := schema.FindByIdent2(sel.State.SelectedMeta(), ident)
	if pos == nil {
		return errors.New("property not found " + ident)
	}
	meta := pos.(schema.HasDataType)
	v, e := SetValue(meta.GetDataType(), value)
	if e != nil {
		return e
	}
	return n.Write(sel, meta, v)
}

func Get(sel *Selection, ident string) (interface{}, error) {
	prop := schema.FindByIdent2(sel.State.SelectedMeta(), ident)
	if prop != nil {
		if schema.IsLeaf(prop) {
			v, err := sel.Node.Read(sel, prop.(schema.HasDataType))
			if err != nil {
				return nil, err
			}
			return v.Value(), nil
		} else {
			return sel.Peek(), nil
		}
	}
	return nil, nil
}

func GetValue(sel *Selection, ident string) (*Value, error) {
	prop := schema.FindByIdent2(sel.State.SelectedMeta(), ident)
	if prop != nil {
		v, err := sel.Node.Read(sel, prop.(schema.HasDataType))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, nil
}

func ClearAll(sel *Selection) error {
	return sel.Node.Event(sel, DELETE)
}

func SelectMetaList(sel *Selection, ident string, autoCreate bool) (*Selection, error) {
	m := schema.FindByIdent2(sel.State.SelectedMeta(), ident)
	var err error
	var child Node
	if m != nil {
		sel.State.SetPosition(m)
		child, err = sel.Node.Select(sel, m.(schema.MetaList), false)
		if err != nil {
			return nil, err
		} else if child == nil && autoCreate {
			child, err = sel.Node.Select(sel, m.(schema.MetaList), true)
			if err != nil {
				return nil, err
			}
		}
	}
	if child != nil {
		return sel.Select(child), nil
	}
	return nil, nil
}

