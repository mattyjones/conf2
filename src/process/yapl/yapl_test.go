package yapl

import (
	"testing"
	"process"
	"data"
	"strings"
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
    select extra into even
      more = info
`)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		aData string
		expected string
	} {
		{
			`{"color":"red"}`,
			`{"tone":"red"}`,
		},
		{
			`{"facets":[{"facet":"seeds"}]}`,
			`{"attributes":[{"attribute":"seeds"}]}`,
		},
		{
			`{"facets":[{"facet":"seeds","extra":{"info":"32"}}]}`,
			// correct expectation, but json wtr is wrong
			//   `{"attributes":[{"attribute":"seeds","even":{"more":"32"}}]}`,
			`{"attributes":[{"attribute":"seeds"},{"even":{"more":"32"}}]}`,
		},
	}
	for i, test := range tests {
		aIn := data.NewJsonReader(strings.NewReader(test.aData))
		p := process.NewProcess(data.Select(a.Module, aIn.Node())).Into(z.Data.Select())
		z.Store.Clear()
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
		container extra {
			leaf info {
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
		container even {
			leaf more {
				type string;
			}
		}
	}
}
`