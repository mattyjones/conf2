package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
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
	in := NewJsonFragmentReader(strings.NewReader(`{"name":"joe"}`))
	var actual bytes.Buffer
	out := NewJsonFragmentWriter(&actual)
	err = Action(NewPath("sayHello"), b, in, out)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(actual.Bytes()))
}