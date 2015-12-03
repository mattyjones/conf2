package restconf

import (
	"data"
	"testing"
)

func TestRestconfBrowserMetaLoad(t *testing.T) {
	rc := NewService()
	b, err := NewData(rc)
	if err != nil {
		t.Fatal(err)
	}
	if err = rc.RegisterBrowser(b); err != nil {
		t.Fatal(err)
	}
	s, err := data.WalkDataPath(b, "modules=restconf/module")
	if err != nil {
		t.Error(err.Error())
	} else if s == nil {
		t.Error("Could not find modules/module")
	}
}
