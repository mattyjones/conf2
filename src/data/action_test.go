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
		if err = input.Push(b.Node()).Insert(); err != nil {
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str: fmt.Sprint("Hello ", yourName)}
		return b.Container(""), nil
	}
	in := NewJsonReader(strings.NewReader(`{"name":"joe"}`)).Node()
	var actual bytes.Buffer
	sel, err := b.Select().Find("sayHello")
	if err != nil {
		t.Fatal(err)
	}
	actionOut, actionErr := sel.Action(in)
	if actionErr != nil {
		t.Fatal(actionErr)
	}
	actionOut.Push(NewJsonWriter(&actual).Node()).Insert()
	AssertStrEqual(t, "joe", yourName.Str)
	AssertStrEqual(t, `{"salutation":"Hello joe"}`, actual.String())
}
