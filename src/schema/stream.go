package schema

import (
	"os"
	"fmt"
	"strings"
)

type DataStream interface {
	Read(p []byte) (n int, err error)
}

type StreamSource interface {
	OpenStream(streamId string) (DataStream, error)
}

type FileStreamSource struct {
	Root string
}

func NewCwdSource() StreamSource {
	cwd,_ := os.Getwd()
	return &FileStreamSource{Root:cwd}
}

type StringSource struct {
	Streamer StringStreamer
}

type StringStreamer func(resource string) (string, error)

type stringStream strings.Reader


func (s *StringSource) OpenStream(resourceId string) (DataStream, error) {
	str, err := s.Streamer(resourceId)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(str), nil
}

func (src *FileStreamSource) OpenStream(resourceId string) (DataStream, error) {
	path := fmt.Sprint(src.Root, "/", resourceId)
	return os.Open(path)
}

type FsError struct {
	Msg string
}

func (e *FsError) Error() string {
	return e.Msg
}

