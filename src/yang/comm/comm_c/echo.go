package comm_c

//export c2io_echoTest
func c2io_echoTest(source DataSource, resourceId *C.char) *C.char {
	ed := &echo{}
	buff := make([]byte, 1024)
	source.ReadData(ed, buff, C.GoString(resourceId))
	return C.CString(ed.data)
}
