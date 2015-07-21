package comm

type DataSink interface {
	WriteData(buffer []byte) int
}

type DataSource interface {
	ReadData(sink DataSink, buffer []byte, resourceId string)
}
