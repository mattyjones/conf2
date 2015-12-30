package process
import "errors"

type Goto struct {
	parent CodeBlock
	Script string
}

func (g *Goto) Parent() CodeBlock {
	return g.parent
}

func (g *Goto) SetParent(parent CodeBlock)  {
	g.parent = parent
}

func (g *Goto) Exec(stack *Stack, table Table) error {
	s, found := stack.Scripts[g.Script]
	if !found {
		return errors.New(g.Script + " script not found")
	}
	return s.Exec(stack, table)
}
