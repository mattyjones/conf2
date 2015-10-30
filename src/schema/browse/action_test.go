package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"schema"
	"fmt"
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
	var yourName *Value
	store.Actions["sayHello"] = func(state *Selection, meta *schema.Rpc, input *Selection) (output *Selection, err error) {
		var read *Selection
		b := NewStoreBrowser(meta.Input, store)
		if read, err = b.Selector(NewPath("")); err != nil {
			return nil, err
		}
		if err = Insert(input, read); err != nil {
			return nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str:fmt.Sprint("Hello ", yourName)}
		b = NewStoreBrowser(meta.Output, store)
		return b.Selector(NewPath(""))
	}
	in := NewJsonReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonWriter(&actual)
	b := NewStoreBrowser(m, store)
	var actionSel, rpcOut, rpcIn *Selection
	if actionSel, err = b.Selector(NewPath("sayHello")); err != nil {
		t.Error(err)
	}
	rpc := actionSel.Position().(*schema.Rpc)
	if rpcIn, err = in.Selector(rpc.Input); err != nil {
		t.Error(err)
	}
	if rpcOut, err = Action(actionSel, rpcIn); err != nil {
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