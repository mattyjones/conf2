package yapl
import "testing"

func TestYaplParser(t *testing.T) {
	//yyDebug = 3
	scripts, err := Load(
`foo
  a = b
  let a = b
`)
	if err != nil {
		t.Error(err)
	}
	if scripts == nil {
		t.Fatal("Nil script")
	}
	if scripts[0].Name != "foo" {
		t.Error("Incorrect script name " + scripts[0].Name)
	}
}
