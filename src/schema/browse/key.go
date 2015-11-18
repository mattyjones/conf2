package browse

import (
	"errors"
	"fmt"
	"schema"
)

var NO_KEYS = make([]*Value, 0)

func CoerseKeys(list *schema.List, keyStrs []string) ([]*Value, error) {
	var err error
	if len(keyStrs) == 0 {
		return NO_KEYS, nil
	}
	if len(list.Keys) != len(keyStrs) {
		msg := fmt.Sprintf("Missing keys on %s", list.GetIdent())
		return NO_KEYS, &browseError{Msg: msg}
	}
	values := make([]*Value, len(keyStrs))
	for i, keyStr := range keyStrs {
		keyProp := schema.FindByIdent2(list, list.Keys[i])
		if keyProp == nil {
			return nil, errors.New(fmt.Sprintf("no key prop %s on %s", list.Keys[i], list.GetIdent()))
		}
		values[i] = &Value{
			Type: keyProp.(schema.HasDataType).GetDataType(),
		}
		err = values[i].CoerseStrValue(keyStr)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

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
