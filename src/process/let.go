package process

type Let struct {
	parent CodeBlock
	Name string
	Expression Expression
}

func (v *Let) Exec(stack *Stack, table Table) (err error) {
	stack.Lets[v.Name], err = v.Expression.Eval(stack, table)
	return err
}

func (v *Let) SetParent(parent CodeBlock) {
	v.parent = parent
}

func (v *Let) Parent() CodeBlock {
	return v.parent
}

