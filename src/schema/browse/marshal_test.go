package browse

import (
	"schema/yang"
	"strings"
	"testing"
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
	c := &MarshalContainer{Obj: &obj}
	var r Node
	r, err = NewJsonReader(strings.NewReader(`{"message":{"hello":"bob"}}`)).Node()
	if err != nil {
		t.Fatal(err)
	}
	sel := NewSelection(c, m)
	err = UpsertByNode(sel, r, sel.Node())
	if err != nil {
		t.Fatal(err)
	}
	if obj.Message.Hello != "bob" {
		t.Fatal("Not selected")
	}
}
