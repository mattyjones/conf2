package browse
import (
	"regexp"
	"conf2"
	"fmt"
)

type Event int

const (
	NEW Event = iota + 1
	BEGIN_EDIT
	END_EDIT
	UNDO_EDIT
	LEAVE
	NEXT
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
	"NEXT",
	"DELETE",
	"UNSELECT",
}

type Events struct {
	listeners []*listener
}

type ListenFunc func() error

func (e Event) String() string {
	return eventNames[e]
}

type listener struct {
	path string
	regex *regexp.Regexp
	event   Event
	handler ListenFunc
}

func (l *listener) String() string {
	if len(l.path) > 0 {
		return fmt.Sprintf("%s:%s=>%p", l.event, l.path, l.handler)
	}
	return fmt.Sprintf("%s:%v=>%p", l.event, l.regex, l.handler)
}

func (impl *Events) dump() {
	for _, l := range impl.listeners {
		conf2.Debug.Print(l.String())
	}
}

func (impl *Events) AddByFullPath(e Event, path string, handler ListenFunc) {
	impl.listeners = append(impl.listeners, &listener{event: e, path: path, handler: handler})
}

func (impl *Events) AddByRegex(e Event, regex *regexp.Regexp, handler ListenFunc) {
	impl.listeners = append(impl.listeners, &listener{event: e, regex: regex, handler: handler})
}

func (impl *Events) Remove(handler ListenFunc) {
	for i, l := range impl.listeners {
		if &l.handler == &handler {
			impl.listeners = append(impl.listeners[:i], impl.listeners[i + 1:]...)
			break
		}
	}
}

func (impl *Events) Fire(path string, e Event) (err error) {
	if len(impl.listeners) > 0 {
		for _, l := range impl.listeners {
			if l.event == e {
				if len(l.path) > 0 {
					if l.path != path {
						continue
					}
				} else if l.regex != nil {
					if ! l.regex.MatchString(path) {
						continue
					}
				}
				if err = l.handler(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
