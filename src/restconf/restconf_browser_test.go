package restconf
import (
	"testing"
	"schema/browse"
)

func TestRestMetaLoad(t *testing.T) {
	rc := &serviceImpl{restconfPath:"/restconf/"}
	rc.registrations = make(map[string]*registration, 5)
	b, err := NewBrowser(rc)
	if err != nil {
		t.Error(err.Error())
	} else {
		s, err := b.RootSelector()
		if err != nil {
			t.Error(err.Error())
		} else {
			p, _ := browse.NewPath("modules/module")
			browse.WalkPath(s, p)
		}
	}
}