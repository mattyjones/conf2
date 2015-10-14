package browse
import (
	"schema"
	"fmt"
	"errors"
)

var NO_KEYS = make([]*Value, 0)

func CoerseKeys(list *schema.List, keyStrs []string) ([]*Value, error) {
	var err error
	if len(keyStrs) == 0 {
		return NO_KEYS, nil
	}
	if len(list.Keys) != len(keyStrs) {
		msg := fmt.Sprintf("Missing keys on %s", list.GetIdent())
		return NO_KEYS, &browseError{Msg:msg}
	}
	values := make([]*Value, len(keyStrs))
	for i, keyStr := range keyStrs {
		keyProp := schema.FindByIdent2(list, list.Keys[i])
		values[i] = &Value {
			Type : keyProp.(schema.HasDataType).GetDataType(),
		}
		err = values[i].CoerseStrValue(keyStr)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func ReadKeys(state *WalkState, s Selection) (values []*Value, err error) {
	if len(state.Key()) > 0 {
		return state.Key(), nil
	}
	list := state.SelectedMeta().(*schema.List)
	values = make([]*Value, len(list.Keys))
	var key *Value
	for i, keyIdent := range list.Keys {
		keyMeta := schema.FindByIdent2(state.SelectedMeta(), keyIdent).(schema.HasDataType)
		if key, err = s.Read(state, keyMeta); err != nil {
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
