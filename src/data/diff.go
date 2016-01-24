package data

import (
	"schema"
	"conf2"
)

func Diff(a Node, b Node) Node {
	n := &MyNode{}
	n.OnSelect = func(sel *Selection, meta schema.MetaList, new bool) (n Node, err error) {
conf2.Debug.Printf("OnSelect %s", meta.GetIdent())
		var aNode, bNode Node
		if aNode, err = a.Select(sel, meta, false); err != nil {
			return nil, err
		}
		if bNode, err = b.Select(sel, meta, false); err != nil {
			return nil, err
		}
		if aNode == nil {
			return nil, nil
		}
		if bNode == nil {
			return aNode, nil
		}
		return Diff(aNode, bNode), nil
	}
	n.OnRead = func(sel *Selection, meta schema.HasDataType) (changedValue *Value, err error) {
conf2.Debug.Printf("OnRead %s", meta.GetIdent())
		var aVal, bVal *Value
		if aVal, err = a.Read(sel, meta); err != nil {
			return nil, err
		}
		if bVal, err = b.Read(sel, meta); err != nil {
			return nil, err
		}
		if aVal == nil {
			if bVal == nil {
				return nil, nil
			}
			return bVal, nil
		}
		if aVal.Equal(bVal) {
			return nil, nil
		}
		return aVal, nil
	}
	return n
}
