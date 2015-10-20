package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"fmt"
	"schema"
)

func TestFindTargetIterator(t *testing.T) {
	mstr := `
module m {
    prefix "";
    namespace "";
	revision 0;
	container a {
	  list aa {
	    key "aaa";
	  	leaf aaa {
	  		type string;
	  	}
	  	container aab {
	  	  leaf aaba {
	  	    type string;
	  	  }
	  	}
	  }
	}
	list b {
	    key "ba";
		leaf ba {
			type string;
		}
		list bb {
			key "bba";
			leaf bba {
		    	type string;
			}
		}
	}
}
`
	module, err := yang.LoadModuleFromByteArray([]byte(mstr), nil);
	if err != nil {
		t.Fatal(err)
	}
	selection := &MySelection{}
	selection.OnNext = func(*WalkState, *schema.List, []*Value, bool) (Selection, error) {
		return selection, nil
	}
	selection.OnSelect = func(*WalkState, schema.MetaList) (Selection, error) {
		return selection, nil
	}
	rootState := NewWalkState(module)
	var s Selection
	var state *WalkState
	tests := []struct {
		path string
		expected string
	}{
		{"", 					"m/<nil>"},
		{"a", 					"m/a/<nil>"},
		{"b", 					"m/b/<nil>"},
		{"b=x",					"m/b=x/<nil>"},
		{"a/aa=key/aab",     	"m/a/aa=key/aab/<nil>"},
	}
	for _, test := range tests {
		s, state, err = WalkPath(rootState, selection, NewPath(test.path))
		if s == nil {
			t.Errorf("Target for %s not found", test.path)
		} else {
			actual := state.String()
			if test.expected != actual {
				t.Errorf("Wrong state path for %s\nExpected:%s\n  Actual:%s", test.path, test.expected, actual)
			}
		}
	}
}


func TestFindTargetAndInsert(t *testing.T) {
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
				t.Error(err)
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
