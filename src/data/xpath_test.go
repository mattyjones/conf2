package data
import (
	"testing"
	"schema/yang"
	"strings"
)


func TestXPath(t *testing.T) {
	moduleStr := `
module m {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	list l1 {
		key "a";
		leaf a {
			type string;
		}
		leaf b {
			type string;
		}
	}
}
	`
	dataStr := `
{
	"l1" : [{
	  "a" : "a-one",
	  "b" : "b-one"
	},{
	  "a" : "a-two",
	  "b" : "b-two"
	}]
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil)
	if err != nil {
		t.Fatal(err)
	}
	in := NewJsonReader(strings.NewReader(dataStr)).Node()
	store := NewBufferStore()
	data := NewStoreData(m, store)
	if err = NodeToNode(in, data.Node(), m).Insert(); err != nil {
		t.Fatal(err)
	}
	sel := NewSelection(data.Node(), data.Schema())
	actual, xpathErr := XPath{}.Get(sel, "l1=a-two/b")
	if xpathErr != nil {
		t.Error(xpathErr)
	}
	if actual != "b-two" {
		t.Error(actual)
	}
}
