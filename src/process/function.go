package process
import (
	"errors"
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
	argumentValues := make([]interface{}, len(f.Arguments) + 2)
	var err error
	argumentValues[0] = stack
	argumentValues[1] = table
	for i, arg := range f.Arguments {
		if argumentValues[i + 2], err = arg.Eval(stack, table); err != nil {
			return nil, err
		}
	}
	return call(fPtr, argumentValues)
}

