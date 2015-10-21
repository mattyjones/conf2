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
	store := NewBufferStore()
	var yourName *Value
	store.Actions["sayHello"] = func(state *WalkState, meta *schema.Rpc, input Selection) (output Selection, outputState *WalkState, err error) {
		var read Selection
		var inputState *WalkState
		b := NewStoreBrowser(meta.Input, store)
		if read, inputState, err = b.Selector(NewPath(""), INSERT); err != nil {
			return nil, nil, err
		}
		if err = Edit(inputState, input, read, INSERT, FullWalk()); err != nil {
			return nil, nil, err
		}
		yourName = store.Values["name"]
		store.Values["salutation"] = &Value{Str:fmt.Sprint("Hello ", yourName)}
		b = NewStoreBrowser(meta.Output, store)
		return b.Selector(NewPath(""), READ)
	}
	in := NewJsonFragmentReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonFragmentWriter(&actual)
	b := NewStoreBrowser(m, store)
	err = Action(NewPath("sayHello"), b, in, out)
	if err != nil {
		t.Fatal(err)
	}
	if yourName.Str != "joe" {
		t.Error("Your name ", yourName)
	}
	if actual.String() != `{"salutation":"Hello joe"}` {
		t.Error(actual.String())
	}
}