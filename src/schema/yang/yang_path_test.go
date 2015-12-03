package yang

import (
	"fmt"
	"testing"
	"schema"
)

// This file should live in schema package, but I really needed to use
// yang calls to setup tests and therefore avoiding circular deps.


func TestPathEmpty(t *testing.T) {
	p, _ := schema.ParsePath("", &schema.Container{})
	if p.Len() != 1 {
		t.Error("expected one segments")
	}
}

func TestPathStringAndEqual(t *testing.T) {
	m, err := LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := [][]string {
		{ "", "m"},
		{"a/b", "m/a/b"},
		{"a/b=x", "m/a/b=x"},
		{"a/b=y/e", "m/a/b=y/e"},
		{"x=9", "m/x=9"},
	}
	for _, test := range tests {
		p, e := schema.ParsePath(test[0], m)
		if e != nil {
			t.Error(e)
		}
		actual := p.String()
		if test[1] != actual {
			t.Errorf("\nExpected: '%s'\n  Actual:'%s'", test[1], actual)
		}

		// Test equals
		p2, _ := schema.ParsePath(test[0], m)
		if ! p.Equal(p2) {
			t.Errorf("%s does not equal itself", test)
		}
	}
}

var mstr = `
module m {
	prefix "";
	namespace "";
	revision 0;
	container a {
		list b {
			key "d";
			leaf d {
				type string;
			}
			container e {
				leaf g {
					type string;
				}
			}
			leaf f {
				type string;
			}
		}
	}
	list x {
		key "y";
		leaf y {
			type int32;
		}
	}
}
`
func TestPathSegment(t *testing.T) {
	m, err := LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		in       string
		expected []string
	}{
		{"a/b", []string{"m", "a", "b"}},
		{"a/b=x", []string{"m", "a", "b"}},
		{"a/b=y/e", []string{"m", "a", "b", "e"}},
		{"a/b?foo=1", []string{"m", "a", "b"}},
	}
	for _, test := range tests {
		p, e := schema.ParsePath(test.in, m)
		if e != nil {
			t.Errorf("Error parsing %s : %s", test.in, e)
		}
		if len(test.expected) != p.Len() {
			t.Errorf("Expected %d segments for %s but got %d", len(test.expected), test.in, p.Len())
		}
		segments := p.Segments()
		for i, e := range test.expected {
			if e != segments[i].Meta().GetIdent() {
				msg := "expected to find \"%s\" as segment number %d in \"%s\" but got \"%s\""
				t.Error(fmt.Sprintf(msg, e, i, test.in, segments[i].Meta().GetIdent()))
			}
		}
	}
}

func TestPathSegmentKeys(t *testing.T) {
	m, err := LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	none := []interface{}{}
	tests := []struct {
		in       string
		expected [][]interface{}
	}{
		{"a/b", [][]interface{}{none, none, none}},
		{"a/b=c/e", [][]interface{}{none, none, []interface{}{"c"}, none}},
		{"x=9", [][]interface{}{none, []interface{}{9}}},
	}
	for _, test := range tests {
		p, e := schema.ParsePath(test.in, m)
		if e != nil {
			t.Errorf("Error parsing %s : %s", test.in, e)
		}
		if len(test.expected) != p.Len() {
			t.Error("wrong number of expected segments for", test.in)
		}
		segments := p.Segments()
		for i, expected := range test.expected {
			for j, key := range expected {
				if segments[i].Key()[j].Value() != key {
					desc := fmt.Sprintf("\"%s\" segment \"%s\" - expected \"%s\" - got \"%s\"",
						test.in, segments[i].Meta().GetIdent(), key, segments[i].Key()[j])
					t.Error(desc)
				}
			}
		}
	}
}

