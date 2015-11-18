package db

import (
	"schema/browse"
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
		oper := browse.NewBufferStore()
		oper.Values["a/aa/aab"] = &browse.Value{Str: "b"}
		operBrowser := browse.NewStoreData(m, oper)

		config := browse.NewBufferStore()
		configBrowser := browse.NewStoreData(m, config)
		config.Values["a/aa/aaa"] = &browse.Value{Str: "a"}
		var pair *DocumentPair
		rootPath := browse.NewPath("")
		operSel, _ := operBrowser.Selector(rootPath)
		configSel, _ := configBrowser.Selector(rootPath)
		if pair, err = NewDocumentPair(operSel, configSel); err != nil {
			t.Error(err)
		}
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}

		oper.Values["a/aa/aac/aaca"] = &browse.Value{Str: "c"}
		var selection *browse.Selection
		selection, err = pair.Selector(browse.NewPath("a/aa/aac"))
		if selection == nil {
			t.Error("nil selection")
		}
	}
	{
		t.Log("Testing Edit")
		edit := browse.NewBufferStore()
		edit.Values["a/ab"] = &browse.Value{Str: "ab"}
		edit.Values["a/aa/aab"] = &browse.Value{Str: "ab"}
		editBrowser := browse.NewStoreData(m, edit)

		oper := browse.NewBufferStore()
		operBrowser := browse.NewStoreData(m, oper)
		config := browse.NewBufferStore()
		configBrowser := browse.NewStoreData(m, config)
		var pair *DocumentPair
		rootPath := browse.NewPath("")
		operSel, _ := operBrowser.Selector(rootPath)
		configSel, _ := configBrowser.Selector(rootPath)
		if pair, err = NewDocumentPair(operSel, configSel); err != nil {
			t.Error(err)
		}
		var in, out *browse.Selection
		p := browse.NewPath("")
		if in, err = editBrowser.Selector(p); err != nil {
			t.Fatal(err)
		}
		if out, err = pair.Selector(p); err != nil {
			t.Fatal(err)
		}
		if err = browse.Upsert(in, out); err != nil {
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
