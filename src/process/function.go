package process
import "errors"

type Function struct {
	Name string
	Arguments []Expression
}

func (f *Function) Eval(stack *Stack, table Table) (interface{}, error) {
	fPtr, found := stack.Funcs[f.Name]
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

