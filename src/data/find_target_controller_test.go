package data

import (
	"bytes"
	"fmt"
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
	node.OnNext = func(*Selection, ListRequest) (Node, []*Value, error) {
		return node, nil, nil
	}
	node.OnSelect = func(*Selection, ContainerRequest) (Node, error) {
		return node, nil
	}
	tests := [][]string {
		{"", "m"},
		{"a","m/a"},
		{"b","m/b"},
		{"b=x","m/b=x"},
		{"a/aa=key/aab","m/a/aa=key/aab"},
	}
	for _, test := range tests {
		t.Log(test[0])
		found := Select(module, node).Find(test[0])
		if found.LastErr != nil {
			t.Error(found.LastErr)
		} else if found.Selection == nil {
			t.Errorf("Target for %s not found", test[0])
		} else {
			actual := found.Selection.path.String()
			if test[1] != actual {
				t.Errorf("Wrong state path\nExpected:%s\n  Actual:%s", test[1], actual)
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

		var rdr Node
		for i, test := range tests {
			rdr = NewJsonReader(strings.NewReader(json)).Node()
			found := Select(module, rdr).Find(test.path)
			if found.LastErr != nil {
				t.Error(err)
			} else if found.Selection == nil {
				t.Error("path not found " + test.path)
			}
			var actualBuff bytes.Buffer
			out := NewJsonWriter(&actualBuff).Node()
			err = found.Push(out).Upsert().LastErr
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
