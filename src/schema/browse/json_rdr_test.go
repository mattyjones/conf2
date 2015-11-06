package browse

import (
	"testing"
	"schema/yang"
	"strings"
)

func TestJsonWalk(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0;
	list hobbies {
		key "name";
	    leaf name {
	      type string;
	    }
		container favorite {
		    leaf common-name {
		      type string;
		    }
		    leaf location {
				type string;
		    }
		}
	}
}
	`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Error("bad module", err)
	} else {
		json := `{"hobbies":[
{"name":"birding", "favorite": {"common-name" : "towhee", "extra":"double-mint", "location":"out back"}},
{"name":"hockey", "favorite": {"common-name" : "bruins", "location" : "Boston"}}
]}`

		tests := [] struct {
			path string
			expectedMeta	string
		} {
			{ "hobbies", 			"json-test/hobbies/<nil>" },
			{ "hobbies=birding", 	"json-test/hobbies=birding/<nil>" },
			{ "hobbies=birding/favorite", "json-test/hobbies=birding/favorite/<nil>" },
		}
		var in, selection *Selection
		for _, test := range tests {
			in, err = NewJsonReader(strings.NewReader(json)).Selector(module)
			selection, err = WalkPath(in, NewPath(test.path))
			if err != nil {
				t.Error("failed to transmit json", err)
			} else if selection == nil {
				t.Error(test.path, "- Target not found, state nil")
			} else if (selection.Path().String() != test.expectedMeta) {
				t.Error(test.path, "-", test.expectedMeta, "!=", selection.String())
			}
		}
	}
}


