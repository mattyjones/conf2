package browse
import (
	"schema"
	"fmt"
)

var NO_KEYS = make([]Value, 0)

func CoerseKeys(list *schema.List, keyStrs []string) ([]Value, error) {
	var err error
	if len(keyStrs) == 0 {
		return NO_KEYS, nil
	}
	if len(list.Keys) != len(keyStrs) {
		msg := fmt.Sprintf("Missing keys on %s", list.GetIdent())
		return NO_KEYS, &browseError{Msg:msg}
	}
	values := make([]Value, len(keyStrs))
	for i, keyStr := range keyStrs {
		keyProp := schema.FindByIdent2(list, list.Keys[i])
		values[i].Type = keyProp.(schema.HasDataType).GetDataType()
		err = values[i].CoerseStrValue(keyStr)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func ReadKeys(s Selection) (values []Value, err error) {
	state := s.WalkState()
	origPosition := state.Position
	list := state.Meta.(*schema.List)
	keys := list.Keys
	values = make([]Value, len(keys))
	for i, key := range keys {
		state.Position = schema.FindByIdent2(list, key)
		if err = s.Read(&values[i]); err != nil {
			return
		}
	}
	state.Position = origPosition
	return
}
