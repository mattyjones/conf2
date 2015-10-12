package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"schema"
	"strings"
)


type Store struct {
	c *mgo.Collection
	oid bson.ObjectId
	values bson.M
}

type StoreIndex interface {
	Next(*browse.WalkState, meta *schema.List, key []*browse.Value, bool) (bool, error)

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

func (s *Store) KeyList(path string) ([]string, err error) {

}


func (s *Store) KeyList(path string) ([]string, err error) {
	// TODO: performance - most efficient way? sort first?
	keysSet := make(map[string]struct{}, 10)
	keyStart := len(path) + 1
	for k, _ := range s.values {
		if strings.HasPrefix(k, path) {
			keyEnd := strings.IndexRune(k[keyStart:], '/')
			if keyEnd < 0 {
				continue
			}
			key := k[keyStart:keyStart + keyEnd]
			keysSet[key] = struct{}{}
		}
	}
	keys := make([]string, len(keysSet))
	var i int
	for k, _ := range keysSet {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (s *Store) Value(key string, dataType *schema.DataType) *browse.Value {
	v, found := s.values[key]
	if found {
		return browse.SetValue(dataType, v), nil
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
	err = s.c.UpsertId(s.oid, s.values)
	return err
}
