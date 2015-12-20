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
	v, e := schema.SetValue(meta.GetDataType(), value)
	if e != nil {
		return e
	}
	return n.Write(sel, meta, v)
}

