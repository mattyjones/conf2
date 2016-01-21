package data

import (
	"schema"
	"schema/yang"
	"testing"
)

func TestDiff(t *testing.T) {
	moduleStr := `
module m {
	prefix "";
	namespace "";
	revision 0;
	container movie {
	    leaf name {
	      type string;
	    }
		container character {
		    leaf name {
		      type string;
		    }
		}
	}
	container car {
		leaf name {
			type string;
		}
		leaf year {
			type int32;
		}
	}
	container videoGame {
		leaf name {
			type string;
		}
	}
}
	`
	var err error
	var m *schema.Module
	if m, err = yang.LoadModuleFromByteArray([]byte(moduleStr), nil); err != nil {
		t.Fatal(err)
	}

	str := schema.NewDataType(nil, "string")

	// new
	a := NewBufferStore()
	a.Values["movie/name"] = &Value{Str: "StarWars"}
	a.Values["movie/character/name"] = &Value{Str: "Hans Solo"}
	a.Values["car/name"] = &Value{Str: "Malibu"}
	aData := NewStoreData(m, a).Node()

	// old
	b := NewBufferStore()
	b.Values["movie/name"] = &Value{Str: "StarWars"}
	laya := &Value{Type: str, Str: "Princess Laya"}
	b.Values["movie/character/name"] = laya
	gtav := &Value{Type: str, Str: "GTA V"}
	b.Values["videoGame/name"] = gtav
	bData := NewStoreData(m, b).Node()

	c := NewBufferStore()
	if err = NewStoreData(m, c).Select().Pull(Diff(bData, aData)).Insert(); err != nil {
		t.Error(err)
	}
	if len(c.Values) != 2 {
		t.Error("Expected 1 value")
	}
	if !laya.Equal(c.Value("movie/character/name", str)) {
		t.Errorf("Unexpected values %v", c.Values)
	}
	if !gtav.Equal(c.Value("videoGame/name", str)) {
		t.Errorf("Unexpected values %v", c.Values)
	}
}
