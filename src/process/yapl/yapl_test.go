package yapl

import (
	"testing"
	"process"
	"schema"
	"data"
)

func TestYaplExec(t *testing.T) {
	a := data.ModuleSetup(moduleA, t)
	z := data.ModuleSetup(moduleZ, t)
	scripts, err := Load(
`main
  u = f
  select b into y
    x = c
`)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		aPath string
		aValue *schema.Value
		expected string
	} {
		{
			"f",
			&schema.Value{Str:"Eff"},
			`{"u":"Eff"}`,
		},
		{
			"b=Cee1/c",
			&schema.Value{Str:"Cee1"},
			`{"y":[{"x":"Cee1"}]}`,
		},
	}
	for i, test := range tests {
		p := process.NewProcess(a.Data.Node(), a.Module).Into(z.Data.Node(), z.Module)
		a.Store.Clear()
		z.Store.Clear()
		a.Store.Values[test.aPath] = test.aValue
		err := p.Run(scripts, "main")
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