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

func TestYangMeta(t *testing.T) {
	ds := &schema.FileStreamSource{Root:"../../../etc"}
	if yangModule, err := yang.LoadModule(ds, "yang-1.0.yang"); err != nil {
		t.Error("yang module", err)
	} else {
		printMeta(yangModule, "")
	}
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
		ds := &schema.FileStreamSource{Root:"../../../etc"}
		if yangModule, err := yang.LoadModule(ds, "yang-1.0.yang"); err != nil {
			t.Error("yang module", err)
		} else {
			var actual bytes.Buffer
			json := NewJsonWriter(&actual)
			out, _ := json.GetSelector()
			metaTx := &YangBrowser{meta:yangModule, module:module}
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
}

//func TestYangDump(t *testing.T) {
//	yang := NewYangBrowser(getYangModule())
//	var actual bytes.Buffer
//	dumper := NewDumper(&actual)
//	var err error
//	var in *MySelection
//	var out *MySelection
//	if in, err = schema.RootSelector(); err != nil {
//		t.Error("failed to dump yang", err)
//	}
//	if out, err = dumper.GetSelector(); err != nil {
//		t.Error("failed to dump yang", err)
//	}
//	if err = Insert(in, out); err != nil {
//		t.Error("failed to dump yang", err)
//	} else {
//		fmt.Print(string(actual.Bytes()))
//	}
//}
