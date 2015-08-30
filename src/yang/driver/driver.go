package driver

import "C"

type driverError struct {
	Msg string
}

func (e *driverError) Error() string {
	return e.Msg
}

//export yangc2_new_driver_error
func yangc2_new_driver_error(err *C.char) error {
	return &driverError{Msg:C.GoString(err)}
}
