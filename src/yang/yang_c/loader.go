package yang_c

import (
	"C"
	"yang"
)

//export LoadModuleFromCByteArray
func LoadModuleFromCByteArray(cdata *C.char, len C.int) {
	// TODO: improve performance by not copying
	gdata := []byte(C.GoStringN(cdata, len))
	yang.LoadModuleFromByteArray(gdata)
}
