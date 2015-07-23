package yang

// #include "yang/yangc2_stream.h"
// extern void* yangc2_open_stream(yangc2_open_stream_impl impl_func, void *source_handle, char *resourceId, void *fs_err);
// extern void yangc2_close_stream(yangc2_close_stream_impl impl_func, void *stream_handle, void *fs_err);
// extern int yangc2_read_stream(yangc2_read_stream_impl impl_func, void *stream_handle, void *buffPtr, int maxAmout, void *fs_err);
import "C"

import (
	"unsafe"
	"os"
	"fmt"
	"io"
)

//export Resource
type Resource interface {
	Read(p []byte) (n int, err error)
	Close() error
}

//export ResourceSource
type ResourceSource interface {
	OpenResource(resourceId string) (Resource, error)
}

type FileDataSource struct {
	Root string
}

// ResourceSource
func (src *FileDataSource) OpenResource(resourceId string) (Resource, error) {
	path := fmt.Sprint(src.Root, resourceId)
	return os.Open(path)
}

type DriverResourceSource struct {
	source_handle unsafe.Pointer
	open_impl C.yangc2_open_stream_impl
	read_impl C.yangc2_read_stream_impl
	close_impl C.yangc2_close_stream_impl
}

type DriverResource struct {
	stream_handle unsafe.Pointer
	read_impl C.yangc2_read_stream_impl
	close_impl C.yangc2_close_stream_impl
}

//export yangc2_new_driver_resource_source
func yangc2_new_driver_resource_source(open_impl C.yangc2_open_stream_impl, read_impl C.yangc2_read_stream_impl, close_impl C.yangc2_close_stream_impl, source_handle unsafe.Pointer) ResourceSource {
	return &DriverResourceSource{
		source_handle: source_handle,
		open_impl: open_impl,
		read_impl: read_impl,
		close_impl: close_impl,
	}
}

type FsError struct {
	Msg string
}

func (e *FsError) Error() string {
	return e.Msg
}

//export yangc2_new_fs_error
func yangc2_new_fs_error(err *C.char) error {
	return &FsError{Msg:C.GoString(err)}
}

func (source *DriverResourceSource) OpenResource(resourceId string) (res Resource, err error) {
	errPtr := unsafe.Pointer(&err)
	stream_handle := C.yangc2_open_stream(source.open_impl, source.source_handle, C.CString(resourceId), errPtr)
	if err != nil {
		return nil, err
	}
	res = &DriverResource{
		stream_handle: stream_handle,
		read_impl: source.read_impl,
		close_impl: source.close_impl,
	}
	return
}

func (res *DriverResource) Read(buff []byte) (n int, err error) {
	errPtr := unsafe.Pointer(&err)
	maxAmount := C.int(len(buff))
	buffPtr := unsafe.Pointer(&buff)
	readAmount := C.yangc2_read_stream(res.read_impl, res.stream_handle, buffPtr, maxAmount, errPtr)
	if readAmount < 0 {
		return 0, io.EOF
	}
	return int(readAmount), err
}

func (res *DriverResource) Close() (err error) {
	errPtr := unsafe.Pointer(&err)
	C.yangc2_close_stream(res.close_impl, res.stream_handle, errPtr)
	return
}
