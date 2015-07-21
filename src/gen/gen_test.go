package gen
import (
	"testing"
	"text/template"
	"bytes"
	"yang"
)

func TestGeneration(t *testing.T) {
	meta := `
module birds {
	namespace "ns";
	revision 99-99-9999 {
	  description "bingo";
	}
	list life-list {
		leaf species {
			type string;
		}
	}
	container list {
		leaf species {
			type string;
		}
	}
}
`
	if m, err := yang.LoadModuleFromByteArray([]byte(meta)); err == nil {
		foo := []string{"hello", "word"}
		data := struct {
			M *yang.Module
			Foo []string
		}{
			m, foo,
		}
		funcs := template.FuncMap {
			"iterate" : yang.ListToArray,
		}
		actual := new(bytes.Buffer)
		if tmpl, err := template.New("test").Funcs(funcs).Parse(`
public class DecoderFactory implements DecoderFactory {
		public DecoderFactory() {
{{ range iterate .M.DataDefs }}
			registerDecoder("{{.GetParent.GetIdent}}/{{.GetIdent}}", new {{.GetIdent}}Decoder());
{{ end }}
		}
}
		`); err == nil {
			err = tmpl.Execute(actual, data)
			t.Log(actual.String())
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}

}
