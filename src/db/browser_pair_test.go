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
	oper := BufferStore{}
	oper["/a/aa/aab"] = &browse.Value{Str:"b"}
	operBrowser := NewConfig(m, oper)
//	state, sel, err := operBrowser.Selector(browse.NewPath(""), browse.READ)
//	browse.Walk(sel, state, browse.WalkAll())
//
	config := BufferStore{}
	configBrowser := NewConfig(m, config)
	config["/a/aa/aaa"] = &browse.Value{Str:"a"}
	pair := NewBrowserPair(operBrowser, configBrowser)
	pair.Init()
	if len(oper) != 2 {
		t.Error("Expected 2 items got ", len(oper))
	}
	edit := BufferStore{}
	edit["/a/ab"] = &browse.Value{Str:"ab"}
	editBrowser := NewConfig(m, edit)
	browse.Upsert(browse.NewPath("a"), editBrowser, pair)
	if len(oper) != 3 {
		t.Error("Expected 3 items got ", len(oper))
	}
}