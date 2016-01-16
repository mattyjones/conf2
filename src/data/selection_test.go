package data

import (
	"testing"
	"strings"
	"schema/yang"
	"regexp"
)

var selectionTestModule = `
module m {
	prefix "";
	namespace "";
	revision 0;
	container message {
		leaf hello {
			type string;
		}
		container deep {
			leaf goodbye {
				type string;
			}
		}
	}
}
`

func TestSelectionEvents(t *testing.T) {
	m, err := yang.LoadModuleFromByteArray([]byte(selectionTestModule), nil)
	if err != nil {
		t.Fatal(err)
	}
	store := NewBufferStore()
	storeData := NewStoreData(m, store)
	sel := NewSelection(storeData.Node(), storeData.Schema())
	var relPathFired bool
	sel.OnPath(NEW, "m/message", func() error {
		relPathFired = true
		return nil
	})
	var regexFired bool
	sel.OnRegex(END_EDIT, regexp.MustCompile(".*"), func() error {
		regexFired = true
		return nil
	})
	json := NewJsonReader(strings.NewReader(`{"message":{"hello":"bob"}}`)).Node()
	err = NodeToSelection(json, sel).Upsert()
	if err != nil {
		t.Fatal(err)
	}
	if !relPathFired {
		t.Fatal("Event not fired")
	}
	if !regexFired {
		t.Fatal("regex not fired")
	}
}

func TestSelectionPeek(t *testing.T) {
	m, err := yang.LoadModuleFromByteArray([]byte(selectionTestModule), nil)
	if err != nil {
		t.Fatal(err)
	}
	var expected = "Justin Bieber Fan Club Member"
	n := &MyNode{Internal:expected}
	sel := NewSelection(n, m)
	actual :=  sel.Peek()
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}
