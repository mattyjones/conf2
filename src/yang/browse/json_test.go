package browse

import (
	"testing"
	"yang"
	"bytes"
	"strings"
	"fmt"
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
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr)); err != nil {
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
			} else if (s.Meta.GetIdent() != test.expectedMeta) {
				t.Error(test.path, "-", test.expectedMeta, "!=", s.Meta.GetIdent())
			}
		}
	}
}


func TestJsonPath(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
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
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr)); err != nil {
		t.Error("bad module", err)
	} else {
		json := `{"birding":{
"lifer":[{"species":"towhee","location":"Hammonasset, CT"},{"species":"robin","location":"East Rock, CT"}],
"reference":{"name":"Peterson's Guide"}
}}`

		tests := [] struct {
			path string
		} {
			//{ "birding" },
			//{ "birding/lifer" },
			{ "birding/lifer=towhee" },
			//{ "birding/reference" },
		}

		for _, test := range tests {
			p, _ := NewPath(test.path)
			inIo := strings.NewReader(json)
			var actualBuff bytes.Buffer
			in, err := NewJsonReader(inIo).GetSelector(module)
			if err != nil {
				t.Error(err)
			}
			if ref, err := WalkPath(in, p); err != nil {
				t.Error(err)
			} else {
				out := NewJsonWriter(&actualBuff)
				to, _ := out.GetSelector()
fmt.Println("json_test:=====================")
				err = Insert(ref, to)
				if err != nil {
					t.Error("failed to transmit json", err)
				} else {
					actual := string(actualBuff.Bytes())
					t.Log("Round Trip:", actual)
				}
			}
		}
	}
}