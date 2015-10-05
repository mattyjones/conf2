package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"fmt"
)

func TestFindTargetController(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "";
	namespace "";
	revision 0;
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
			inIo := strings.NewReader(json)
			in := &JsonReader{In:inIo, Meta:module}
			if err != nil {
				t.Error(err)
			}
			var actualBuff bytes.Buffer
			out := NewJsonFragmentWriter(&actualBuff)
			err = Upsert(NewPath(test.path), in, out)
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
