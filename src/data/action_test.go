package data

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
	store.Actions["sayHello"] = func(state *Selection, meta *schema.Rpc, input *Selection) (output Node, err error) {
		if err = input.Selector().Push(b.Node()).Insert().LastErr; err != nil {
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str: fmt.Sprint("Hello ", yourName)}
		return b.Container(""), nil
	}
	in := NewJsonReader(strings.NewReader(`{"name":"joe"}`)).Node()
	var actual bytes.Buffer
	sel := b.Select().Find("sayHello")
	if sel.LastErr != nil {
		t.Fatal(sel.LastErr)
	}
	actionOut, actionErr := sel.Selection.Action(in)
	if actionErr != nil {
		t.Fatal(actionErr)
	}
	if err = actionOut.Selector().Push(NewJsonWriter(&actual).Node()).Insert().LastErr; err != nil {
		t.Fatal(err)
	}
	AssertStrEqual(t, "joe", yourName.Str)
	AssertStrEqual(t, `{"salutation":"Hello joe"}`, actual.String())
}
