package process

type Operation interface {
	SetParent(CodeBlock)
	Parent() CodeBlock
	Exec(stack *Stack, table Table) error
}

type CodeBlock interface {
	Operation
	AddOperation(op Operation)
}

// Simply a named code block,and entry point of execution
type Script struct {
	Name string
	operations []Operation
}

func (s *Script) Parent() CodeBlock {
	return nil
}

func (s *Script) SetParent(parent CodeBlock)  {
	panic("Cannot set parent on script")
}

func (s *Script) Operations() []Operation {
	return s.operations
}

func (s *Script) AddOperation(op Operation) {
	if s.operations == nil {
		s.operations = make([]Operation, 0, 1)
	}
	s.operations = append(s.operations, op)
}

func execOperations(stack *Stack, ops []Operation, table Table) (err error) {
	for _, op := range ops {
		if err = op.Exec(stack, table); err != nil {
			return err
		}
	}
	return
}

func (script *Script) Exec(stack *Stack, table Table) (err error) {
	return execOperations(stack, script.operations, table)
}


