package mongo
import (
	"testing"
	"schema/yang"
	"gopkg.in/mgo.v2"
	"schema"
	"schema/browse"
	"strings"
	"bytes"
)

func testCollection(cname string) *mgo.Collection {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	c := session.DB("test").C(cname)
	c.DropCollection()
	return c
}


func TestMongoBrowserWrite(t *testing.T) {
	mstr := `
module animal {
	prefix "";
	namespace "";
	revision 0;
	container bird {
		list species {
			key "name";
			leaf name {
				type string;
			}
			container origin {
				leaf country {
					type string;
				}
			}
		}
	}
}
`
	var m *schema.Module
	var err error
	m, err = yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	c := testCollection("TestMongoBrowserWrite")
	var b *MongoBrowser
	b = NewMongoBrowser(m, c)
	data := `{"bird":{"species":[{"name":"sparrow","origin":{"country":"US"}},{"name":"canary"}]}}`
	var dataBrowser *browse.JsonReader
	dataBrowser = browse.NewJsonReader(strings.NewReader(data))
	var w browse.Selection
	var state *browse.WalkState
	w, state, err = b.WriteSelector(browse.NewPath(""), browse.UPSERT)
	if err != nil {
		t.Fatal(err)
	}
	r, _ := dataBrowser.GetSelector(state)

	err = browse.Upsert(state, r, w, browse.WalkAll())
	if err != nil {
		t.Fatal(err)
	}

	var read browse.Selection
	if read, state, err = b.ReadSelector(browse.NewPath("bird/species=sparrow/origin")); err != nil {
		t.Fatal(err)
	} else if read == nil {
		t.Fatal("Could not find item")
	}
	state.InsideList()
	var actualBytes bytes.Buffer
	jsonOut := browse.NewJsonWriter(&actualBytes)
	out, _ := jsonOut.GetSelector()
	browse.Insert(state, read, out, browse.WalkAll())
	actual := actualBytes.String()
	expected := `{"country":"US"}`
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}