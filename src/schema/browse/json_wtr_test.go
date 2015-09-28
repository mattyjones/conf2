package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"fmt"
)


func TestJsonPath(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	container birding {
		list lifer {
		    key "species";
			leaf species {
				type string;
			}
			leaf location {
			    type string;
			}
		}
		container reference {
		    leaf name {
		        type string;
		    }
		}
	}
}
	`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Error("bad module", err)
	} else {
		json := `{"birding":{
"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}],
"reference":{"name":"Peterson's Guide"}
}}`

		tests := [] struct {
			path string
			expected string
		} {
			{ "", strings.Replace(json, "\n", "", -1)},
			{ "birding", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}],"reference":{"name":"Peterson's Guide"}}`},
			{ "birding/lifer=towhee", `{"species":"towhee","location":"Hammonasset, CT"}` },
			{ "birding?depth=1", `{"lifer":[],"reference":{}}` },
			{ "birding/lifer", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}]}` },
			{ "birding/lifer?depth=1", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}]}` },
			{ "birding/reference", `{"name":"Peterson's Guide"}` },
		}

		for i, test := range tests {
			p, _ := ParsePath(test.path)
			inIo := strings.NewReader(json)
			var actualBuff bytes.Buffer
			in, err := NewJsonReader(inIo).GetSelector(module, false)
			if err != nil {
				t.Error(err)
			}
			if ref, err := WalkPath(in, p); err != nil {
				t.Error(err)
			} else {
				out := NewJsonWriter(&actualBuff)
				to, _ := out.GetSelector()
				var cntlr WalkController
				if cntlr, err = p.WalkTargetController(); err != nil {
					t.Error(err)
				} else {
					err = Upsert(ref, to, cntlr)
					if err != nil {
						t.Error("failed to transmit json", err)
					} else {
						actual := string(actualBuff.Bytes())
						if actual != test.expected {
							msg := fmt.Sprintf("Failed subtest #%d - '%s'\nExpected:'%s'\n  Actual:'%s'",
								i + 1, test.path, test.expected, actual)
							t.Error(msg)
						}
					}
				}
			}
		}
	}
}

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
		w := NewJsonWriter(&json)
		s, _ := w.GetSelector()
		r, _ := b.RootSelector()
		err = Upsert(r, s, NewExhaustiveController())
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
