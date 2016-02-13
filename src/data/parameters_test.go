package data

import (
	"schema/yang"
	"testing"
	"schema"
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
		default "y";
	}
	leaf c {
		type string;
	}
}
`
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	p := &Parameters{}
	obj := struct {
		A string
		B string
		C string
	} {}
	p.Collect("c", &Value{Type: schema.NewDataType(nil, "string"), Str: "z"})
	p.Record("b")
	n := MarshalContainer(&obj)
	sel := Select(m, n)
	err = p.Finish(sel, n)
	if err != nil {
		t.Error(err)
	}
	if obj.A != "x" {
		t.Error("wrong parameter default value ", obj.A)
	}
	if obj.B != "" {
		t.Error("wrong parameter default value ", obj.B)
	}
	if obj.C != "z" {
		t.Error("wrong parameter default value ", obj.C)
	}
}
