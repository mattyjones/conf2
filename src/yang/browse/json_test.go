package browse

import (
	"testing"
	"yang"
	"bytes"
	"strings"
)

func TestJson(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	list hobbies {
		container birding {
			leaf favorite-species {
				type string;
			}
		}
		container hockey {
			leaf favorite-team {
				type string;
			}
		}
	}
}
	`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr)); err != nil {
		t.Error("bad module", err)
	} else {
		json := `{"hobbies":[
{"birding":{"favorite-species":"towhee","extra":"double-mint"}},
{"birding":{"favorite-species":"robin"}}
]}`
		inIo := strings.NewReader(json)
		var actualBuff bytes.Buffer
		out := NewJsonWriter(&actualBuff)
		to, _ := out.GetSelector()
		if err != nil {
			t.Error(err)
		}
		in, err := NewJsonReader(inIo).GetSelector(module)
		if err != nil {
			t.Error(err)
		}
		err = Insert(in, to)
		if err != nil {
			t.Error("failed to transmit json", err)
		} else {
			actual := string(actualBuff.Bytes())
			t.Log("Round Trip:", actual)
			expected := strings.Replace(`{"hobbies":[
{"birding":{"favorite-species":"towhee"}},
{"birding":{"favorite-species":"robin"}}
]}`, "\n", "", -1)
			if actual != expected {
				t.Error(actual, "!=", expected)
			}
		}
	}
}