package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"schema"
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
        leaf salutaion {
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
	var saidHello bool
	store.Actions["sayHello"] = func(state *WalkState, meta *schema.Rpc) (input Selection, output Selection, err error) {
		
		saidHello = true
		return nil, nil, nil

	}
	in := NewJsonFragmentReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonFragmentWriter(&actual)
	b := NewStoreBrowser(m, store)
	err = Action(NewPath("sayHello"), b, in, out)
	if err != nil {
		t.Fatal(err)
	}
	if ! saidHello {
		t.Error("Never said hello")
	}
}