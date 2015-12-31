package process

type Let struct {
	parent CodeBlock
	Name string
	Expression Expression
}

func (l *Let) Exec(stack *Stack, table Table) (error) {
	v, err := l.Expression.Eval(stack, table)
	stack.Let(l.Name, v)
	return err
}

func (v *Let) SetParent(parent CodeBlock) {
	v.parent = parent
}

func (v *Let) Parent() CodeBlock {
	return v.parent
}

