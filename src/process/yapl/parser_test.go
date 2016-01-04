package yapl
import "testing"

func TestYaplParser(t *testing.T) {
	//yyDebug = 4
	scripts, err := Load(
`foo
  a = b
  let a = b
  select z into q
    x = x(z)
    y = w()
  if p
    x = x(z,b,c(z))
    goto bleep

bleep
    z = s
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
	if _, found := scripts["foo"]; !found {
		t.Error("'foo' script not found")
	}
}
