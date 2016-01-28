package process
import (
	"fmt"
	"conf2"
)

type Select struct {
	parent CodeBlock
	On string
	Into string
	Script string
	operations []Operation
}

func (s *Select) Parent() CodeBlock {
	return s.parent
}

func (s *Select) SetParent(parent CodeBlock)  {
	s.parent = parent
}

func (s *Select) Operations() []Operation {
	return s.operations
}

func (s *Select) AddOperation(op Operation) {
	if s.operations == nil {
		s.operations = make([]Operation, 0, 1)
	}
	op.SetParent(s)
	s.operations = append(s.operations, op)
}

func (selct *Select) join(table Table) (*Join, error) {
	var on, into Table
	var err error
	if len(selct.On) > 0 {
		on, err = table.Select(selct.On, false)
		if err != nil {
			return nil, err
		}
		if on == nil {
			return nil, nil
		}
	} else {
		on = table
	}
	if len(selct.Into) > 0 {
		into, err = table.Select(selct.Into, true)
		if err != nil {
			return nil, err
		}
		if into == nil {
			return nil, conf2.NewErr("Unable to create node for " + selct.Into)
		}
	} else {
		into = table
	}

	return &Join{
		On: on,
		Into: into,
	}, nil
}

func (selct *Select) Exec(stack *Stack, table Table) error {
	join, err := selct.join(table)
	if join == nil {
		return nil
	}
	if err != nil {
		return err
	}
	sub := &Stack{
		Parent: stack,
	}
	var subscript *Script
	if len(selct.Script) > 0 {
		var found bool
		subscript, found = stack.Scripts[selct.Script]
		if !found {
			return conf2.NewErr(fmt.Sprintf("Script %s not found", selct.Script))
		}
	} else if len(selct.operations) == 0 {
		return nil
	}
	err = join.Next()
	for join.HasNext() && err == nil {
		if subscript != nil {
			err = subscript.Exec(sub, join)
		} else {
			err = execOperations(stack, selct.operations, join)
		}
		if err == nil {
			err = join.Next()
		}
	}
	return err
}
