package browse

import (
"schema"
"strings"
)

type Data interface {
	Selector(path *Path) (*Selection, error)
	Schema() schema.MetaList
}

//
// Example:
//  DataValue(data, "foo=10.bar.blah.0")
//
func DataValue(data Data, path string) (interface{}, error) {
	var sel *Selection
	var p *Path
	var propIdent string
	valSep := strings.LastIndex(path, "/")
	if valSep > 0 {
		p = NewPath(path[:valSep])
		propIdent = path[valSep:]
	} else {
		p = NewPath("")
		propIdent = path
	}
	sel, err := data.Selector(p)
	if err != nil {
		return nil, err
	}
	propMeta := schema.FindByIdent2(sel.State.SelectedMeta(), propIdent)
	var val *Value
	val, err = sel.Node.Read(sel, propMeta.(schema.HasDataType))
	if err != nil || val == nil {
		return nil, err
	}
	return val.Value(), nil
}
