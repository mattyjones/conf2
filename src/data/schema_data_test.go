package data

import (
	"bytes"
	"fmt"
	"schema"
	"schema/yang"
	"testing"
)

func printMeta(m schema.Meta, level string) {
	fmt.Printf("%s%s\n", level, m.GetIdent())
	if nest, isNest := m.(schema.MetaList); isNest {
		if len(level) >= 16 {
			panic("Max level reached")
		}
		i2 := schema.NewMetaListIterator(nest, false)
		for i2.HasNextMeta() {
			printMeta(i2.NextMeta(), level+"  ")
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
	m, err := yang.LoadModuleFromByteArray([]byte(moduleStr), nil)
	if err != nil {
		t.Fatal("bad module", err)
	}
	b := NewSchemaData(m, false)
	var actual bytes.Buffer
	json := NewJsonWriter(&actual).Node()
	if err = NodeToNode(b.Node(), json, b.Schema()).Insert(); err != nil {
		t.Error(err)
	} else {
		t.Log("Round Trip:", string(actual.Bytes()))
	}
}

// TODO: support typedefs - simpleyang datatypes that use typedefs return format=0
func DISABLED_TestYangWrite(t *testing.T) {
	simple, err := yang.LoadModuleFromByteArray([]byte(yang.TestDataSimpleYang), nil)
	if err != nil {
		t.Fatal(err)
	}
	from := NewSchemaData(simple, false)
	to := NewSchemaData(nil, false)
	edit, err2 := PathToPath(from, to, "")
	if err2 != nil {
		t.Fatal(err2)
	}
	err = edit.Upsert()
	if err != nil {
		t.Fatal(err)
	}
	// dump original and clone to see if anything is missing
	diff := Diff(from.Node(), to.Node())
	var out bytes.Buffer
	diffOut := NewJsonWriter(&out).Node()
	NodeToNode(diff, diffOut, from.Schema()).Insert()
	t.Log(out.String())
//
//	var expected string
//	var actual string
//	expected, err = DumpModule(from)
//	if err != nil {
//		t.Error(err)
//	}
//	actual, err = DumpModule(to)
//	if err != nil {
//		t.Error(err)
//	}
//	if actual != expected {
//		t.Log("Different")
//						t.Log(expected)
//						t.Log("Actual")
//						t.Log(actual)
//						t.Fail()
//	}
}

func DumpModule(b *SchemaData) (string, error) {
	var buff bytes.Buffer
	err := NodeToNode(b.Node(), NewDumper(&buff).Node(), b.Schema()).Insert()
	if err != nil {
		return "", err
	}
	return string(buff.Bytes()), nil
}
