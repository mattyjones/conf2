package data
import (
	"regexp"
	"conf2"
	"fmt"
	"schema"
)

type Event int

const (
	NEW Event = iota + 1
	BEGIN_EDIT
	END_EDIT
	UNDO_EDIT
	LEAVE
	DELETE
	UNSELECT
)
var eventNames = []string {
	"N/A",
	"NEW",
	"BEGIN_EDIT",
	"END_EDIT",
	"UNDO_EDIT",
	"LEAVE",
	"DELETE",
	"UNSELECT",
}

type Events interface {
	AddListener(*Listener)
	RemoveListener(*Listener)
	Fire(path *schema.Path, e Event) error
}

type EventsImpl struct {
	Parent    Events
	listeners []*Listener
}

type ListenFunc func() error

func (e Event) String() string {
	return eventNames[e]
}

type Listener struct {
	path string
	regex *regexp.Regexp
	event   Event
	handler ListenFunc
}

func (l *Listener) String() string {
	if len(l.path) > 0 {
		return fmt.Sprintf("%s:%s=>%p", l.event, l.path, l.handler)
	}
	return fmt.Sprintf("%s:%v=>%p", l.event, l.regex, l.handler)
}

func (impl *EventsImpl) dump() {
	for _, l := range impl.listeners {
		conf2.Debug.Print(l.String())
	}
}

func (impl *EventsImpl) AddListener(l *Listener) {
	impl.listeners = append(impl.listeners, l)
}

func (impl *EventsImpl) RemoveListener(l *Listener) {
	for i, candidate := range impl.listeners {
		if l == candidate {
			impl.listeners = append(impl.listeners[:i], impl.listeners[i + 1:]...)
			break
		}
	}
}

func (impl *EventsImpl) Fire(path *schema.Path, e Event) (err error) {
	if len(impl.listeners) > 0 {
		pathStr := path.String()
		for _, l := range impl.listeners {
			if l.event == e {
				if len(l.path) > 0 {
					if l.path != pathStr {
						continue
					}
				} else if l.regex != nil {
					if ! l.regex.MatchString(pathStr) {
						continue
					}
				}
				if err = l.handler(); err != nil {
					return err
				}
			}
		}
	}
	if impl.Parent != nil {
		if err = impl.Parent.Fire(path, e); err != nil {
			return err
		}
	}
	return nil
}
