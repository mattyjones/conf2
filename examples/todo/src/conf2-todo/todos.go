package main
import (
	"time"
	"schema/browse"
	"schema"
)

type Todos struct {
	Tasks []*Task
}

func (todos *Todos) Manage() browse.Node {
	s := &browse.MyNode{}
	var index int
	s.OnNext = func(sel *browse.Selection, meta *schema.List, new bool, key []*browse.Value, first bool) (browse.Node, error) {
		var task *Task
		if len(key) > 0 {
			task = todos.Tasks[key[0].Int]
		} else if new {
			task = Task{}
			if todos.Tasks == nil {
				todos.Tasks = []Task{task}
			} else {
				todos.Tasks = append(todos.Tasks, task)
			}
		} else {
			if first {
				index = 0
			}
			if index < len(todos.Tasks) {
				task = todos.Tasks[index]
				sel.State.SetKey(&browse.Value{Int:index})
			}
		}
		if task != nil {
			return task.Select(), nil
		}
		return nil, nil
	}
	return s
}

type Task struct {
	Summary string
	Description string
	DueDate time.Duration
	Keywords []string
	timer *time.Timer
}

func (task *Task) Select() browse.Node {
	s := &browse.MyNode{}
	s.Read = func(sel *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "dueDate":
			return &browse.Value{Int64:task.DueDate}
		}
		return browse.ReadField(meta, task)
	}
	s.Write = func(sel *browse.Selection, meta schema.HasDataType, v *browse.Value) error {
		switch meta.GetIdent() {
		case "dueDate":
			task.DueDate = v.Int64
			if task.timer != nil {
				task.timer.Reset(task.DueDate)
			}
		}
		return nil
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) error {
		switch s {
		case browse.NEW:
			task.timer = time.NewTimer(task.DueDate)
		case browse.DELETE:
			if task.timer != nil {
				task.timer.Stop()
			}
		}
		return
	}

	return s
}