package browse

import (
	"testing"
	"schema/yang"
	"bytes"
	"schema"
	"strings"
)
func TestKeyListBuilderInBufferStore(t *testing.T) {
	tests := []struct {
		path string
		expected string
	} {
		{ "a/a", "" },
		{ "a/b", "x" },
		{ "a/c", "y|z" },
	}
	s := &storeSelector{}
	store := NewBufferStore()
	s.store = store
	v := &Value{}
	store.Values["a/a/c"] = v
	store.Values["a/b=x/c"] = v
	store.Values["a/c=y/c"] = v
	store.Values["a/c=y/c"] = v
	store.Values["a/c=z/q/f=yy/fg=gf/gf"] = v
	meta := &schema.List{Ident:"c", Keys:[]string{"k"}}
	meta.AddMeta(&schema.Leaf{Ident:"k", DataType:&schema.DataType{Format:schema.FMT_STRING}})
	for _, test := range tests {
		keys, err := s.store.KeyList(test.path, meta)
		if err != nil {
			t.Error(err)
		}
		actual := strings.Join(keys, "|")
		if actual != test.expected {
			t.Errorf("Test failed for path %s\nExpected:%s\n  Actual:%s", test.path, test.expected, actual)
		}
	}
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

func TestStoreBrowserKeyValueRead(t *testing.T) {
	store := NewBufferStore()
	m := keyValuesTestModule()
	kv := NewStoreBrowser(m, store)
	store.Values["a/aa/aaa"] = &Value{Str:"hi"}
	store.Values["b=x/ba"] = &Value{Str:"x"}
	var actualBytes bytes.Buffer
	json := NewJsonWriter(&actualBytes, m)
	err := Insert(NewPath(""), kv, json)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes.Bytes())
	expected := `{"a":{"aa":{"aaa":"hi"}},"b":[{"ba":"x"}]}`
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}

func TestStoreBrowserValueEdit(t *testing.T) {
	store := NewBufferStore()
	m := keyValuesTestModule()
	kv := NewStoreBrowser(m, store)
	inputJson := `{"a":{"aa":{"aaa":"hi"}},"b":[{"ba":"x"}]}`
	json := NewJsonReader(strings.NewReader(inputJson), m)
	err := Insert(NewPath(""), json, kv)
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Values) != 2 {
		t.Error("Expected 2 items")
	}
	expectations := []struct {
		path string
		value string
	}{
		{"a/aa/aaa", "hi"},
		{"b=x/ba", "x"},
	}
	if len(expectations) != len(store.Values) {
		t.Errorf("Expected %d items but got", len(expectations), len(store.Values))
	}
	for _, expected := range expectations {
		v, found := store.Values[expected.path]
		if !found {
			t.Error("Could not find item", expected.path, "\nItems: ", store)
		} else if v.Str != expected.value {
			t.Error("Expected value to be %s but was %s", expected.value, v.Str)
		}
	}
}

func TestStoreBrowserKeyValueEdit(t *testing.T) {
	store := NewBufferStore()
	m := keyValuesTestModule()
	kv := NewStoreBrowser(m, store)
	store.Values["b=x/ba"] = &Value{Str:"z"}

	// change key
	inputJson2 := `{"ba":"y"}`
	json2 := NewJsonFragmentReader(strings.NewReader(inputJson2))
	err := Update(NewPath("b=x"), json2, kv)
	if err != nil {
		t.Fatal(err)
	}
	if v, newKeyExists := store.Values["b=y/ba"]; !newKeyExists {
		t.Error("Edit key value not made")
	} else if v.Str != "y" {
		t.Error("Wrong key value")
	}
	if _, oldKeyExists := store.Values["/b=x/ba"]; oldKeyExists {
		t.Error("Old key was not removed")
	}
}