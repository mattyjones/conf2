package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"schema"
	"strings"
	"db"
)

type StoreEntry struct {
	Path string
	Values bson.M
}

type Store struct {
	c *mgo.Collection
	entry StoreEntry
	selector interface{}
}

func NewStore(collection *mgo.Collection, rootPath string) *Store {
	store := &Store{
		c : collection,
		selector : bson.M{ "path" : rootPath},
	}
	store.entry.Path = rootPath
	return store
}

func (s *Store) HasValues(path string) bool {
	// TODO: performance - most efficient way? sort first?
	for k, _ := range s.entry.Values {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (s *Store) KeyList(path string, meta *schema.List) ([]string, error) {
	builder := db.NewKeyListBuilder(path)
	for k, _ := range s.entry.Values {
		builder.ParseKey(k)
	}
	return builder.List(), nil
}

func (s *Store) Value(key string, dataType *schema.DataType) (*browse.Value, error) {
	v, found := s.entry.Values[key]
	if found {
		return browse.SetValue(dataType, v)
	}
	return nil, nil
}

func (s *Store) SetValue(key string, v *browse.Value) error {
	s.entry.Values[key] = v.Value()
	return nil
}

func (s *Store) Load() (err error) {
	err = s.c.Find(s.selector).One(&s.entry)
	if err == mgo.ErrNotFound {
		s.entry.Values = make(bson.M, 10)
		return nil
	}
	return nil
}

func (s *Store) Save() (err error) {
	_, err = s.c.Upsert(s.selector, &s.entry)
	return err
}
