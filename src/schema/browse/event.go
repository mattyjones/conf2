package browse
import "conf2"

type EventHandler func(sel *Selection, e Event) error

type Event int

const (
	INIT Event = iota
	BEGIN_EDIT
	END_EDIT
	UNDO_EDIT
)

type Events struct {
	listeners []*listener
}

type listener struct {
	event   Event
	handler EventHandler
}

func (impl *Events) StartEdit(sel *Selection) (err error) {
	sel.Node().Event(sel, BEGIN_EDIT)
	impl.fire(sel, BEGIN_EDIT)
	return
}

func (impl *Events) Listen(e Event, handler EventHandler) {
	listener := &listener{handler: handler, event: e}
	impl.listeners = append(impl.listeners, listener)
}

func (impl *Events) RemoveListener(handler EventHandler) {
	// TODO
}

func (impl *Events) fire(sel *Selection, e Event) (err error) {
	if len(impl.listeners) == 0 {
		return
	}
	for _, l := range impl.listeners {
		if l.event == e {
			if err = l.handler(sel, e); err != nil {
				return err
			}
		}
	}
	return
}

func (impl *Events) EndEdit(sel *Selection) (err error) {
	sel.Node().Event(sel, END_EDIT)
	impl.fire(sel, END_EDIT)
	return
}

func (impl *Events) UndoEdit(sel *Selection) (err error) {
	sel.Node().Event(sel, UNDO_EDIT)
	impl.fire(sel, UNDO_EDIT)
	return
}
