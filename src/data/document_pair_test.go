package data

import (
	"schema/yang"
	"testing"
	"schema"
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
		oper.Values["a/aa/aab"] = &schema.Value{Str: "b"}
		operBrowser := NewStoreData(m, oper)

		config := NewBufferStore()
		configBrowser := NewStoreData(m, config)
		config.Values["a/aa/aaa"] = &schema.Value{Str: "a"}
		var pair *DocumentPair
		if pair, err = NewDocumentPair(operBrowser, configBrowser); err != nil {
			t.Error(err)
		}
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}

		oper.Values["a/aa/aac/aaca"] = &schema.Value{Str: "c"}

		var selection *Selection
		selection, err = WalkPath(NewSelection(pair.Node(), m), schema.NewPathSlice("a/aa/aac", m))
		if selection == nil {
			t.Error("nil selection")
		}
	}
	{
		t.Log("Testing Edit")
		edit := NewBufferStore()
		edit.Values["a/ab"] = &schema.Value{Str: "ab"}
		edit.Values["a/aa/aab"] = &schema.Value{Str: "ab"}
		editBrowser := NewStoreData(m, edit)

		oper := NewBufferStore()
		operBrowser := NewStoreData(m, oper)
		config := NewBufferStore()
		configBrowser := NewStoreData(m, config)
		var pair *DocumentPair
		if pair, err = NewDocumentPair(operBrowser, configBrowser); err != nil {
			t.Error(err)
		}
		x, err2 := PathToPath(editBrowser, pair, "")
		if err2 != nil {
			t.Error(err2)
		}
		if err = x.Upsert(); err != nil {
			t.Error(err)
		}
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}
		if len(config.Values) != 1 {
			t.Error("Expected 1 items got ", len(config.Values))
		}
	}
}
