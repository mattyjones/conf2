package process
import "testing"

func sayHello(stack *Stack, table Table, name string) string {
	return "hello " + name
}

func TestScriptFunc(t *testing.T) {
	s := &Stack{}
	s.Funcs = map[string]interface{} {
		"SayHello" : sayHello,
	}
	f := &Function{Name :"SayHello"}
	f.Arguments = []Expression {
		&Primative{Str: "charlie"},
	}
	actual, err := f.Eval(s, nil)
	if err != nil {
		t.Error(err)
	}
	expected := "hello charlie"
	if actual != expected {
		t.Errorf("Expected:%s\n  Actual:%s\n", expected, actual)
	}
}
