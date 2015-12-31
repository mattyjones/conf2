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
  tone = color
  select facets into attributes
    let label = "orangey-"
    attribute = facet
    if data.value
      additional.detail = concat(label, data.value)
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
			"color",
			&schema.Value{Str:"red"},
			`{"tone":"red"}`,
		},
		{
			"facets=seeds/facet",
			&schema.Value{Str:"seeds"},
			`{"attributes":[{"attribute":"seeds"}]}`,
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
module apple {
	prefix "";
	namespace "";
	revision 0;
	leaf color {
		type string;
	}
	list facets {
		key "facet";
		leaf facet {
			type string;
		}
		container data {
			leaf value {
				type string;
			}
		}
	}
}
`

var moduleZ = `
module orange {
	prefix "";
	namespace "";
	revision 0;
	leaf tone {
		type string;
	}
	list attributes {
		key "attribute";
		leaf attribute {
			type string;
		}
		container additional {
			leaf detail {
				type string;
			}
		}
	}
}
`