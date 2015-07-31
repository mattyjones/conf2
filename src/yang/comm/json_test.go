package comm
import (
	"testing"
	"yang"
	"yang/browse"
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
		rcvr := NewJsonReceiver(&actualBuff)
		out := rcvr.GetSelector()
		in := NewJsonTransmitter(inIo).GetSelector()
		v := browse.NewVisitor(in)
		v.Out = browse.NewVisitor(out)
		err = in(browse.READ_VALUE, module, v)
		if err != nil {
			t.Error("failed to transmit json", err)
		} else {
			rcvr.Flush();
			actual := string(actualBuff.Bytes())
			t.Log("Round Trip:", actual)
			expected := "{\"hobbies\":[{\"birding\":{\"favorite-species\":\"towhee\"}}]}"
			if actual != expected {
				t.Error(actual, "!=", expected)
			}
		}
	}
}