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
	store.Actions["sayHello"] = func(state *Selection, meta *schema.Rpc, input *Selection) (output *Selection, err error) {
		readInput := NewSelection(b.Container(""), meta.Input)
		if err = Insert(input, readInput); err != nil {
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str: fmt.Sprint("Hello ", yourName)}
		return NewSelection(b.Container(""), meta.Output), nil
	}
	var actionSel, rpcOut *Selection
	if actionSel, err = b.Selector(NewPath("sayHello")); err != nil {
		t.Error(err)
	}
	rpc := actionSel.State.Position().(*schema.Rpc)
	var readJson Node
	if readJson, err = NewJsonReader(strings.NewReader(`{"name":"joe"}`)).Node(); err != nil {
		t.Fatal(err)
	}
	in := NewSelection(readJson, rpc.Input)
	if rpcOut, err = Action(actionSel, in); err != nil {
		t.Error(err)
	}
	var actual bytes.Buffer
	out := NewJsonWriter(&actual).Selector(rpcOut.State)
	if err = Insert(rpcOut, out); err != nil {
		t.Error(err)
	}
	if yourName.Str != "joe" {
		t.Error("Your name ", yourName)
	}
	if actual.String() != `{"salutation":"Hello joe"}` {
		t.Error(actual.String())
	}
}
