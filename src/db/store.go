package db
import (
	"schema/browse"
	"schema"
)

type Store interface {
	Load() error
	Save() error
	HasValues(path string) bool
	Value(path string, typ *schema.DataType) (*browse.Value, error)
	SetValue(path string, v *browse.Value) error
	KeyList(path string, meta *schema.List) ([]string, error)
	RenameKey(oldPath string, newPath string)
}