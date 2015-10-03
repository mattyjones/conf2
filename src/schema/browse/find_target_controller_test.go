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
			p := NewPath(test.path)
			inIo := strings.NewReader(json)
			var actualBuff bytes.Buffer
			state := NewWalkState(module)
			in, err := NewJsonReader(inIo).GetSelector(state)
			if err != nil {
				t.Error(err)
			}
			if ref, walkedState, err := WalkPath(state, in, p); err != nil {
				t.Error(err)
			} else {
				out := NewJsonWriter(&actualBuff)
				to, _ := out.GetSelector()
				cntlr := NewFullWalkFromPath(p)
				err = Upsert(walkedState, ref, to, cntlr)
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
