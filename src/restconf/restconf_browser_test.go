package restconf
import (
	"testing"
	"schema/browse"
)

func TestRestconfBrowserMetaLoad(t *testing.T) {
	rc := &serviceImpl{restconfPath:"/restconf/"}
	rc.registrations = make(map[string]*registration, 5)
	b, err := NewBrowser(rc)
	if err != nil {
		t.Error(err.Error())
	} else {
		s, _, err := b.Selector(browse.NewPath("modules"), browse.READ)
		if err != nil {
			t.Error(err.Error())
		} else if s == nil {
			t.Error("Could not find modules/module")
		}
	}
}