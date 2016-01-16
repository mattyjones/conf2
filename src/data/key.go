package data

import (
	"errors"
	"fmt"
	"schema"
)

func ReadKeys(sel *Selection) (values []*Value, err error) {
	if len(sel.State.Key()) > 0 {
		return sel.State.Key(), nil
	}
	list := sel.State.SelectedMeta().(*schema.List)
	values = make([]*Value, len(list.Keys))
	var key *Value
	for i, keyIdent := range list.Keys {
		keyMeta := schema.FindByIdent2(sel.State.SelectedMeta(), keyIdent).(schema.HasDataType)
		if key, err = sel.Node.Read(sel, keyMeta); err != nil {
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

