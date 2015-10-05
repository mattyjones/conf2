package adapt
import (
	"testing"
	"schema"
	"schema/browse"
	"schema/yang"
	"strings"
	"bytes"
)

func TestBridge(t *testing.T) {
	externalData := `{"a":{"b":"hi","c":"bye"}}`
	var err error
	var externalModule *schema.Module
	var internalModule *schema.Module
	if externalModule, err = yang.LoadModuleFromByteArray([]byte(externalYang), nil); err != nil {
		t.Error(err)
	} else if internalModule, err = yang.LoadModuleFromByteArray([]byte(internalYang), nil); err != nil {
		t.Error(err)
	} else {
		var actualBuff bytes.Buffer
		internal := browse.NewJsonWriter(&actualBuff, internalModule)
		b := NewBridge(internal, externalModule)
		a := b.Mapping.AddMapping("a", "x")
		a.AddMapping("b", "y")
		input := browse.NewJsonReader(strings.NewReader(externalData), externalModule)
		err = browse.Upsert(browse.NewPath(""), input, b)
		if err != nil {
			t.Error(err)
		} else {
			actual := string(actualBuff.Bytes())
			expected := `{"x":{"y":"hi","c":"bye"}}`
			if actual != expected {
				t.Errorf("\nExpected:\"%s\"\n  Actual:\"%s\"", expected, actual)
			}
		}
	}
}

var externalYang = `
module external {
	prefix "";
	namespace "";
	revision 0;
	container a {
		leaf b {
			type string;
		}
		leaf c {
			type string;
		}
	}
}
`

var internalYang = `
module internal {
	prefix "";
	namespace "";
	revision 0;
	container x {
		leaf y {
			type string;
		}
		leaf c {
			type string;
		}
	}
}
`
