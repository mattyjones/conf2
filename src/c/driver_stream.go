package c

// #include "conf2/stream.h"
// extern void* conf2_open_stream(conf2_open_stream_impl impl_func, void *source_handle, char *resourceId, void *fs_err);
// extern int conf2_read_stream(conf2_read_stream_impl impl_func, void *stream_handle, void *buffPtr, int maxAmout, void *fs_err);
import "C"

import (
	"unsafe"
	"io"
	"schema"
)

type DriverStreamSource struct {
	sourceHandle *ApiHandle
	open_impl C.conf2_open_stream_impl
	read_impl C.conf2_read_stream_impl
}

func (dss *DriverStreamSource) Close() error {
	return dss.sourceHandle.Close()
}

type DriverStream struct {
	streamHandle *ApiHandle
	read_impl C.conf2_read_stream_impl
}


func (res *DriverStream) Close() error {
	return res.streamHandle.Close()
}

//export conf2_new_driver_resource_source
func conf2_new_driver_resource_source(stream_source_hnd_id unsafe.Pointer, open_impl C.conf2_open_stream_impl,
		read_impl C.conf2_read_stream_impl) unsafe.Pointer {
	stream_source_hnd, found := ApiHandles()[stream_source_hnd_id]
	if !found {
		panic("stream source handle not found")
	}
	return NewGoHandle(&DriverStreamSource{
		sourceHandle : stream_source_hnd,
		open_impl : open_impl,
		read_impl : read_impl,
	}).ID
}

func (source *DriverStreamSource) OpenStream(resourceId string) (res schema.DataStream, err error) {
	errPtr := unsafe.Pointer(&err)
	streamHandleId := C.conf2_open_stream(source.open_impl, source.sourceHandle.ID, C.CString(resourceId), errPtr)
	if err != nil {
		return nil, err
	}
	streamHandle, found := ApiHandles()[streamHandleId]
	if !found {
		panic("Stream handle not found")
	}
	res = &DriverStream{
		streamHandle: streamHandle,
		read_impl: source.read_impl,
	}
	return
}

func (res *DriverStream) Read(buff []byte) (n int, err error) {
	errPtr := unsafe.Pointer(&err)
	maxAmount := C.int(len(buff))
	buffPtr := unsafe.Pointer(&buff)
	readAmount := C.conf2_read_stream(res.read_impl, res.streamHandle.ID, buffPtr, maxAmount, errPtr)
	if readAmount < 0 {
		return 0, io.EOF
	}
	return int(readAmount), err
}
