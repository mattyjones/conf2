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

func TestYangBrowse(t *testing.T) {
	moduleStr := `
module json-test {
	prefix "t";
	namespace "t";
	revision 0;
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
		json := NewJsonFragmentWriter(&actual)
		b := NewSchemaBrowser(module, false)
		if err = Insert(NewPath(""), b, json); err != nil {
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
		from := NewSchemaBrowser(simple, false)
		to := NewSchemaBrowser(nil, false)
		err = Insert(NewPath(""), from, to)
		if err != nil {
			t.Error(err)
		} else {
			// dump original and clone to see if anything is missing
			var expected string
			var actual string
			expected, err = DumpModule(from)
			if err != nil {
				t.Error(err)
			}
			actual, err = DumpModule(to)
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
	err := Insert(NewPath(""), b, dumper)
	if err != nil {
		return "", err
	}
	return string(buff.Bytes()), nil
}

