package data
import (
	"testing"
	"schema"
)

func TestExtend(t *testing.T) {
	child := &MyNode{Label: "Bloop"}
	n := &MyNode{
		Label: "N",
		OnRead: func(*Selection, schema.HasDataType) (*schema.Value, error) {
			return &schema.Value{Str:"Hello"}, nil
		},
		OnSelect: func(s *Selection, meta schema.MetaList, new bool) (Node, error) {
			return child, nil
		},
	}
	x := Extend{
		Label: "Bleep",
		Node: n,
		OnRead: func(p Node, s *Selection, m schema.HasDataType) (*schema.Value, error) {
			v, _ := p.Read(s, m)
			return &schema.Value{Str:v.Str + " World"}, nil
		},
	}
	actualValue, _  := x.Read(nil, nil)
	if actualValue.Str != "Hello World" {
		t.Error(actualValue.Str)
	}
	actualChild, _ := x.Select(nil, nil, false)
	if actualChild.String() != "(Bloop) <- Bleep" {
		t.Error(actualChild.String())
	}
}