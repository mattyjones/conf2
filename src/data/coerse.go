package data

import (
	"fmt"
	"errors"
	"schema"
)

var NO_KEYS = make([]*Value, 0)

func CoerseKeys(list *schema.List, keyStrs []string) ([]*Value, error) {
	var err error
	if len(keyStrs) == 0 {
		return NO_KEYS, nil
	}
	if len(list.Keys) != len(keyStrs) {
		return NO_KEYS, errors.New("Missing keys on " + list.GetIdent())
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