package process
import (
	"errors"
	"bytes"
	"fmt"
)

type Function struct {
	Name string
	Arguments []Expression
}

func (f *Function) Eval(stack *Stack, table Table) (interface{}, error) {
	fPtr, found := stack.Func(f.Name)
	if !found {
		return nil, errors.New("Function not found " + f.Name)
	}
	argumentValues := make([]interface{}, len(f.Arguments))
	var err error
	for i, arg := range f.Arguments {
		if argumentValues[i], err = arg.Eval(stack, table); err != nil {
			return nil, err
		}
	}
	return call(fPtr, argumentValues)
}

var builtins map[string]interface{}
func init() {
	builtins = map[string]interface{} {
		"concat" : func(args ...interface{}) string {
			var buff bytes.Buffer
			for _, arg := range args {
				if tostr, ok := arg.(fmt.Stringer); ok {
					buff.WriteString(tostr.String())
				} else {
					buff.WriteString(fmt.Sprintf("%v", arg))
				}
			}
			return buff.String()
		},
	}
}
