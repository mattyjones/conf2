package browse

import (
	"schema/yang"
	"testing"
)

func TestParameters(t *testing.T) {
	mstr := `
module m {
	prefix "";
	namespace "";
	revision 0;
	leaf a {
		type string;
		default "x";
	}
	leaf b {
		type string;
	}
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	p := NewParameters(m)
	v := p.Value("a")
	if v.Str != "x" {
		t.Error("wrong parameter default value ", v.Str)
	}
}
