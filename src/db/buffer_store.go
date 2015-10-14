package db
import (
	"schema/browse"
	"schema"
	"strings"
)

// Store key values in memory.  Useful for testing or moving temporary data

type BufferStore map[string]*browse.Value

func (kvs BufferStore) Load() error {
	return nil
}

func (kvs BufferStore) Clear() error {
	for k, _ := range kvs {
		delete(kvs, k)
	}
	return nil
}

func (kvs BufferStore) HasValues(path string) bool {
	for k, _ := range kvs {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (kvs BufferStore) Save() error {
	return nil
}

func (kvs BufferStore) KeyList(key string, meta *schema.List) ([]string, error) {
	builder := NewKeyListBuilder(key)
	for k, _ := range kvs {
		builder.ParseKey(k)
	}
	return builder.List(), nil
}


func (kvs BufferStore) Value(key string, dataType *schema.DataType) (*browse.Value, error) {
	if v, found := kvs[key]; found {
		v.Type = dataType
		return v, nil
	}
	return nil, nil
}

func (kvs BufferStore) SetValue(key string, v *browse.Value) error {
	kvs[key] = v
	return nil
}

