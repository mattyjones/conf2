package yang

import (
	"C"
	"io/ioutil"
	"fmt"
)

//export LoadModuleFromCByteArray
func LoadModuleFromCByteArray(cdata *C.char, len C.int) {
	// improve performance but not copying
	gdata := []byte(C.GoStringN(cdata, len))
	LoadModuleFromByteArray(gdata)
}

func LoadModule(yangfile string) (*Module, error) {
	data, err := ioutil.ReadFile(yangfile)
	if err == nil {
		return LoadModuleFromByteArray(data)
	}
	return nil, err
}

func LoadModuleFromByteArray(data []byte) (*Module, error) {
	l := lex(string(data))
	err_code := yyParse(l)
	if err_code == 0 {
		d := l.stack.Peek()
		return d.(*Module), nil
	}
	return nil, yangError{fmt.Sprintf("Error %d loading yang file", err_code)}
}