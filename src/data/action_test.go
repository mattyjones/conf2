package data

import (
	"bytes"
	"fmt"
	"schema"
	"schema/yang"
	"strings"
	"testing"
	"conf2"
)

func TestAction(t *testing.T) {
	y := `
module m {
	prefix "";
	namespace "";
	revision 0000-00-00 {
	  description "";
    }
    rpc sayHello {
      input {
        leaf name {
          type string;
        }
      }
      output {
        leaf salutation {
          type string;
        }
      }
    }
}`
	m, err := yang.LoadModuleFromByteArray([]byte(y), nil)
	if err != nil {
		t.Fatal(err)
	}
	// lazy trick, we stick all data, input, output into one bucket
	store := NewBufferStore()
	b := NewStoreData(m, store)
	var yourName *Value
	store.Actions["sayHello"] = func(state *Selection, meta *schema.Rpc, input Node) (output Node, err error) {
		if err = NodeToNode(input, b.Container(""), meta.Input).Insert(); err != nil {
conf2.Debug.Printf("here")
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str: fmt.Sprint("Hello ", yourName)}
		return b.Container(""), nil
	}
	in := NewJsonReader(strings.NewReader(`{"name":"joe"}`)).Node()
	var actual bytes.Buffer
	out := NewJsonWriter(&actual).Node()
	if err = PathAction(b, "sayHello", in, out); err != nil {
		t.Fatal(err)
	}
	if yourName.Str != "joe" {
		t.Error("Your name ", yourName)
	}
	if actual.String() != `{"salutation":"Hello joe"}` {
		t.Error(actual.String())
	}
}
