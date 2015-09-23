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

func ReadKeys(s Selection) (values []*Value, err error) {
fmt.Printf("key - READING kyes\n")
	state := s.WalkState()
	origPosition := state.Position
	list := state.Meta.(*schema.List)
	keyIdents := list.Keys
	values = make([]*Value, len(keyIdents))
	var key *Value
	for i, keyIdent := range keyIdents {
		state.Position = schema.FindByIdent2(list, keyIdent)
		keyMeta := state.Position.(schema.HasDataType)
		if key, err = s.Read(keyMeta); err != nil {
			return nil, err
		}
		if key == nil {
			return nil, errors.New(fmt.Sprint("Key value is nil for ", keyIdent))
		}
		key.Type = keyMeta.GetDataType()
		values[i] = key
	}
	state.Position = origPosition
	return
}
