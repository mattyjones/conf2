package data

import (
	"schema/yang"
	"testing"
)

func TestBrowserPair(t *testing.T) {
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
			leaf aab {
				config "false";
				type string;
			}
			container aac {
				leaf aaca {
					type string;
				}
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
	{
		t.Log("Testing Init")
		oper := NewBufferStore()
		oper.Values["a/aa/aab"] = &Value{Str: "b"}
		operBrowser := NewStoreData(m, oper)

		config := NewBufferStore()
		configBrowser := NewStoreData(m, config)
		config.Values["a/aa/aaa"] = &Value{Str: "a"}
		var pair *DocumentPair
		if pair, err = NewDocumentPair(operBrowser, configBrowser); err != nil {
			t.Error(err)
		}
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}

		oper.Values["a/aa/aac/aaca"] = &Value{Str: "c"}
		var selection *Selection
		selection, err = pair.Selector(NewPath("a/aa/aac"))
		if selection == nil {
			t.Error("nil selection")
		}
	}
	{
		t.Log("Testing Edit")
		edit := NewBufferStore()
		edit.Values["a/ab"] = &Value{Str: "ab"}
		edit.Values["a/aa/aab"] = &Value{Str: "ab"}
		editBrowser := NewStoreData(m, edit)

		oper := NewBufferStore()
		operBrowser := NewStoreData(m, oper)
		config := NewBufferStore()
		configBrowser := NewStoreData(m, config)
		var pair *DocumentPair
		if pair, err = NewDocumentPair(operBrowser, configBrowser); err != nil {
			t.Error(err)
		}
		var in, out *Selection
		p := NewPath("")
		if in, err = editBrowser.Selector(p); err != nil {
			t.Fatal(err)
		}
		if out, err = pair.Selector(p); err != nil {
			t.Fatal(err)
		}
		if err = Upsert(in, out); err != nil {
			t.Fatal(err)
		}
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}
		if len(config.Values) != 1 {
			t.Error("Expected 1 items got ", len(config.Values))
		}
	}
}