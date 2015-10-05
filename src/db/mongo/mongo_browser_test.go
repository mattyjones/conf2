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
	db := NewMongoBrowser(m, c)
	data := `{"bird":{"species":[{"name":"sparrow","origin":{"country":"US"}},{"name":"canary"}]}}`
	in := &browse.JsonReader{
		In : strings.NewReader(data),
		Meta : m,
	}
	err = browse.Upsert(browse.NewPath(""), in, db)
	if err != nil {
		t.Fatal(err)
	}

	var actualBytes bytes.Buffer
	out := browse.NewJsonFragmentWriter(&actualBytes)
	err = browse.Upsert(browse.NewPath("bird/species=sparrow/origin"), db, out)
	if err != nil {
		t.Fatal(err)
	}
	actual := actualBytes.String()
	expected := `{"country":"US"}`
	if actual != expected {
		t.Errorf("\nExpected:%s\n  Actual:%s", expected, actual)
	}
}