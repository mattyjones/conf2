package yang

import (
	"os"
	"fmt"
)

type Resource interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type ResourceSource interface {
	OpenResource(resourceId string) (Resource, error)
}

type FileDataSource struct {
	Root string
}

func (src *FileDataSource) OpenResource(resourceId string) (Resource, error) {
	path := fmt.Sprint(src.Root, "/", resourceId)
	return os.Open(path)
}

type FsError struct {
	Msg string
}

func (e *FsError) Error() string {
	return e.Msg
}

