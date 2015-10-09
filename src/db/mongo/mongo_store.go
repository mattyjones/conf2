package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"strings"
)

type Store struct {
	c *mgo.Collection
	oid bson.ObjectId
}

func NewStore(collection *mgo.Collection, oid bson.ObjectId) *Store {
	return &Store{
		c : collection,
		oid : oid,
	}
}

func (s *Store) Load() (values map[string]*browse.Value, err error) {
	var results bson.M
	s.c.FindId(s.oid).One(&results)
	values = make(map[string]*browse.Value, len(results))
	for path, v := range results {
		values[path] = &Value{}
	}
	if err = q.One(&results); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Store) Clear(p *browse.Path) error {
	escPath := strings.Replace("/", "\\/", p.URL, -1)
	pattern := fmt.Sprint("/^", escPath, "/")
	selector := bson.M {
		"$regex" : pattern,
	}
	q := s.c.FindId(s.oid)
	q.Select(selector)
	var results bson.M
	if err = q.One(&results); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Store) Upsert(vals map[string]*browse.Value) (err error) {
	records := make(bson.M, len(vals))
	for path, val := range vals {
		records[path] = val.Value()
	}
	s.c.UpsertId(s.oid, records)

	return nil
}
