package process
import (
	"testing"
	"schema/yang"
	"data"
	"schema"
	"bytes"
)

type ScriptTestSetup struct {
	Module *schema.Module
	Store *data.BufferStore
	Data *data.StoreData
}

func ScriptSetup(mstr string, t *testing.T) (setup *ScriptTestSetup) {
	setup = &ScriptTestSetup{}
	var err error
	setup.Module, err = yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	setup.Store = data.NewBufferStore()
	setup.Data = data.NewStoreData(setup.Module, setup.Store)
	return
}

func (setup *ScriptTestSetup) ToString(t *testing.T) string {
	var actualBuff bytes.Buffer
	out := data.NewJsonWriter(&actualBuff).Node()
	err := data.NodeToNode(setup.Data.Node(), out, setup.Data.Schema()).Insert()
	if err != nil {
		t.Error(err)
	}
	return actualBuff.String()
}

func TestScriptExec(t *testing.T) {
	a := ScriptSetup(moduleA, t)
	z := ScriptSetup(moduleZ, t)
	main := &Script{Name:"main"}
	txy := &Select{
		On : "b",
		Into : "y",
	}
	txy.AddOperation(&Set{
		Name : "x",
		Expression : &Primative{Var:"c"},
	})
	main.AddOperation(&Set{
		Name : "u",
		Expression : &Primative{Var:"f"},
	})
	main.AddOperation(txy)
	tests := []struct {
		aPath string
		aValue *data.Value
		expected string
	} {
		{
			"f",
			&data.Value{Str:"Eff"},
			`{"u":"Eff"}`,
		},
		{
			"b=Cee1/c",
			&data.Value{Str:"Cee1"},
			`{"y":[{"x":"Cee1"}]}`,
		},
	}
	for i, test := range tests {
		p := NewProcess(a.Data.Node(), a.Module).Into(z.Data.Node(), z.Module)
		a.Store.Clear()
		z.Store.Clear()
		a.Store.Values[test.aPath] = test.aValue
		err := p.RunScript(main)
		if err != nil {
			t.Error(err)
		} else {
			actual := z.ToString(t)
			expected := test.expected
			if actual != test.expected {
				t.Errorf("test #%d\nExpected:%s\n  Actual:%s", i + 1, expected, actual)
			}
		}
	}
}

var moduleA = `
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

var moduleZ = `
module z {
	prefix "";
	namespace "";
	revision 0;
	leaf u {
		type string;
	}
	list y {
		key "x";
		leaf x {
			type string;
		}
		container w {
			leaf v {
				type string;
			}
		}
	}
}
`