package data

import (
	"errors"
	"fmt"
	"schema"
)

func ReadKeys(sel *Selection) (values []*Value, err error) {
	if len(sel.path.key) > 0 {
		return sel.path.key, nil
	}
	list := sel.path.meta.(*schema.List)
	values = make([]*Value, len(list.Key))
	var key *Value
	for i, keyIdent := range list.Key {
		keyMeta := schema.FindByIdent2(sel.path.meta, keyIdent).(schema.HasDataType)
		if key, err = sel.node.Read(sel, keyMeta); err != nil {
			return nil, err
		}
		if key == nil {
			return nil, errors.New(fmt.Sprint("Key value is nil for ", keyIdent))
		}
		key.Type = keyMeta.GetDataType()
		values[i] = key
	}
	return
}

