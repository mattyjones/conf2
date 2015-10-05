package browse
import (
	"testing"
	"schema/yang"
	"bytes"
)


func TestJsonWriterListInList(t *testing.T) {
	moduleStr := `
module m {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	list l1 {
		list l2 {
		    key "a";
			leaf a {
				type string;
			}
			leaf b {
			    type string;
			}
		}
	}
}
	`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Error("bad module", err)
	} else {
		b := NewBucketBrowser(module)
		l1 := make([]map[string]interface{}, 1)
		l1[0] = make(map[string]interface{}, 1)
		l2 := make([]map[string]interface{}, 1)
		l2[0] = make(map[string]interface{}, 1)
		l2[0]["a"] = "hi"
		l2[0]["b"] = "bye"
		l1[0]["l2"] = l2
		b.Bucket["l1"] = l1
		var json bytes.Buffer
		w := NewJsonFragmentWriter(&json)
		err = Upsert(NewPath(""), b, w)
		if err != nil {
			t.Fatal(err)
		}
		actual := string(json.Bytes())
		expected := `{"l1":[{"l2":[{"a":"hi","b":"bye"}]}]}`
		if actual != expected {
			t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
		}
	}
}
