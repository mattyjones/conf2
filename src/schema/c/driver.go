package c

import "C"

type driverError struct {
	Msg string
}

func (e *driverError) Error() string {
	return e.Msg
}

//export conf2_new_driver_error
func conf2_new_driver_error(err *C.char) error {
	return &driverError{Msg:C.GoString(err)}
}
