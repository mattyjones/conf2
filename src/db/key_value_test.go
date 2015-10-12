package db

import (
	"testing"
	"schema/yang"
	"schema/browse"
	"bytes"
	"schema"
	"strings"
	"sort"
)

func TestKeyValueBuildList(t *testing.T) {
	tests := []struct {
		path string
		expected string
	} {
		{ "/a/a", "" },
		{ "/a/b", "x" },
		{ "/a/c", "y|z" },
	}
	s := &KeyValuesSelector{}
	store := make(KeyValueStore, 10)
	s.store = store
	v := &browse.Value{}
	store["/a/a/c"] = v
	store["/a/b=x/c"] = v
	store["/a/c=y/c"] = v
	store["/a/c=y/c"] = v
	store["/a/c=z/q/f=yy/fg=gf/gf"] = v
	meta := &schema.List{Ident:"c", Keys:[]string{"k"}}
	meta.AddMeta(&schema.Leaf{Ident:"k", DataType:&schema.DataType{Format:schema.FMT_STRING}})
	for _, test := range tests {
		keys, err := s.store.KeyList(test.path, meta)
		if err != nil {
			t.Error(err)
		}
		actual := joinTestValues(keys)
		if actual != test.expected {
			t.Errorf("Test failed for path %s\nExpected:%s\n  Actual:%s", test.path, test.expected, actual)
		}
	}
}

func joinTestValues(keys [][]*browse.Value) string {
	strs := make([]string, 0, len(keys))
	for _, key := range keys {
		for _, v := range key {
			strs = append(strs, v.String())
		}
	}
	sort.Strings(strs)
	return strings.Join(strs, "|")
}

func keyValuesTestModule() *schema.Module {
	mstr := `
module m {
	prefix "";
	namespace "";
	revision 0;
	container a {
		container aa {
			leaf aaa {
				type string;
			}
		}
		leaf ab {
			type string;
		}
	}
	list b {
		key "ba";
		leaf ba {
			type string;
		}
		container bb {
			leaf bba {
				type string;
			}
		}
		list bc {
			key "bca";
			leaf bca {
				type string;
			}
		}
	}
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		panic(err)
	}
	return m
}

func TestKeyValueRead(t *testing.T) {
	store := make(KeyValueStore, 100)
	m := keyValuesTestModule()
	kv := NewKeyValues(m, store)
	store["/a/aa/aaa"] = &browse.Value{Str:"hi"}
	store["/b=x/ba"] = &browse.Value{Str:"x"}
	var actualBytes bytes.Buffer
	json := browse.NewJsonWriter(&actualBytes, m)
	browse.Insert(browse.NewPath(""), kv, json)
	actual := string(actualBytes.Bytes())
	expected := `{"a":{"aa":{"aaa":"hi"}},"b":[{"ba":"x"}]}`
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}

func TestKeyValueEdit(t *testing.T) {
	store := make(KeyValueStore, 100)
	m := keyValuesTestModule()
	kv := NewKeyValues(m, store)
	inputJson := `{"a":{"aa":{"aaa":"hi"}},"b":[{"ba":"x"}]}`
	json := browse.NewJsonReader(strings.NewReader(inputJson), m)
	err := browse.Insert(browse.NewPath(""), json, kv)
	if err != nil {
		t.Fatal(err)
	}
	if len(store) != 2 {
		t.Error("Expected 2 items")
	}
	expectations := []struct {
		path string
		value string
	} {
		{"/a/aa/aaa", "hi"},
		{"/b=x/ba", "x"},
	}
	if len(expectations) != len(store) {
		t.Errorf("Expected %d items but got", len(expectations), len(store))
	}
	for _, expected := range expectations {
		v, found := store[expected.path]
		if !found {
			t.Error("Could not find item", expected.path)
		} else if v.Str != expected.value {
			t.Error("Expected value to be %s but was %s", expected.value, v.Str)
		}
	}
}