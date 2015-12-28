package process
import (
	"testing"
	"schema/yang"
)

func TestTableGet(t *testing.T)  {
	_, err := yang.LoadModuleFromByteArray([]byte(tableModuleStr), nil)
	if err != nil {
		t.Fatal("not implemented")
	}
}

var tableModuleStr = `
module a {
	prefix "";
	namespace "";
	revision 0;
	leaf f {
		type string;
	}
	list b {
		key "c";
		leaf c {
			type string;
		}
		container d {
			leaf e {
				type string;
			}
		}
	}
}
`

