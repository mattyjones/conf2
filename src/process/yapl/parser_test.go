package yapl
import "testing"

func TestYaplParser(t *testing.T) {
	//yyDebug = 3
	scripts, err := Load(
`foo
  a = b
  let a = b
  select z into q
    x = x
  if p
    z = "y"
    goto bleep

bleep
  glop = blop
`)
	if err != nil {
		t.Error(err)
	}
	if scripts == nil {
		t.Fatal("Nil script")
	}
	if len(scripts) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(scripts))
	}
	if scripts[0].Name != "foo" {
		t.Error("Incorrect script name " + scripts[0].Name)
	}
}
