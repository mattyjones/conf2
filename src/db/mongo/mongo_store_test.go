package mongo
import (
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"schema/browse"
	"schema"
)

func TestMongoStore(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	col := session.DB("test").C("TestMongoStore")
	oid := bson.ObjectIdHex("5618045b36846627ce000001")
	store := NewStore(col, oid)
	strType := &schema.DataType{Format:schema.FMT_STRING}
	vals := map[string]*browse.Value {
		"a/b/f" : &browse.Value{Type:strType, Str:"hi"},
		"a/b/d" : &browse.Value{Type:strType, Str:"xxx"},
	}
	err = store.Upsert(browse.NewPath("/x/y"), vals)
	if err != nil {
		t.Error(err)
	}


}
