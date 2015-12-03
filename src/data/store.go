package data

import (
	"schema"
)

type Store interface {
	Load() error
	Save() error
	HasValues(path string) bool
	Value(path string, typ *schema.DataType) *schema.Value
	SetValue(path string, v *schema.Value) error
	KeyList(path string, meta *schema.List) ([]string, error)
	RenameKey(oldPath string, newPath string)
	Action(path string) (ActionFunc, error)
	RemoveAll(path string) error
}
