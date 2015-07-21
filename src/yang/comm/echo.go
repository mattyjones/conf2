package comm

import "C"

// For API testing purposes
// When testing a resource loader you can call this function to
// trigger your resource loader to ask for a resource.  Service
// will send you back data as a string

type echo struct {
	data string
}


func (e *echo) WriteData(buffer []byte) int {
	// simple assumption, assumes data is < buffsize
	e.data = string(buffer)
	return 0
}

