package c2io
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
		json := "{\"hobbies\":[{\"birding\":{\"favorite-species\":\"towhee\",\"extra\":\"double-mint\"}}]}"
		inIo := strings.NewReader(json)
		var actualBuff bytes.Buffer
		out := NewJsonReceiver(&actualBuff)
		in := JsonTransmitter{in:inIo, out:out, metaRoot:module}
		if err = in.Transmit(); err != nil {
			t.Error("failed to transmit json", err)
		} else {
			actual := string(actualBuff.Bytes())
			t.Log("Round Trip:", actual)
			expected := "{\"hobbies\":[{\"birding\":{\"favorite-species\":\"towhee\"}}]}"
			if actual != expected {
				t.Error(actual, "!=", expected)
			}
		}
	}
}