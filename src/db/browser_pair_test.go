package db
import (
	"testing"
	"schema/yang"
	"schema/browse"
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
		oper.Values["a/aa/aab"] = &browse.Value{Str:"b"}
		operBrowser := browse.NewStoreBrowser(m, oper)

		config := browse.NewBufferStore()
		configBrowser := browse.NewStoreBrowser(m, config)
		config.Values["a/aa/aaa"] = &browse.Value{Str:"a"}
		pair := NewBrowserPair(operBrowser, configBrowser)
		pair.Init()
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}
	}
	{
		t.Log("Testing Edit")
		edit := browse.NewBufferStore()
		edit.Values["a/ab"] = &browse.Value{Str:"ab"}
		edit.Values["a/aa/aab"] = &browse.Value{Str:"ab"}
		editBrowser := browse.NewStoreBrowser(m, edit)

		oper := browse.NewBufferStore()
		operBrowser := browse.NewStoreBrowser(m, oper)
		config := browse.NewBufferStore()
		configBrowser := browse.NewStoreBrowser(m, config)
		pair := NewBrowserPair(operBrowser, configBrowser)

		browse.Upsert(browse.NewPath("a"), editBrowser, pair)
		if len(oper.Values) != 2 {
			t.Error("Expected 2 items got ", len(oper.Values))
		}
		if len(config.Values) != 1 {
			t.Error("Expected 1 items got ", len(config.Values))
		}
	}
}