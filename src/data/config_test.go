package data
import (
	"testing"
	"schema/yang"
	"schema"
	"strings"
)

func TestConfig(t *testing.T) {
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
	var err error
	var m *schema.Module
	if m, err = yang.LoadModuleFromByteArray([]byte(mstr), nil); err != nil {
		t.Fatal(err)
	}
	operStore := NewBufferStore()
	oper := NewStoreData(m, operStore)
	//oper.Values["a/aa/aaa"] = &Value{Str:":hello"}
	persistStore := NewBufferStore()
	persist := NewStoreData(m, persistStore)
	//oper.Values["a/aa/aaa"] = &Value{Str:":hello"}
	edit := `{"a":{"aa":{"aaa":"hello"}}}`

	sel := oper.Select().Fork(Config(persist.Node(), oper.Node()))
	if err = sel.Pull(NewJsonReader(strings.NewReader(edit)).Node()).Insert(); err != nil {
		t.Error(err)
	}
	if len(operStore.Values) != 1 {
		t.Error(len(operStore.Values))
	}
	if len(persistStore.Values) != 1 {
		t.Error(len(persistStore.Values))
	}
}


