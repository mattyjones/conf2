package browse
import (
	"log"
	"yang"
)

type DebuggingWriter struct {
	Delegate Writer
}

func (w *DebuggingWriter) EnterContainer(m yang.MetaList) error {
	log.Println("Entering Container", m.GetIdent())
	return w.Delegate.EnterContainer(m)
}

func (w *DebuggingWriter) ExitContainer(m yang.MetaList) error {
	log.Println("Exiting Container", m.GetIdent())
	return w.Delegate.ExitContainer(m)
}

func (w *DebuggingWriter) EnterList(m *yang.List) error {
	log.Println("Entering List", m.GetIdent())
	return w.Delegate.EnterList(m)
}

func (w *DebuggingWriter) ExitList(m *yang.List) error {
	log.Println("Exiting List", m.GetIdent())
	return w.Delegate.ExitList(m)
}

func (w *DebuggingWriter) UpdateValue(m yang.Meta, v *Value) error {
	log.Println("Updating Value", m.GetIdent())
	return w.Delegate.UpdateValue(m, v)
}

type NullWriter struct {
}

func (NullWriter) EnterContainer(yang.MetaList) error {
	return nil
}

func (NullWriter) ExitContainer(yang.MetaList) error {
	return nil
}

func (NullWriter) EnterList(*yang.List) error {
	return nil
}

func (NullWriter) ExitList(*yang.List) error {
	return nil
}

func (NullWriter) EnterListItem(*yang.List) error {
	return nil
}

func (NullWriter) ExitListItem(*yang.List) error {
	return nil
}

func (NullWriter) UpdateValue(meta yang.Meta, val *Value) error {
	return nil
}