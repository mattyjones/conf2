package c

import "C"

const (
	TRUE_SHORT = C.short(1)
	FALSE_SHORT = C.short(0)
)

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
