package yang

import (
	"C"
	"io/ioutil"
)

// For API testing purposes
// When testing a resource loader you can call this function to
// trigger your resource loader to ask for a resource.  Service
// will send you back data as a string

//export yangc2_echo_test
func yangc2_echo_test(source ResourceSource, resourceIdStr *C.char) *C.char {
	resourceId := C.GoString(resourceIdStr)
	if res, err := source.OpenResource(resourceId); err != nil {
		return C.CString(err.Error())
	} else {
		if data, err := ioutil.ReadAll(res); err != nil {
			return C.CString(err.Error())
		} else {
			return C.CString(string(data))
		}
	}
}
