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
			func(state *WalkState, meta schema.HasDataType) (*Value, error) {
				return &Value{
					Int : 0,
					Str : "a",
				}, nil
			},
    	},
		{
			`leaf-list c { type enumeration { enum a; enum b; } }`,
			`{"c":["a","b"]}`,
			func(state *WalkState, meta schema.HasDataType) (*Value, error) {
				return &Value{
					Intlist : []int{0, 1},
					Strlist : []string {"a", "b"},
				}, nil
			},
		},
	}
	for _, test := range tests {
		state, s := makeMySelection(t, test.yang)
		s.OnRead = test.read
		actual := tojson(t, state, s)
		if actual != test.expected {
			msg := fmt.Sprintf("Expected:\"%s\" Actual:\"%s\"", test.expected, actual)
			t.Error(msg)
		}
	}
}

func makeMySelection(t *testing.T, yangFragment string) (*WalkState, *MySelection) {
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
		state := NewWalkState(module)
		state.SetPosition(module.GetFirstMeta())
		return state, &MySelection{}
	}
	return nil, nil
}

func tojson(t *testing.T, state *WalkState, s *MySelection) string {
	var actual bytes.Buffer
	json := NewJsonWriter(&actual)
	out, _ := json.GetSelector()
	err := Insert(state, s, out, WalkAll())
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
		in, state, err := metaTx.RootSelector()
		if err != nil {
			t.Error(err)
		}
		if err = Insert(state, in, out, WalkAll()); err != nil {
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
		from, state, _ := fromBrowser.RootSelector()
		toBrowser := NewSchemaBrowser(nil, false)
		to, _, _ := toBrowser.RootSelector()
		err = Insert(state, from, to, WalkAll())
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
				t.Log("Different")
//				t.Log(expected)
//				t.Log("Actual")
//				t.Log(actual)
//				t.Fail()
			}
		}
	}
}

func DumpModule(b *SchemaBrowser) (string, error) {
	var buff bytes.Buffer
	dumper := NewDumper(&buff)
	s, state, _ := b.RootSelector()
	ds, _ := dumper.GetSelector()
	err := Insert(state, s, ds, WalkAll())
	if err != nil {
		return "", err
	}
	return string(buff.Bytes()), nil
}

