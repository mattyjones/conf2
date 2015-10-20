package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"schema"
	"errors"
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
	store.Actions["sayHello"] = func(state *WalkState, meta *schema.Rpc) (input Selection, output Selection, err error) {
		panic("YEAH")
		return nil, nil, errors.New("Got here!")

	}
	in := NewJsonFragmentReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonFragmentWriter(&actual)
	rpc := schema.FindByIdent2(m.DataDefs(), "sayHello")
	if rpc == nil {
		t.Error("No rpc")
	}
	b := NewStoreBrowser(m, store)
	err = Action(NewPath("sayHello"), b, in, out)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(actual.Bytes()))
}