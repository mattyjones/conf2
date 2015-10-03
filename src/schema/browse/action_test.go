package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"schema"
)

func DISABLE_TestAction(t *testing.T) {
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
	b := NewBucketBrowser(m)
	s, state, _ := b.RootSelector()
	target, state, err := WalkPath(state, s, NewPath("sayHello"))
	in := NewJsonReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonWriter(&actual)
	rpc := state.SelectedMeta().(*schema.Rpc)
	var outSel, inSel Selection
	rpcState := NewWalkState(rpc.Input)
	inSel, _ = in.GetSelector(rpcState)
	outSel, _ = out.GetSelector()
	err = Action(rpcState, target, inSel, outSel)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(actual.Bytes()))
}