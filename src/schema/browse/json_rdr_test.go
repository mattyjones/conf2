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
	revision 0000-00-00 {
		description "x";
	}
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
			{ "hobbies", 			"hobbies" },
			{ "hobbies=birding", 	"hobbies" },
			{ "hobbies=birding/favorite", "favorite" },
		}
		for _, test := range tests {

			inIo := strings.NewReader(json)
			if err != nil {
				t.Error(err)
			}
			in, err := NewJsonReader(inIo).GetSelector(module)
			if err != nil {
				t.Error(err)
			}
			p, _ := NewPath(test.path)
			s, err := WalkPath(in, p)
			if err != nil {
				t.Error("failed to transmit json", err)
			} else if (s.WalkState().Meta.GetIdent() != test.expectedMeta) {
				t.Error(test.path, "-", test.expectedMeta, "!=", s.WalkState().Meta.GetIdent())
			}
		}
	}
}


