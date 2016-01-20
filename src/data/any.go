package data

import (
	"encoding/json"
	"bytes"
)

type AnyData interface {
	String() (string, error)
	Node() Node
}

type AnyJson struct {
	container map[string]interface{}
}

func (any *AnyJson) Node() Node {
	return JsonContainerReader(any.container)
}

func (any *AnyJson) String() (string, error) {
	bytes, err := json.Marshal(any.container)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type AnySelection struct {
	Selection *Selection
}

func (any *AnySelection) Node() Node {
	return any.Selection.Node()
}

func (any *AnySelection) String() (string, error) {
	var out bytes.Buffer
	if err := any.Selection.Push(NewJsonWriter(&out).Node()).Insert(); err != nil {
		return "", err
	}
	return out.String(), nil
}

