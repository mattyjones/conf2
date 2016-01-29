package data

import (
	"encoding/json"
	"bytes"
	"strings"
	"io"
	"io/ioutil"
)

type AnyData interface {
	String() (string, error)
	Node() Node
}

type AnyReader struct {
	Reader io.Reader
}

func (any *AnyReader) Node() Node {
	return NewJsonReader(any.Reader).Node()
}

func (any *AnyReader) String() (string, error) {
	b, err := ioutil.ReadAll(any.Reader)
	return string(b), err
}

type AnyJsonString struct {
	Json string
}

func (any *AnyJsonString) Node() Node {
	return NewJsonReader(strings.NewReader(any.Json)).Node()
}

func (any *AnyJsonString) String() (string, error) {
	return any.Json, nil
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

