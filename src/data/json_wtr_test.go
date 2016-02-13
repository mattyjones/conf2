package data

import (
	"bytes"
	"schema/yang"
	"testing"
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
	if m, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Error("bad module", err)
	} else {
		root := map[string]interface{}{
			"l1": []map[string]interface{}{
				map[string]interface{}{"l2" : []map[string]interface{}{
					map[string]interface{}{
							"a" : "hi",
							"b" : "bye",
						},
					},
				},
			},
		}
		b := BucketData{Meta: m, Root: root}
		var json bytes.Buffer
		if err = b.Select().Selector().Push(NewJsonWriter(&json).Node()).Upsert().LastErr; err != nil {
			t.Fatal(err)
		}
		actual := string(json.Bytes())
		expected := `{"l1":[{"l2":[{"a":"hi","b":"bye"}]}]}`
		if actual != expected {
			t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
		}
	}
}
