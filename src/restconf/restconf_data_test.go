package restconf

import (
	"schema/browse"
	"testing"
)

func TestRestconfBrowserMetaLoad(t *testing.T) {
	rc := &Service{restconfPath: "/restconf/"}
	rc.registrations = make(map[string]*registration, 5)
	b, err := NewData(rc)
	if err != nil {
		t.Error(err.Error())
	} else {
		s, err := b.Selector(browse.NewPath("modules"))
		if err != nil {
			t.Error(err.Error())
		} else if s == nil {
			t.Error("Could not find modules/module")
		}
	}
}
