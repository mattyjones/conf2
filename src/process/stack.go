package process

type Stack struct {
	Scripts map[string]*Script
	Lets    map[string]interface{}
	Parent  *Stack
	Funcs   map[string]interface{}
	LastErr error
}

func (stack *Stack) ResolveValue(valname string, table Table) interface{} {
	if v, hasValue := stack.Lets[valname]; hasValue {
		return v
	}
	v, err := table.Get(valname)
	if err != nil {
		stack.LastErr = err
	}
	return v
}

func (stack *Stack) ClearLastError() error {
	e := stack.LastErr
	stack.LastErr = nil
	return e
}

