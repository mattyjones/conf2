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
		out := NewJsonReceiver(&actualBuff)
		//dbg := &browse.DebuggingWriter{Delegate:out}
		if err != nil {
			t.Error(err)
		}
		in, err := NewJsonTransmitter(inIo).GetSelector(module)
		if err != nil {
			t.Error(err)
		}
		err = browse.Transfer(in, out)
		if err != nil {
			t.Error("failed to transmit json", err)
		} else {
			out.Flush();
			actual := string(actualBuff.Bytes())
			t.Log("Round Trip:", actual)
			expected := "{\"hobbies\":[{\"birding\":{\"favorite-species\":\"towhee\"}}]}"
			if actual != expected {
				t.Error(actual, "!=", expected)
			}
		}
	}
}