package data

import (
	"encoding/json"
	"schema"
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

type AnyNode struct {
	Any  Node
	Meta schema.MetaList
}

func (any *AnyNode) Node() Node {
	return any.Any
}

func (any *AnyNode) String() (string, error) {
	var out bytes.Buffer
	w := NewJsonWriter(&out)
	err := NodeToNode(any.Any, w.Node(), any.Meta).Insert()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

