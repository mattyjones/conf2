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
		oper := BufferStore{}
		oper["a/aa/aab"] = &browse.Value{Str:"b"}
		operBrowser := NewStoreBrowser(m, oper)

		config := BufferStore{}
		configBrowser := NewStoreBrowser(m, config)
		config["a/aa/aaa"] = &browse.Value{Str:"a"}
		pair := NewBrowserPair(operBrowser, configBrowser)
		pair.Init()
		if len(oper) != 2 {
			t.Error("Expected 2 items got ", len(oper))
		}
	}
	{
		t.Log("Testing Edit")
		edit := BufferStore{}
		edit["a/ab"] = &browse.Value{Str:"ab"}
		edit["a/aa/aab"] = &browse.Value{Str:"ab"}
		editBrowser := NewStoreBrowser(m, edit)

		oper := BufferStore{}
		operBrowser := NewStoreBrowser(m, oper)
		config := BufferStore{}
		configBrowser := NewStoreBrowser(m, config)
		pair := NewBrowserPair(operBrowser, configBrowser)

		browse.Upsert(browse.NewPath("a"), editBrowser, pair)
		if len(oper) != 2 {
			t.Error("Expected 2 items got ", len(oper))
		}
		if len(config) != 1 {
			t.Error("Expected 1 items got ", len(config))
		}
	}
}