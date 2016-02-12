package data
import (
	"testing"
	"schema"
)

func TestExtend(t *testing.T) {
	child := &MyNode{Label: "Bloop"}
	n := &MyNode{
		Label: "Blop",
		OnRead: func(*Selection, schema.HasDataType) (*Value, error) {
			return &Value{Str:"Hello"}, nil
		},
		OnSelect: func(s *Selection, r ContainerRequest) (Node, error) {
			return child, nil
		},
	}
	x := Extend{
		Label: "Bleep",
		Node: n,
		OnRead: func(p Node, s *Selection, m schema.HasDataType) (*Value, error) {
			v, _ := p.Read(s, m)
			return &Value{Str:v.Str + " World"}, nil
		},
	}
	actualValue, _  := x.Read(nil, nil)
	if actualValue.Str != "Hello World" {
		t.Error(actualValue.Str)
	}
	if x.String() != "(Blop) <- Bleep" {
		t.Error(x.String())
	}
	actualChild, _ := x.Select(nil, ContainerRequest{})
	if actualChild.String() != "Bloop" {
		t.Error(actualChild.String())
	}
}