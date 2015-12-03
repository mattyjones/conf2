package data

import (
	"schema"
)

func Diff(a Node, b Node) Node {
	n := &MyNode{}
	n.OnSelect = func(state *Selection, meta schema.MetaList, new bool) (n Node, err error) {
		var aNode, bNode Node
		if aNode, err = a.Select(state, meta, false); err != nil {
			return nil, err
		}
		if bNode, err = b.Select(state, meta, false); err != nil {
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
	n.OnRead = func(state *Selection, meta schema.HasDataType) (changedValue *schema.Value, err error) {
		var aVal, bVal *schema.Value
		if aVal, err = a.Read(state, meta); err != nil {
			return nil, err
		}
		if bVal, err = b.Read(state, meta); err != nil {
			return nil, err
		}
		if aVal.Equal(bVal) {
			return nil, nil
		}
		return aVal, nil
	}
	return n
}
