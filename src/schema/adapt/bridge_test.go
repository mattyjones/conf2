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
		mapping := NewBridgeMapping("")
		a := mapping.AddMapping("a", "x")
		a.AddMapping("b", "y")
		external := browse.NewJsonReader(strings.NewReader(externalData))
		var actualBuff bytes.Buffer
		internal := browse.NewJsonWriter(&actualBuff)
		var from browse.Selection
		externalState := browse.NewWalkState(externalModule)
		from, err = external.GetSelector(externalState)
		if err != nil {
			t.Error(err)
		} else {
			toJson, _ := internal.GetSelector()
			b := &Bridge{}
			internalState := browse.NewWalkState(internalModule)
			to, _ := b.selectBridge(toJson, internalState, mapping)
			err = browse.Upsert(externalState, from, to, browse.WalkAll())
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
