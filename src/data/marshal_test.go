package data

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
	c := MarshalContainer(&obj)
	r := NewJsonReader(strings.NewReader(`{"message":{"hello":"bob"}}`)).Node()
	if err = NodeToNode(r, c, m).Upsert(); err != nil {
		t.Fatal(err)
	}
	if obj.Message.Hello != "bob" {
		t.Fatal("Not selected")
	}
}
