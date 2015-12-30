package process

type Set struct {
	parent CodeBlock
	Name string
	Expression Expression
}

func (s *Set) SetParent(parent CodeBlock) {
	s.parent = parent
}

func (s *Set) Parent() CodeBlock {
	return s.parent
}

func (set *Set) Exec(stack *Stack, table Table) error {
	v, err := set.Expression.Eval(stack, table)
	if err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	return table.Set(set.Name, v)
}

