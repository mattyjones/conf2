package process
import (
	"data"
	"conf2"
	"fmt"
)

type Process struct {
	table Table
}

func NewProcess(sel *data.Selection) *Process {
	return &Process{
		table : &NodeTable{
			Corner: sel,
		},
	}
}

func (p *Process) Into(into *data.Selection) *Process {
	p.table = &Join{
		On : p.table,
		Into : &NodeTable{
			Corner: into,
			autoCreate: true,
		},
	}
	return p
}

func (p *Process) RunScript(script *Script) (err error) {
	return p.Run(map[string]*Script{script.Name : script}, script.Name)
}

func (p *Process) Run(scripts map[string]*Script, main string) (err error) {
	var found bool
	script, found := scripts[main]
	if !found {
		return conf2.NewErr(fmt.Sprintf("Script %s not found", main))
	}
	stack := &Stack{
		Scripts : scripts,
	}
	err = p.table.Next()
	for p.table.HasNext() && err == nil {
		err = script.Exec(stack, p.table)
		if err == nil {
			err = p.table.Next()
		}
	}
	return err
}
