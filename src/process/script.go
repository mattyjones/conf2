package process

import (
	"data"
	"text/template"
	"bytes"
	"conf2"
	"fmt"
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

type Operation struct {
	If *If
	Command Command
}

type Command interface {
	Exec(stack *Stack) error
}

type If struct {
	Expression string
	Operations []*Operation
}

func (iF *If) IsTrue(stack *Stack) bool {
	return true
}

type Let struct {
	Name string
	Value string
}

func (v *Let) Exec(stack *Stack) error {
	stack.Lets[v.Name] = v.Value
	return nil
}

type Set struct {
	Name string
	Value string
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

type Select struct {
	On string
	Into string
	Script string
	Operations []*Operation
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
	Operations []*Operation
}

func (script *Script) Exec(stack *Stack) (err error) {
	more, err := stack.Join.Iterate()
	for more && err == nil {
		for _, op := range script.Operations {
			if op.If != nil {
				if ! op.If.IsTrue(stack) {
					continue
				}
			}
			if err = op.Command.Exec(stack); err != nil {
				return err
			}
		}
		more, err = stack.Join.Iterate()
	}
	return err
}
