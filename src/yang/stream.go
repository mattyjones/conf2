package yang

import (
	"os"
	"fmt"
)

type DataStream interface {
	Resource
	Read(p []byte) (n int, err error)
}

type StreamSource interface {
	Resource
	OpenStream(streamId string) (DataStream, error)
}

type FileStreamSource struct {
	Root string
}

func (src *FileStreamSource) OpenStream(resourceId string) (DataStream, error) {
	path := fmt.Sprint(src.Root, "/", resourceId)
	return os.Open(path)
}

func (src *FileStreamSource) Close() (error) {
	// closes automatically
	return nil
}


type FsError struct {
	Msg string
}

func (e *FsError) Error() string {
	return e.Msg
}

