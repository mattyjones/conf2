package data

import (
	"bytes"
	"fmt"
	"schema"
	"schema/yang"
	"strings"
	"testing"
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
	module, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	node := &MyNode{}
	node.OnNext = func(*Selection, *schema.List, bool, []*Value, bool) (Node, error) {
		return node, nil
	}
	node.OnSelect = func(*Selection, schema.MetaList, bool) (Node, error) {
		return node, nil
	}
	root := NewSelection(node, module)
	var selection *Selection
	tests := []struct {
		path     string
		expected string
	}{
		{"", "m"},
		{"a", "m/a"},
		{"b", "m/b"},
		{"b=x", "m/b=x"},
		{"a/aa=key/aab", "m/a/aa=key/aab"},
	}
	for _, test := range tests {
		selection, err = WalkPath(root, NewPath(test.path))
		if selection == nil {
			t.Errorf("Target for %s not found", test.path)
		} else {
			actual := selection.String()
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

		tests := []struct {
			path     string
			expected string
		}{
			{"", strings.Replace(json, "\n", "", -1)},
			{"birding", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}],"reference":{"name":"Peterson's Guide"}}`},
			{"birding/lifer=towhee", `{"species":"towhee","location":"Hammonasset, CT"}`},
			{"birding?depth=1", `{"lifer":[],"reference":{}}`},
			{"birding/lifer", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}]}`},
			{"birding/lifer?depth=1", `{"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}]}`},
			{"birding/reference", `{"name":"Peterson's Guide"}`},
		}

		var in *Selection
		var rdr Node
		for i, test := range tests {
			inIo := strings.NewReader(json)
			if rdr, err = NewJsonReader(inIo).Node(); err != nil {
				t.Fatal(err)
			}
			in = NewSelection(rdr, module)
			if err != nil {
				t.Error(err)
			}
			p := NewPath(test.path)
			if in, err = WalkPath(in, p); err != nil {
				t.Error(err)
			}
			var actualBuff bytes.Buffer
			out := NewJsonWriter(&actualBuff).Selector(in.State)
			err = ControlledUpsert(in, out, LimitedWalk(p.Query))
			if err != nil {
				t.Error(err)
			} else {
				actual := string(actualBuff.Bytes())
				if actual != test.expected {
					msg := fmt.Sprintf("Failed subtest #%d - '%s'\nExpected:'%s'\n  Actual:'%s'",
						i+1, test.path, test.expected, actual)
					t.Error(msg)
				}
			}
		}
	}
}