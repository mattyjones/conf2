package data
import (
	"testing"
	"schema/yang"
	"strings"
	"bytes"
	"encoding/json"
)

var mstr = `
module m {
	namespace "";
	prefix "";
	revision 0;
	container a {
		container b {
			leaf x {
				type string;
			}
		}
	}
	list p {
		key "k";
		leaf k {
			type string;
		}
		container q {
			leaf s {
				type string;
			}
		}
		list r {
			leaf z {
				type int32;
			}
		}
	}
}
`

func TestBucketWrite(t *testing.T) {
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		data string
		path string
	} {
		{
			`{"a":{"b":{"x":"waldo"}}}`,
			"a.b.x",
		},
		{
			`{"p":[{"k":"walter"},{"k":"waldo"},{"k":"weirdo"}]}`,
			"p.1.k",
		},
	}
	for _, test := range tests {
		bd := &BucketData {Meta: m}
		var in, sel *Selection
		if sel, err = bd.Selector(NewPath("")); err != nil {
			t.Fatal(err)
		}
		in, err = NewJsonReader(strings.NewReader(test.data)).Selector(sel.State)
		if err = Insert(in, sel); err != nil {
			t.Error(err)
		}
		actual := MapValue(bd.Root, test.path)
		if actual != "waldo" {
			t.Error(actual)
		}
	}
}

func TestBucketRead(t *testing.T) {
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		data map[string]interface{}
		expected string
	} {
		{
			map[string]interface{}{
				"a" : map[string]interface{}{
					"b" : map[string]interface{}{
						"x" : "waldo",
					},
				},
			},
			`{"a":{"b":{"x":"waldo"}}}`,
		},
		{
			map[string]interface{}{
				"p" : []map[string]interface{}{
					map[string]interface{}{"k" :"walter"},
					map[string]interface{}{"k" :"waldo"},
					map[string]interface{}{"k" :"weirdo"},
				},
			},
			`{"p":[{"k":"walter"},{"k":"waldo"},{"k":"weirdo"}]}`,
		},
	}

	for _, test := range tests {
		bd := &BucketData {Meta: m, Root : test.data}
		var sel *Selection
		if sel, err = bd.Selector(NewPath("")); err != nil {
			t.Fatal(err)
		}
		var buff bytes.Buffer
		out := NewJsonWriter(&buff).Selector(sel.State)
		if err = Insert(sel, out); err != nil {
			t.Error(err)
		}
		actual := buff.String()
		if actual != test.expected {
			t.Errorf("\nExpected:%s\n  Actual:%s", test.expected, actual)
		}
	}
}

func TestBucketDecode(t *testing.T) {
	var err error
	dataJson := `{"a":{"b":{"x":"waldo"}},"p":[{"k":"walter"},{"k":"waldo"},{"k":"weirdo"}]}`
	var data map[string]interface{}
	if err = json.Unmarshal([]byte(dataJson), &data); err != nil {
		t.Error(err)
	}
	if MapValue(data, "a.b.x") != "waldo" {
		t.Error("can't find waldo")
	}
	if MapValue(data, "p.1.k") != "waldo" {
		t.Error("can't find waldo")
	}
}
