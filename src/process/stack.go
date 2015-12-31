package process

type Stack struct {
	Scripts map[string]*Script
	Lets    map[string]interface{}
	Parent  *Stack
	Funcs   map[string]interface{}
	LastErr error
}

func (stack *Stack) Let(key string, value interface{}) {
	if stack.Lets == nil {
		stack.Lets = make(map[string]interface{}, 1)
	}
	stack.Lets[key] = value
}

func (stack *Stack) Func(key string) (f interface{}, found bool) {
	if stack.Funcs == nil {
		if stack.Parent == nil {
			stack.Funcs = builtins
		} else {
			return stack.Parent.Func(key)
		}
	}
	f, found = stack.Funcs[key]
	return
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

