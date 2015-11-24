package browse

import (
	"testing"
	"strings"
	"schema/yang"
	"regexp"
)

func TestSelectionEvents(t *testing.T) {
	mstr := `
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
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	store := NewBufferStore()
	storeData := NewStoreData(m, store)
	sel, _ := storeData.Selector(NewPath(""))
	//out := NewJsonWriter(&null).Container()
	var r Node
	r, err = NewJsonReader(strings.NewReader(`{"message":{"hello":"bob"}}`)).Node()
	if err != nil {
		t.Fatal(err)
	}
	//sel := NewSelection(out, m)
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
	in := NewSelectionFromState(r, sel.State)
	err = Upsert(in, sel)
	if err != nil {
		t.Fatal(err)
	}
	if !relPathFired {
		t.Fatal("Event not fired")
	}
	if !regexFired {
		t.Fatal("regex not fired")
	}
//
//	sel, err = storeData.Selector(NewPath("message"))
//	if err != nil {
//		t.Fatal(err)
//	}
}