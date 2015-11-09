package browse
import (
	"schema"
	"strings"
	"fmt"
)

// Store key values in memory.  Useful for testing or moving temporary data

type BufferStore struct {
	Values map[string]*Value
	Actions map[string]ActionFunc
}

func NewBufferStore() *BufferStore {
	return &BufferStore{
		Values : make(map[string]*Value, 10),
		Actions : make(map[string]ActionFunc, 10),
	}
}

func (kvs *BufferStore) Load() error {
	return nil
}

func (kvs *BufferStore) Clear() error {
	for k, _ := range kvs.Values {
		delete(kvs.Values, k)
	}
	return nil
}

func (kvs *BufferStore) HasValues(path string) bool {
	for k, _ := range kvs.Values {
		if strings.HasPrefix(k, path) {
			return true
		}
	}
	return false
}

func (kvs *BufferStore) Save() error {
	return nil
}

func (kvs *BufferStore) KeyList(key string, meta *schema.List) ([]string, error) {
	builder := NewKeyListBuilder(key)
	for k, _ := range kvs.Values {
		builder.ParseKey(k)
	}
	return builder.List(), nil
}

func (kvs *BufferStore) Action(key string) (ActionFunc, error) {
	return kvs.Actions[key], nil
}

func (kvs *BufferStore) Value(key string, dataType *schema.DataType) (*Value) {
	if v, found := kvs.Values[key]; found {
		v.Type = dataType
		return v
	}
	return nil
}

func (kvs *BufferStore) SetValue(key string, v *Value) error {
	kvs.Values[key] = v
	return nil
}

func (kvs *BufferStore) RemoveAll(path string) error {
	for k, _ := range kvs.Values {
		if strings.HasPrefix(k, path) {
			delete(kvs.Values, k)
		}
	}
	for k, _ := range kvs.Actions {
		if strings.HasPrefix(k, path) {
			delete(kvs.Actions, k)
		}
	}
	return nil
}

func (kvs *BufferStore) RenameKey(oldPath string, newPath string) {
	for k, v := range kvs.Values {
		if strings.HasPrefix(k, oldPath) {
			newKey := fmt.Sprint(newPath, k[len(oldPath):])
			delete(kvs.Values, k)
			kvs.Values[newKey] = v
		}
	}
	for k, v := range kvs.Actions {
		if strings.HasPrefix(k, oldPath) {
			newKey := fmt.Sprint(newPath, k[len(oldPath):])
			delete(kvs.Actions, k)
			kvs.Actions[newKey] = v
		}
	}
}

