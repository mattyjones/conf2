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
	oper := NewBufferStore()
	operNode := NewStoreData(m, oper).Node()
	//oper.Values["a/aa/aaa"] = &schema.Value{Str:":hello"}
	config := NewBufferStore()
	configNode := NewStoreData(m, config).Node()
	//oper.Values["a/aa/aaa"] = &schema.Value{Str:":hello"}
	edit := `{"a":{"aa":{"aaa":"hello"}}}`
	in := NewJsonReader(strings.NewReader(edit)).Node()
	err = NodeToNode(in, Config(configNode, operNode), m).Insert()
	if err != nil {
		t.Error(err)
	}
	if len(oper.Values) != 1 {
		t.Error(len(oper.Values))
	}
	if len(config.Values) != 1 {
		t.Error(len(config.Values))
	}
}


