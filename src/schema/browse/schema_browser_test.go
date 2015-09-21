package browse
import (
	"testing"
	"schema"
	"schema/yang"
	"bytes"
	"fmt"
)

func printMeta(m schema.Meta, level string) {
	fmt.Printf("%s%s\n", level, m.GetIdent())
	if nest, isNest := m.(schema.MetaList); isNest {
		if len(level) >= 16 {
			panic("Max level reached")
		}
		i2 := schema.NewMetaListIterator(nest, false)
		for i2.HasNextMeta() {
			printMeta(i2.NextMeta(), level + "  ")
		}
	}
}

func TestYangBrowserRead(t *testing.T) {
	tests := []struct {
		yang string
		expected string
		read ReadFunc
	} {
		{
			`leaf c { type enumeration { enum a; enum b; } }`,
			`{"c":"a"}`,
			func(v *Value) error {
				v.Int = 0
				v.Str = "a"
				return nil
			},
    	},
		{
			`leaf-list c { type enumeration { enum a; enum b; } }`,
			`{"c":["a","b"]}`,
			func(v *Value) error {
				v.Intlist = []int {0, 1}
				v.Strlist = []string {"a", "b"}
				return nil
			},
		},
	}
	for _, test := range tests {
		s := makeMySelection(t, test.yang)
		s.OnRead = test.read
		s.State.Found = true
		actual := tojson(t, s)
		if actual != test.expected {
			msg := fmt.Sprintf("Expected:\"%s\" Actual:\"%s\"", test.expected, actual)
			t.Error(msg)
		}
	}
}

func makeMySelection(t *testing.T, yangFragment string) *MySelection {
	moduleStr := `
module test {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	%s
}
`
	yangStr := fmt.Sprintf(moduleStr, yangFragment)
	if module, err := yang.LoadModuleFromByteArray([]byte(yangStr), nil); err != nil {
		t.Error(err.Error())
	} else {
		s := &MySelection{}
		s.State.Meta = module
		s.State.Position = module.GetFirstMeta()
		return s
	}
	return nil
}

func tojson(t *testing.T, s *MySelection) string {
	var actual bytes.Buffer
	json := NewJsonWriter(&actual)
	out, _ := json.GetSelector()
	err := Insert(s, out, NewExhaustiveController())
	if err != nil {
		t.Error(err)
	}
	return string(actual.Bytes())
}

func TestYangBrowse(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0000-00-00 {
		description "x";
	}
	list hobbies {
		container birding {
			leaf favorite-species {
				type string;
			}
		}
		container hockey {
			leaf favorite-team {
				type string;
			}
		}
	}
}`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Error("bad module", err)
	} else {
		var actual bytes.Buffer
		json := NewJsonWriter(&actual)
		out, _ := json.GetSelector()
		metaTx := NewSchemaBrowser(module, false)
		in, err := metaTx.RootSelector()
		if err != nil {
			t.Error(err)
		}
		if err = Insert(in, out, NewExhaustiveController()); err != nil {
			t.Error("failed to transmit json", err)
		} else {
			t.Log("Round Trip:", string(actual.Bytes()))
		}
	}
}

func TestYangWrite(t *testing.T) {
	simple, err := yang.LoadModule(schema.NewCwdSource(), "../testdata/simple.yang")
	if err != nil {
		t.Error(err)
	} else {
		fromBrowser := NewSchemaBrowser(simple, false)
		from, _ := fromBrowser.RootSelector()
		toBrowser := NewSchemaBrowser(nil, false)
		to, _ := toBrowser.RootSelector()
		err = Insert(from, to, NewExhaustiveController())
		if err != nil {
			t.Error(err)
		} else {
			// dump original and clone to see if anything is missing
			var expected string
			var actual string
			expected, err = DumpModule(fromBrowser)
			if err != nil {
				t.Error(err)
			}
			actual, err = DumpModule(toBrowser)
			if err != nil {
				t.Error(err)
			}
			if actual != expected {
				t.Log("Expected")
				t.Log(expected)
				t.Log("Actual")
				t.Log(actual)
				t.Fail()
			}
		}
	}
}

func DumpModule(b *SchemaBrowser) (string, error) {
	var buff bytes.Buffer
	dumper := NewDumper(&buff)
	s, _ := b.RootSelector()
	ds, _ := dumper.GetSelector()
	err := Insert(s, ds, NewExhaustiveController())
	if err != nil {
		return "", err
	}
	return string(buff.Bytes()), nil
}

