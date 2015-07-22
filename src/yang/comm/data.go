package comm

// #include "yang/comm/yangc2_comm_stream.h"
// extern int yangc2_comm_read_stream(yangc2_comm_read_stream_impl impl_func, void *sinkPtr, void *source_handle, void *bufPtr, char *resourceId);
import "C"

import (
	"unsafe"
)

//export DataSink
type DataSink interface {
	WriteData(buffer []byte) int
}

//export DataSource
type DataSource interface {
	ReadData(sink DataSink, buffer []byte, resourceId string)
}

type DriverDataSource struct {
	source_handle unsafe.Pointer
	read_impl C.yangc2_comm_read_stream_impl
}

//export yangc2_comm_new_driver_data_source
func yangc2_comm_new_driver_data_source(read_impl C.yangc2_comm_read_stream_impl, source_handle unsafe.Pointer) DataSource {
	return &DriverDataSource{
		source_handle: source_handle,
		read_impl: read_impl,
	}
}

func (source *DriverDataSource) ReadData(sink DataSink, buff []byte, resourceId string) {
	sinkPtr := unsafe.Pointer(&sink)
	buffPtr := unsafe.Pointer(&buff)
	C.yangc2_comm_read_stream(source.read_impl, sinkPtr, source.source_handle, buffPtr, C.CString(resourceId))
}

//export yangc2_comm_datasink_writedata
func yangc2_comm_datasink_writedata(sink DataSink, buff []byte) int {
	return sink.WriteData(buff)
}
