package process

type Expression interface {
	Eval(stack *Stack, table Table) (interface{}, error)
}

type Primative struct {
	Str string
	Num int
	Var string
}

func (p *Primative) Eval(stack *Stack, table Table) (interface{}, error) {
	if len(p.Str) > 0 {
		return p.Str, nil
	}
	if len(p.Var) > 0 {
		return stack.ResolveValue(p.Var, table), nil
	}
	return p.Num, nil
}

