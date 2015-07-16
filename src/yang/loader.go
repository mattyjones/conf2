package yang

import (
	"C"
	"io/ioutil"
	"fmt"
)

//export LoadModuleFromCByteArray
func LoadModuleFromCByteArray(cdata *C.char, len C.int) {
	// TODO: improve performance by not copying
	gdata := []byte(C.GoStringN(cdata, len))
	loadModuleFromByteArray(gdata)
}

func LoadModule(resolver ResourceResolver, yangfile string) (*Module, error) {
	data, err := resolver.LoadResource(yangfile)
	if err == nil {
		return loadModuleFromByteArray(data)
	}
	return nil, err
}

func loadModuleFromByteArray(data []byte) (*Module, error) {
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

type ResourceResolver interface {
	LoadResource(resource string) ([]byte, error)
}

type FileResolver struct {
}

func (fs *FileResolver) LoadResource(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}
