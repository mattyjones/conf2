package yang

import "C"

import (
	"io/ioutil"
	"fmt"
)

func LoadModuleFromByteArray(data []byte) (*Module, error) {
	l := lex(string(data))
	err_code := yyParse(l)
	if err_code != 0 || l.lastError != nil {
		if l.lastError == nil {
			// Developer - Find out why there's no error
			msg := fmt.Sprint("Error parsing, code ", string(err_code))
			l.lastError = &yangError{msg}

		}
		return nil, l.lastError
	}

	d := l.stack.Peek()
	return d.(*Module), nil
}

func LoadModule(source ResourceSource, yangfile string) (*Module, error) {
	if res, err := source.OpenResource(yangfile); err != nil {
		return nil, err
	} else {
		defer res.Close()
		if data, err := ioutil.ReadAll(res); err != nil {
			return nil, err
		} else {
			return LoadModuleFromByteArray(data)
		}
	}
}

//export yangc2_load_module_from_resource_source
func yangc2_load_module_from_resource_source(source ResourceSource, resourceStr *C.char) error {
	resourceId := C.GoString(resourceStr)
	_, err := LoadModule(source, resourceId)
	return err
}
