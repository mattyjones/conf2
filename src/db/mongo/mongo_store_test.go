package mongo
import (
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoStore(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	col := session.DB("test").C("TestMongoStore")
	store := NewStore(col, "abc")
	store.entry.Values = bson.M {
		"a/b/f" : "hi",
		"a/b/d" : "xxx",
	}
	err = store.Save()
	if err != nil {
		t.Error(err)
	}
	store.entry.Values["a/b/e"] = "zzz"
	err = store.Save()
	if err != nil {
		t.Error(err)
	}

	dup := NewStore(col, "abc")
	dup.Load()
	if len(dup.entry.Values) != 3 {
		t.Error("Expected 3 values got ", len(dup.entry.Values))
	}
}
