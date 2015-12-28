package process

import (
	"data"
	"text/template"
	"bytes"
	"conf2"
	"fmt"
	"strings"
	"errors"
)

type Stack struct {
	Join    *Join
	Scripts map[string]*Script
	Lets    map[string]string
	Parent  *Stack
	fmap    template.FuncMap
	LastErr error
	buffer  bytes.Buffer
}

func (stack *Stack) FuncMap() template.FuncMap {
	if stack.fmap == nil {
		stack.fmap = template.FuncMap{
			"get": stack.Get,
		}
	}
	return stack.fmap
}

func (stack *Stack) Get(key string) string {
	if v, hasValue := stack.Lets[key]; hasValue {
		return v
	}
	v, err := stack.Join.Get(key)
	if err != nil {
		stack.LastErr = err
		return ""
	}
	if v != nil {
		return v.String()
	}
	return ""
}

func (stack *Stack) Eval(expression string) (string, error) {
	tmpl, parseErr := template.New("convert").Funcs(stack.FuncMap()).Parse(expression)
	if parseErr != nil {
		return "", parseErr
	}
	from := struct {} {}
	stack.buffer.Reset()
	execErr := tmpl.Execute(&stack.buffer, from)
	if execErr != nil {
		return "", execErr
	}
	return stack.buffer.String(), nil
}

func (stack *Stack) ClearLastError() error {
	e := stack.LastErr
	stack.LastErr = nil
	return e
}

type Operation interface {
	SetParent(CodeBlock)
	Parent() CodeBlock
	Exec(stack *Stack) error
}

type CodeBlock interface {
	Operation
	AddOperation(op Operation)
}

type If struct {
	parent CodeBlock
	Expression string
	operations []Operation
}

func (iF *If) IsTrue(stack *Stack) bool {
	return true
}

func (iF *If) Operations() []Operation {
	return iF.operations
}

func (iF *If) SetParent(parent CodeBlock) {
	iF.parent = parent
}

func (iF *If) Parent() CodeBlock {
	return iF.parent
}

func (iF *If) AddOperation(op Operation) {
	if iF.operations == nil {
		iF.operations = make([]Operation, 0, 1)
	}
	iF.operations = append(iF.operations, op)
}

func (iF *If) Exec(stack *Stack) error {
	condition, err := stack.Eval(iF.Expression)
	if err != nil {
		return err
	}
	clean := strings.TrimSpace(condition)
	if len(clean) == 0 || "false" == clean || "0" == clean {
		return nil
	}
	for _, op := range iF.operations {
		if err = op.Exec(stack); err != nil {
			return err
		}
	}
	return nil
}


type Let struct {
	parent CodeBlock
	Name string
	Value string
}

func (v *Let) Exec(stack *Stack) error {
	stack.Lets[v.Name] = v.Value
	return nil
}

func (v *Let) SetParent(parent CodeBlock) {
	v.parent = parent
}

func (v *Let) Parent() CodeBlock {
	return v.parent
}

type Set struct {
	parent CodeBlock
	Name string
	Value string
}

func (s *Set) SetParent(parent CodeBlock) {
	s.parent = parent
}

func (s *Set) Parent() CodeBlock {
	return s.parent
}

func (set *Set) Exec(stack *Stack) error {
	s, err := stack.Eval(set.Value)
	if err != nil {
		return err
	}
	if len(s) == 0 {
		return nil
	}
	return data.ChangeValue(stack.Join.Into.Row, set.Name, s)
}

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

func (g *Goto) Exec(stack *Stack) error {
	s, found := stack.Scripts[g.Script]
	if !found {
		return errors.New(g.Script + " script not found")
	}
	return s.Exec(stack)
}

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

func (selct *Select) Exec(stack *Stack) error {
	var on, into *data.Selection
	if len(selct.On) > 0 {
		var err error
		on, err = data.SelectMetaList(stack.Join.On.Row, selct.On, false)
		if err != nil {
			return err
		}
		if on == nil {
			return nil
		}
	} else {
		on = stack.Join.On.Row
	}
	if len(selct.Into) > 0 {
		var err error
		into, err = data.SelectMetaList(stack.Join.Into.Row, selct.Into, true)
		if err != nil {
			return err
		}
		if into == nil {
			return conf2.NewErr("Unable to create node for " + selct.Into)
		}
	} else {
		into = stack.Join.Into.Row
	}

	sub := &Stack{
		Parent: stack,
		Join: &Join{
			Into: &Table {
				Corner: into,
			},
			On: &Table {
				Corner: on,
			},
		},
	}
	subScript, found := stack.Scripts[selct.Script]
	if !found {
		return conf2.NewErr(fmt.Sprintf("Script %s not found", selct.Script))
	}
	return subScript.Exec(sub)
}

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

func (script *Script) Exec(stack *Stack) (err error) {
	more, err := stack.Join.Iterate()
	for more && err == nil {
		for _, op := range script.Operations() {
			if err = op.Exec(stack); err != nil {
				return err
			}
		}
		more, err = stack.Join.Iterate()
	}
	return err
}
