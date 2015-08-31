package browse
import (
	"testing"
	"yang"
	"bytes"
	"os"
	"fmt"
)
func TestScratch(t *testing.T) {
	cwd, _ := os.Getwd()
	t.Log(cwd)
}

func printMeta(m yang.Meta, level string) {
	fmt.Printf("%s%s\n", level, m.GetIdent())
	if nest, isNest := m.(yang.MetaList); isNest {
		if len(level) >= 16 {
			panic("Max level reached")
		}
		i2 := yang.NewMetaListIterator(nest, false)
		for i2.HasNextMeta() {
			printMeta(i2.NextMeta(), level + "  ")
		}
	}
}

func TestYangMeta(t *testing.T) {
	ds := &yang.FileStreamSource{Root:"../../../etc"}
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
}
	`
	if module, err := yang.LoadModuleFromByteArray([]byte(moduleStr)); err != nil {
		t.Error("bad module", err)
	} else {
		ds := &yang.FileStreamSource{Root:"../../../etc"}
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
			if err = Insert(in, out); err != nil {
				t.Error("failed to transmit json", err)
			} else {
				t.Log("Round Trip:", string(actual.Bytes()))
			}
		}
	}
}
