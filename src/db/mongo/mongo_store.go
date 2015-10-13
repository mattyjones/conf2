package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"schema"
	"strings"
	"db"
)


type Store struct {
	c *mgo.Collection
	oid bson.ObjectId
	values bson.M
}

func NewStore(collection *mgo.Collection, oid bson.ObjectId) *Store {
	return &Store{
		c : collection,
		oid : oid,
	}
}

func (s *Store) HasValues(path string) bool {
	// TODO: performance - most efficient way? sort first?
	for k, _ := range s.values {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (s *Store) KeyList(path string, meta *schema.List) ([]string, error) {
	builder := db.NewKeyListBuilder(path)
	for k, _ := range s.values {
		builder.ParseKey(k)
	}
	return builder.List(), nil
}

func (s *Store) Value(key string, dataType *schema.DataType) (*browse.Value, error) {
	v, found := s.values[key]
	if found {
		return browse.SetValue(dataType, v)
	}
	return nil, nil
}

func (s *Store) SetValue(key string, v *browse.Value) error {
	s.values[key] = v.Value()
	return nil
}

func (s *Store) Load() (err error) {
	err = s.c.FindId(s.oid).One(&s.values)
	return err
}

func (s *Store) Save() (err error) {
	_, err = s.c.UpsertId(s.oid, s.values)
	return err
}
