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
	if err = Select(m, c).Selector().Pull(r).Upsert().LastErr; err != nil {
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
	sel := Select(m, d).Find("messages")
	if err = sel.Push(marshaller.Node()).Upsert().LastErr; err != nil {
		t.Fatal(err)
	}
	if objs["bob"].Id != "bob" {
		t.Fatal("Not inserted")
	}
	n := marshaller.Node()
	r := ListRequest{
		Meta: m.DataDefs().GetFirstMeta().(*schema.List),
		First: true,
	}
	r.Key = SetValues(r.Meta.KeyMeta(), "bob")
	foundByKeyNode, _, nextByKeyErr := n.Next(sel.Selection, r)
	if nextByKeyErr != nil {
		t.Fatal(nextByKeyErr)
	}
	if foundByKeyNode == nil {
		t.Error("lookup by key failed")
	}
	r.Key = []*Value{}
	foundFirstNode, _, nextFirstErr := n.Next(sel.Selection, r)
	if nextFirstErr != nil {
		t.Fatal(nextFirstErr)
	}
	if foundFirstNode == nil {
		t.Error("lookup by next failed")
	}
}