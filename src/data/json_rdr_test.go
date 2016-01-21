package data

import (
	"schema/yang"
	"strings"
	"testing"
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
	module, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil)
	if err != nil {
		t.Fatal(err)
	}
	json := `{"hobbies":[
{"name":"birding", "favorite": {"common-name" : "towhee", "extra":"double-mint", "location":"out back"}},
{"name":"hockey", "favorite": {"common-name" : "bruins", "location" : "Boston"}}
]}`
	tests := []string {
		"hobbies",
		"hobbies=birding",
		"hobbies=birding/favorite",
	}
	for _, test := range tests {
		rdr := NewJsonReader(strings.NewReader(json)).Node()
		found, selErr := NewSelection(module, rdr).Find(test)
		if selErr != nil {
			t.Error("failed to transmit json", err)
		} else if found == nil {
			t.Error(test, "- Target not found, state nil")
		} else {
			actual := found.Path().String()
			if actual != "json-test/" + test {
				t.Error("json-test/" + test, "!=", actual)
			}
		}
	}
}
