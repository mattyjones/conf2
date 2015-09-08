package browse
import (
	"testing"
	"yang"
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
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr)); err != nil {
		t.Error("bad module", err)
	} else {
		json := `{"birding":{
"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}],
"reference":{"name":"Peterson's Guide"}
}}`

		tests := [] struct {
			path string
		} {
			//{ "birding" },
			//{ "birding/lifer" },
			{ "birding/lifer=towhee" },
			//{ "birding/reference" },
		}

		for _, test := range tests {
			p, _ := NewPath(test.path)
			inIo := strings.NewReader(json)
			var actualBuff bytes.Buffer
			in, err := NewJsonReader(inIo).GetSelector(module)
			if err != nil {
				t.Error(err)
			}
			if ref, err := WalkPath(in, p); err != nil {
				t.Error(err)
			} else {
				out := NewJsonWriter(&actualBuff)
				to, _ := out.GetSelector()
				fmt.Println("json_test:=====================")
				err = Insert(ref, to)
				if err != nil {
					t.Error("failed to transmit json", err)
				} else {
					actual := string(actualBuff.Bytes())
					t.Log("Round Trip:", actual)
				}
			}
		}
	}
}