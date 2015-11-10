package browse

import (
	"bytes"
	"fmt"
	"schema"
	"schema/yang"
	"strings"
	"testing"
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
	store.Actions["sayHello"] = func(state *Selection, meta *schema.Rpc, input Node) (output *Selection, err error) {
		readInput := NewSelection(b.Container(""), meta.Input)
		if err = InsertByNode(readInput, input, readInput.Node()); err != nil {
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str: fmt.Sprint("Hello ", yourName)}
		return NewSelection(b.Container(""), meta.Output), nil
	}
	var in Node
	if in, err = NewJsonReader(strings.NewReader(`{"name":"joe"}`)).Node(); err != nil {
		t.Fatal(err)
	}
	var actual bytes.Buffer
	out := NewJsonWriter(&actual)
	var actionSel, rpcOut *Selection
	if actionSel, err = b.Selector(NewPath("sayHello")); err != nil {
		t.Error(err)
	}
	if rpcOut, err = Action(actionSel, in); err != nil {
		t.Error(err)
	}
	if err = Insert(rpcOut, out.Selector(rpcOut)); err != nil {
		t.Error(err)
	}
	if yourName.Str != "joe" {
		t.Error("Your name ", yourName)
	}
	if actual.String() != `{"salutation":"Hello joe"}` {
		t.Error(actual.String())
	}
}
