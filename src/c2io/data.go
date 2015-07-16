package c2io

// #include "c2io/c2io_stream.h"
// extern int c2io_read_stream(c2io_read_stream_impl impl_func, void *sinkPtr, void *source_handle, void *bufPtr, char *resourceId);
import "C"

import (
	"unsafe"
)

//export DataSink
type DataSink interface {
	WriteData(buffer []byte) int
}

//export DataSink
type DataSource interface {
	ReadData(sink DataSink, buffer []byte, resourceId string)
}


type DriverDataSource struct {
	source_handle unsafe.Pointer
	read_impl C.c2io_read_stream_impl
}

//export c2io_NewDriverDataSource
func c2io_NewDriverDataSource(read_impl C.c2io_read_stream_impl, source_handle unsafe.Pointer) DataSource {
	return &DriverDataSource{
		source_handle: source_handle,
		read_impl: read_impl,
	}
}

func (source *DriverDataSource) ReadData(sink DataSink, buff []byte, resourceId string) {
	sinkPtr := unsafe.Pointer(&sink)
	buffPtr := unsafe.Pointer(&buff)
	C.c2io_read_stream(source.read_impl, sinkPtr, source.source_handle, buffPtr, C.CString(resourceId))
}

//export c2io_DataSink_WriteData
func c2io_DataSink_WriteData(sink DataSink, buff []byte) int {
	return sink.WriteData(buff)
}
