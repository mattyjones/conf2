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
	json := `{"a":{"b":"hi","c":"bye"}}`
	var err error
	var m1 *schema.Module
	var m2 *schema.Module
	if m1, err = yang.LoadModuleFromByteArray([]byte(m1Str), nil); err != nil {
		t.Error(err)
	} else if m2, err = yang.LoadModuleFromByteArray([]byte(m2Str), nil); err != nil {
		t.Error(err)
	} else {
		mapp := NewMap()
		mapp.Mappings["a"] = Mapping{To:"x"}
		mapp.Mappings["b"] = Mapping{To:"y"}
		jsonRdr := browse.NewJsonReader(strings.NewReader(json))
		var actualBuff bytes.Buffer
		jsonWtr := browse.NewJsonWriter(&actualBuff)
		var from browse.Selection
		from, err = jsonRdr.GetSelector(m1)
		if err != nil {
			t.Error(err)
		} else {
			toJson, _ := jsonWtr.GetSelector()
			toJson.WalkState().Meta = m2
			b := &Bridge{Map: mapp}
			to, _ := b.enterBridge(toJson)
			browse.Insert(from, to, browse.NewExhaustiveController())
			actual := string(actualBuff.Bytes())
			expected := `{"x":{"y":"hi","c":"bye"}}`
			if actual != expected {
				t.Errorf("\nExpected:\"%s\"\n  Actual:\"%s\"", expected, actual)
			}
		}
	}
}

var m1Str = `
module m1 {
	prefix "p";
	namespace "n";
	description "d";
	revision 0000-00-00 {
		description "d";
	}
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

var m2Str = `
module m2 {
	prefix "p";
	namespace "n";
	revision 0000-00-00 {
		description "d";
	}
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
