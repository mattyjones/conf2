package data

import (
	"bytes"
	"fmt"
	"schema"
	"schema/yang"
	"testing"
	"os"
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
	b := NewSchemaData(m, false).Select()
	var actual bytes.Buffer
	if err = b.Selector().Push(NewJsonWriter(&actual).Node()).Insert().LastErr; err != nil {
		t.Error(err)
	} else {
		t.Log("Round Trip:", string(actual.Bytes()))
	}
}

// TODO: support typedefs - simpleyang datatypes that use typedefs return format=0
func TestYangWrite(t *testing.T) {
	simple, err := yang.LoadModuleFromByteArray([]byte(yang.TestDataSimpleYang), nil)
	if err != nil {
		t.Fatal(err)
	}
	from := NewSchemaData(simple, false)
	to := NewSchemaData(nil, false)
	err = from.Select().Selector().Push(to.Node()).Upsert().LastErr
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout.WriteString("\n*********** O R I G I N A L **********\n")
	orig, _ := os.Create("original.json")
	defer orig.Close()
	from.Select().Selector().Push(NewJsonWriter(orig).Node()).Insert()

	os.Stdout.WriteString("\n*********** C O P Y **********\n")
	new, _ := os.Create("new.json")
	defer new.Close()
	to.Select().Selector().Push(NewJsonWriter(new).Node()).Insert()

	// dump original and clone to see if anything is missing
	diff := Diff(from.Node(), to.Node())
	diffSel := from.Select().Fork(diff)
	var out bytes.Buffer
	diffSel.Selector().Push(NewJsonWriter(&out).Node()).Insert()
	t.Log(out.String())
}
