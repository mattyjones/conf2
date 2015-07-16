package c2io

import "C"

// For API testing purposes
// When testing a resource loader you can call this function to
// trigger your resource loader to ask for a resource.  Service
// will send you back data as a string

type echo struct {
	data string
}

//export c2io_echoTest
func c2io_echoTest(source DataSource, resourceId *C.char) *C.char {
	ed := &echo{}
	buff := make([]byte, 1024)
	source.ReadData(ed, buff, C.GoString(resourceId))
	return C.CString(ed.data)
}

func (e *echo) WriteData(buffer []byte) int {
	// simple assumption, assumes data is < buffsize
	e.data = string(buffer)
	return 0
}

