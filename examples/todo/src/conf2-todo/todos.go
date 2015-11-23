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
			index = key[0].Int
			task = todos.Tasks[index]
		} else if new {
			task = &Task{}
			index = len(todos.Tasks)
			if todos.Tasks == nil {
				todos.Tasks = []*Task{task}
			} else {
				todos.Tasks = append(todos.Tasks, task)
			}
		} else {
			if first {
				index = 0
			} else {
				index++
			}
			if index < len(todos.Tasks) {
				task = todos.Tasks[index]
				keyMeta := meta.KeyMeta()[0]
				sel.State.SetKey([]*browse.Value{&browse.Value{Type:keyMeta.GetDataType(), Int:index}})
			}
		}
		if task != nil {
			return task.Select(index), nil
		}
		return nil, nil
	}
	return s
}

type Status int
const (
	StatusTodo Status = iota
	StatusDone
)

type Task struct {
	Summary string
	Status Status
	Description string
	DueDate time.Duration
	Keywords []string
	timer *time.Timer
}

func (task *Task) Select(id int) browse.Node {
	s := &browse.MyNode{}
	s.OnRead = func(sel *browse.Selection, meta schema.HasDataType) (*browse.Value, error) {
		switch meta.GetIdent() {
		case "id":
			return  &browse.Value{Int:id}, nil
		case "dueDate":
			return &browse.Value{Int64:int64(task.DueDate)}, nil
		}
		return browse.ReadField(meta, task)
	}
	s.OnWrite = func(sel *browse.Selection, meta schema.HasDataType, v *browse.Value) error {
		switch meta.GetIdent() {
		case "id":
			// Not allowed
		case "dueDate":
			task.DueDate = time.Duration(v.Int64)
			if task.timer != nil {
				task.timer.Reset(task.DueDate)
			}
		default:
			return browse.WriteField(meta, task, v)
		}
		return nil
	}
	s.OnEvent = func(sel *browse.Selection, e browse.Event) error {
		switch e {
// This is what i want to change timers after all fields have been updated
//		case browse.UPDATE:
//
		case browse.NEW:
			if task.Status != StatusDone {
				task.timer = time.NewTimer(task.DueDate)
			}
		case browse.DELETE:
			if task.timer != nil {
				task.timer.Stop()
			}
		}
		return nil
	}

	return s
}