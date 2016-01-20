package data

import (
	"schema/yang"
	"strings"
	"testing"
	"schema"
)

type TestMessage struct {
	Message struct {
		Hello string
	}
}

func TestMarshal(t *testing.T) {
	mstr := `
module m {
	prefix "";
	namespace "";
	revision 0;
	container message {
		leaf hello {
			type string;
		}
	}
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	var obj TestMessage
	c := MarshalContainer(&obj)
	r := NewJsonReader(strings.NewReader(`{"message":{"hello":"bob"}}`)).Node()
	if err = NewSelection(m, c).Pull(r).Upsert(); err != nil {
		t.Fatal(err)
	}
	if obj.Message.Hello != "bob" {
		t.Fatal("Not selected")
	}
}

type TestMessageItem struct {
	Id string
}

func TestMarshalIndex(t *testing.T) {
	mstr := `
module m {
	prefix "";
	namespace "";
	revision 0;
	list messages {
		key "id";
		leaf id {
			type string;
		}
	}
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	objs := make(map[string]*TestMessageItem)
	marshaller := &MarshalMap{
		Map: objs,
		OnNewItem: func() interface{} {
			return &TestMessageItem{}
		},
		OnSelectItem: func(item interface{}) Node {
			return MarshalContainer(item)
		},
	}
	d := NewJsonReader(strings.NewReader(`{"messages":[{"id":"bob"},{"id":"barb"}]}`)).Node()
	sel := NewSelection(m, d).Require("messages")
	if err = sel.Push(marshaller.Node()).Upsert(); err != nil {
		t.Fatal(err)
	}
	if objs["bob"].Id != "bob" {
		t.Fatal("Not inserted")
	}
	n := marshaller.Node()
	messagesMeta := m.DataDefs().GetFirstMeta().(*schema.List)
	key := SetValues(messagesMeta.KeyMeta(), "bob")
	foundByKeyNode, nextByKeyErr := n.Next(sel, messagesMeta, false, key, true)
	if nextByKeyErr != nil {
		t.Fatal(nextByKeyErr)
	}
	if foundByKeyNode == nil {
		t.Error("lookup by key failed")
	}

	foundFirstNode, nextFirstErr := n.Next(sel, messagesMeta, false, []*Value{}, true)
	if nextFirstErr != nil {
		t.Fatal(nextFirstErr)
	}
	if foundFirstNode == nil {
		t.Error("lookup by next failed")
	}
}