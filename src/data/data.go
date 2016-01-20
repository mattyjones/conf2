package data

type Data interface {
	Select() *Selection
}

//
// Example:
//  DataValue(data, "foo=10/bar/blah.0")
//
//func DataValue(data Data, path string) (interface{}, error) {
//	var sel *Selection
//	var p *PathSlice
//	var propIdent string
//	valSep := strings.LastIndex(path, "/")
//	if valSep > 0 {
//		p = NewPathSlice(path[:valSep], data.Schema())
//		propIdent = path[valSep:]
//	} else {
//		p = NewPathSlice("", data.Schema())
//		propIdent = path
//	}
//	var err error
//	sel, err = WalkPath(NewSelection(data.Node(), data.Schema()), p)
//	if err != nil {
//		return nil, err
//	}
//	propMeta := schema.FindByIdent2(sel.State.SelectedMeta(), propIdent)
//	var val *Value
//	val, err = sel.Node.Read(sel, propMeta.(schema.HasDataType))
//	if err != nil || val == nil {
//		return nil, err
//	}
//	return val.Value(), nil
//}
