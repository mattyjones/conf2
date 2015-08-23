package driver

// #include "yang/driver/yangc2_stream.h"
// extern void* yangc2_open_stream(yangc2_open_stream_impl impl_func, void *source_handle, char *resourceId, void *fs_err);
// extern void yangc2_close_stream(yangc2_close_stream_impl impl_func, void *stream_handle, void *fs_err);
// extern int yangc2_read_stream(yangc2_read_stream_impl impl_func, void *stream_handle, void *buffPtr, int maxAmout, void *fs_err);
import "C"

import (
	"unsafe"
	"io"
	"yang"
	"fmt"
)

//export ResourceHandle
type ResourceHandle interface {
     yang.ResourceSource
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
func yangc2_new_driver_resource_source(open_impl C.yangc2_open_stream_impl, read_impl C.yangc2_read_stream_impl,
		close_impl C.yangc2_close_stream_impl, source_handle unsafe.Pointer) ResourceHandle {
	return &DriverResourceSource{
		source_handle: source_handle,
		open_impl: open_impl,
		read_impl: read_impl,
		close_impl: close_impl,
	}
}

func (source *DriverResourceSource) OpenResource(resourceId string) (res yang.Resource, err error) {
	errPtr := unsafe.Pointer(&err)
	stream_handle := C.yangc2_open_stream(source.open_impl, source.source_handle, C.CString(resourceId), errPtr)
	if err != nil {
fmt.Println("ERR OpenResource", resourceId);
		return nil, err
	}
fmt.Println("NO ERR OpenResource", resourceId);
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
